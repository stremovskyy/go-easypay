/*
 * MIT License
 *
 * Copyright (c) 2024 Anton Stremovskyy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/stremovskyy/go-easypay/consts"
	"github.com/stremovskyy/go-easypay/easypay"
	"github.com/stremovskyy/go-easypay/log"
	"github.com/stremovskyy/recorder"
)

type Client struct {
	client         *http.Client
	options        *Options
	logger         *log.Logger
	xmlLogger      *log.Logger
	applePayLogger *log.Logger
	recorder       recorder.Recorder
}

func (c *Client) Api(apiRequest *easypay.Request) (*easypay.Response, error) {
	return c.sendRequest(apiRequest, c.logger, true)
}

func (c *Client) NotRecordedApi(apiRequest *easypay.Request) (*easypay.Response, error) {
	return c.sendRequest(apiRequest, c.logger, false)
}

func (c *Client) sendRequest(apiRequest *easypay.Request, logger *log.Logger, record bool) (*easypay.Response, error) {
	requestID := uuid.New().String()
	logger.Debug("Request ID: %v", requestID)
	logger.Debug("Request URL: %v", apiRequest.Url)

	jsonBody, err := json.Marshal(apiRequest)
	if err != nil {
		logger.Error("cannot marshal request: %v", err)
		return nil, fmt.Errorf("cannot marshal request: %v", err)
	}
	if jsonBody != nil {
		logger.Debug("Request: %v", string(jsonBody))
	}

	ctx := context.WithValue(context.Background(), "request_id", requestID)
	tags := tagsRetriever(apiRequest)

	if record {
		metricsMap := make(map[string]string)
		metricsMap["url"] = apiRequest.Url
		tim := time.Now()
		metricsMap["start_timestamp"] = tim.Format("2006-01-02 15:04:05")
		defer func() {
			metricsMap["end_timestamp"] = time.Now().Format("2006-01-02 15:04:05")
			metricsMap["duration"] = fmt.Sprintf("%s", time.Since(tim).String())

			if c.recorder != nil {
				err := c.recorder.RecordMetrics(ctx, nil, requestID, metricsMap, tags)
				if err != nil {
					c.logger.Error("cannot record metrics: %v", err)
				}
			}
		}()
	}

	req, err := http.NewRequest("POST", apiRequest.Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error("cannot create request: %v", err)
		return nil, err
	}

	signature := computeSignature(apiRequest.SecretKey, string(jsonBody))

	c.setHeaders(req, requestID, apiRequest.Headers, signature)

	if record && c.recorder != nil {
		err = c.recorder.RecordRequest(ctx, nil, requestID, jsonBody, tags)
		if err != nil {
			logger.Error("cannot record request: %v", err)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if record && c.recorder != nil {
			err = c.recorder.RecordError(ctx, nil, requestID, err, tags)
			if err != nil {
				logger.Error("cannot record error: %v", err)
			}
		}

		logger.Error("cannot send request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("cannot read response: %v", err)
		return nil, err
	}

	logger.Debug("Response: %v", string(raw))
	logger.Debug("Response status: %v", resp.StatusCode)

	if record && c.recorder != nil {
		err = c.recorder.RecordResponse(ctx, nil, requestID, raw, tags)
		if err != nil {
			logger.Error("cannot record response: %v", err)
		}
	}

	response, err := easypay.UnmarshalJSONResponse(raw)
	if err != nil {
		logger.Error("cannot unmarshal response: %v", err)
		return nil, err
	}

	if !apiRequest.SkipGeneratingError && response.GetError() != nil {
		return nil, fmt.Errorf("easypay error: %v", response.GetError())
	}

	return response, nil
}

func tagsRetriever(request *easypay.Request) map[string]string {
	tags := make(map[string]string)

	if request.OrderID != nil {
		tags["order_id"] = fmt.Sprintf("%s", *request.OrderID)
	}

	if request.TransactionID != nil {
		tags["transaction_id"] = fmt.Sprintf("%d", *request.TransactionID)
	}

	return tags
}

func (c *Client) setHeaders(req *http.Request, requestID string, headers map[string]string, signature string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("locale", "ua")
	req.Header.Set("User-Agent", "GO EASYPAY/"+consts.Version)
	req.Header.Set("X-Request-ID", requestID)
	req.Header.Set("Api-Version", consts.ApiVersion)
	req.Header.Set("Sign", signature)

	if headers != nil {
		for k, v := range headers {
			c.logger.Debug("Header: %v: %v", k, v)
			req.Header.Set(k, v)
		}
	}
}

func (c *Client) SetClient(cl *http.Client) {
	c.client = cl
}

func (c *Client) WithRecorder(rec recorder.Recorder) *Client {
	c.recorder = rec

	return c
}

func (c *Client) SetRecorder(r recorder.Recorder) {
	c.recorder = r
}

func NewClient(options *Options) *Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: options.KeepAlive,
	}

	tr := &http.Transport{
		MaxIdleConns:       options.MaxIdleConns,
		IdleConnTimeout:    options.IdleConnTimeout,
		DisableCompression: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
	}

	cl := &http.Client{
		Transport: tr,
		Timeout:   options.Timeout,
	}

	return &Client{
		client:         cl,
		options:        options,
		logger:         log.NewLogger("easypay HTTP:"),
		applePayLogger: log.NewLogger("easypay HTTP:"),
		xmlLogger:      log.NewLogger("easypay HTTP XML:"),
	}
}

func computeSignature(secretKey, requestBody string) string {
	// Concatenate the secret key and request body
	data := secretKey + requestBody

	// Calculate the SHA256 hash
	hash := sha256.New()
	hash.Write([]byte(data))
	hashedData := hash.Sum(nil)

	// Encode the hashed data to base64
	signature := base64.StdEncoding.EncodeToString(hashedData)

	return signature
}
