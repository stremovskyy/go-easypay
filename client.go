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

package go_easypay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/stremovskyy/go-easypay/consts"
	"github.com/stremovskyy/go-easypay/easypay"
	"github.com/stremovskyy/go-easypay/internal/http"
	"github.com/stremovskyy/go-easypay/log"
	"github.com/stremovskyy/recorder"
)

type client struct {
	easypayClient *http.Client
	app           *easypay.App
}

func (c *client) SetLogLevel(levelDebug log.Level) {
	log.SetLevel(levelDebug)
}

func NewDefaultClient() Easypay {
	return &client{
		easypayClient: http.NewClient(http.DefaultOptions()),
	}
}

func NewClientWithRecorder(rec recorder.Recorder) Easypay {
	return &client{
		easypayClient: http.NewClient(http.DefaultOptions()).WithRecorder(rec),
	}
}

func NewClient(options ...Option) Easypay {
	c := &client{
		easypayClient: http.NewClient(http.DefaultOptions()),
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *client) VerificationLink(request *Request) (*url.URL, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	if c.app == nil || !c.app.IsValid() {
		err := c.createApp(request.Merchant)
		if err != nil {
			return nil, fmt.Errorf("cannot create App: %v", err)
		}
	}

	pageID, err := c.createPageID(request)
	if err != nil {
		return nil, fmt.Errorf("cannot create Page ID: %v", err)
	}

	createTokenRequest := easypay.NewRequest(
		consts.CardTokenCreateURL,
		easypay.WithPartnerKeyHeader(request.Merchant.getPartnerKey()),
		easypay.WithAppIDHeader(c.app.AppID()),
		easypay.WithPageIDHeader(pageID),
		easypay.WithSecretKey(request.Merchant.GetSecretKey()),
		easypay.WithPhone(request.GetPaymentID()),
		easypay.WithRedirects(request.GetRedirects()),
		easypay.WithWebhook(request.GetWebhookURL()),
	)

	apiResponse, err := c.easypayClient.Api(createTokenRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get API response: %v", err)
	}

	u, err := url.Parse(apiResponse.ForwardUrl)
	if err != nil {
		return nil, fmt.Errorf("cannot parse URL: %v", err)
	}

	return u, nil
}

func (c *client) Status(request *Request) (*easypay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	if c.app == nil || !c.app.IsValid() {
		err := c.createApp(request.Merchant)
		if err != nil {
			return nil, fmt.Errorf("cannot create App: %v", err)
		}
	}

	pageID, err := c.createPageID(request)
	if err != nil {
		return nil, fmt.Errorf("cannot create Page ID: %v", err)
	}

	cancelRequest := easypay.NewRequest(
		consts.CheckOrderStateURL,
		easypay.WithPartnerKeyHeader(request.Merchant.getPartnerKey()),
		easypay.WithAppIDHeader(c.app.AppID()),
		easypay.WithPageIDHeader(pageID),
		easypay.WithSecretKey(request.Merchant.GetSecretKey()),
		easypay.WithRootServiceKey(request.Merchant.GetServiceKey()),
		easypay.WithTransactionID(request.GetTransactionID()),
		easypay.WithRootOrderID(request.GetPaymentID()),
		easypay.WithoutError(),
	)

	apiResponse, err := c.easypayClient.Api(cancelRequest)
	if err != nil {
		return nil, fmt.Errorf("error while creating payment: %v", err)
	}

	return apiResponse, nil
}

func (c *client) PaymentURL(request *Request) (*easypay.Response, error) {
	panic("implement me")
}

func (c *client) Payment(request *Request) (*easypay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	if c.app == nil || !c.app.IsValid() {
		err := c.createApp(request.Merchant)
		if err != nil {
			return nil, fmt.Errorf("cannot create App: %v", err)
		}
	}

	pageID, err := c.createPageID(request)
	if err != nil {
		return nil, fmt.Errorf("cannot create Page ID: %v", err)
	}

	requestOptions := []func(*easypay.Request){
		easypay.WithPartnerKeyHeader(request.Merchant.getPartnerKey()),
		easypay.WithAppIDHeader(c.app.AppID()),
		easypay.WithPageIDHeader(pageID),
		easypay.WithSecretKey(request.Merchant.GetSecretKey()),
		easypay.WithAdditionalWebhook(request.GetWebhookURL()),
		easypay.WithOrderID(request.GetPaymentID()),
		easypay.WithAmount(request.GetAmount()),
		easypay.WithDescription(request.GetDescription()),
		easypay.WithServiceKey(request.Merchant.GetServiceKey()),
		easypay.WithOneTimePayment(true),
		easypay.WithBankingDetails(request.GetBankingDetails()),
	}

	if request.IsMobile() {
		if request.IsApplePay() {
			requestOptions = append(requestOptions, easypay.WithApplePayContainer(request.GetAppleContainer()))
			requestOptions = append(requestOptions, easypay.WithPaymentInstrumentMerchantID(request.GetAppleMerchantID()))
		} else {
			requestOptions = append(requestOptions, easypay.WithGooglePayToken(request.GetGoogleToken()))
		}
	} else {
		requestOptions = append(requestOptions, easypay.WithCardToken(request.GetCardToken()))
		requestOptions = append(requestOptions, easypay.WithCardTokenID(request.GetCardTokenID()))
	}

	paymentRequest := easypay.NewRequest(
		consts.CreateOrderURL,
		requestOptions...,
	)

	apiResponse, err := c.easypayClient.Api(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("error while creating payment: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Hold(request *Request) (*easypay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	if c.app == nil || !c.app.IsValid() {
		err := c.createApp(request.Merchant)
		if err != nil {
			return nil, fmt.Errorf("cannot create App: %v", err)
		}
	}

	pageID, err := c.createPageID(request)
	if err != nil {
		return nil, fmt.Errorf("cannot create Page ID: %v", err)
	}

	requestOptions := []func(*easypay.Request){
		easypay.WithPaymentOperation(consts.PaymentOperationPaymentHold),
		easypay.WithPartnerKeyHeader(request.Merchant.getPartnerKey()),
		easypay.WithAppIDHeader(c.app.AppID()),
		easypay.WithPageIDHeader(pageID),
		easypay.WithSecretKey(request.Merchant.GetSecretKey()),
		easypay.WithAdditionalWebhook(request.GetWebhookURL()),
		easypay.WithOrderID(request.GetPaymentID()),
		easypay.WithAmount(request.GetAmount()),
		easypay.WithDescription(request.GetDescription()),
		easypay.WithServiceKey(request.Merchant.GetServiceKey()),
		easypay.WithOneTimePayment(true),
		easypay.WithBankingDetails(request.GetBankingDetails()),
	}

	if request.IsMobile() {
		if request.IsApplePay() {
			requestOptions = append(requestOptions, easypay.WithApplePayContainer(request.GetAppleContainer()))
			requestOptions = append(requestOptions, easypay.WithPaymentInstrumentMerchantID(request.GetAppleMerchantID()))
		} else {
			requestOptions = append(requestOptions, easypay.WithGooglePayToken(request.GetGoogleToken()))
		}
	} else {
		requestOptions = append(requestOptions, easypay.WithCardToken(request.GetCardToken()))
		requestOptions = append(requestOptions, easypay.WithCardTokenID(request.GetCardTokenID()))
	}

	holdRequest := easypay.NewRequest(
		consts.CreateOrderURL,
		requestOptions...,
	)

	apiResponse, err := c.easypayClient.Api(holdRequest)
	if err != nil {
		return nil, fmt.Errorf("error while creating payment: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Capture(request *Request) (*easypay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	if c.app == nil || !c.app.IsValid() {
		err := c.createApp(request.Merchant)
		if err != nil {
			return nil, fmt.Errorf("cannot create App: %v", err)
		}
	}

	pageID, err := c.createPageID(request)
	if err != nil {
		return nil, fmt.Errorf("cannot create Page ID: %v", err)
	}

	CaptureRequest := easypay.NewRequest(
		consts.UnHoldURL,
		easypay.WithPartnerKeyHeader(request.Merchant.getPartnerKey()),
		easypay.WithAppIDHeader(c.app.AppID()),
		easypay.WithPageIDHeader(pageID),
		easypay.WithSecretKey(request.Merchant.GetSecretKey()),
		easypay.WithTransactionID(request.GetTransactionID()),
		easypay.WithRootAmount(request.GetAmount()),
		easypay.WithRootOrderID(request.GetPaymentID()),
		easypay.WithWebhook(request.GetWebhookURL()),
		easypay.WithRootServiceKey(request.Merchant.GetServiceKey()),
	)

	apiResponse, err := c.easypayClient.Api(CaptureRequest)
	if err != nil {
		return nil, fmt.Errorf("error while creating payment: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Refund(request *Request) (*easypay.Response, error) {
	if request == nil {
		return nil, ErrRequestIsNil
	}

	if c.app == nil || !c.app.IsValid() {
		err := c.createApp(request.Merchant)
		if err != nil {
			return nil, fmt.Errorf("cannot create App: %v", err)
		}
	}

	pageID, err := c.createPageID(request)
	if err != nil {
		return nil, fmt.Errorf("cannot create Page ID: %v", err)
	}

	cancelRequest := easypay.NewRequest(
		consts.CancelOrderURL,
		easypay.WithPartnerKeyHeader(request.Merchant.getPartnerKey()),
		easypay.WithAppIDHeader(c.app.AppID()),
		easypay.WithPageIDHeader(pageID),
		easypay.WithSecretKey(request.Merchant.GetSecretKey()),
		easypay.WithRootServiceKey(request.Merchant.GetServiceKey()),
		easypay.WithTransactionID(request.GetTransactionID()),
		easypay.WithRootOrderID(request.GetPaymentID()),
		easypay.WithRootAmount(request.GetAmount()),
		easypay.WithWebhook(request.GetWebhookURL()),
	)

	apiResponse, err := c.easypayClient.Api(cancelRequest)
	if err != nil {
		return nil, fmt.Errorf("error while creating payment: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Credit(request *Request) (*easypay.Response, error) {
	panic("implement me")
}

func (c *client) createApp(merchant *Merchant) error {
	createAppRequest := easypay.NewRequest(
		consts.CreateAppURL,
		easypay.WithPartnerKeyHeader(merchant.getPartnerKey()),
	)

	response, err := c.easypayClient.NotRecordedApi(createAppRequest)
	if err != nil {
		return fmt.Errorf("cannot get App response: %v", err)
	}

	c.app = response.App()

	return nil
}

func (c *client) createPageID(request *Request) (*string, error) {
	pageRequest := easypay.NewRequest(
		consts.CreatePageURL,
		easypay.WithPartnerKeyHeader(request.Merchant.getPartnerKey()),
		easypay.WithAppIDHeader(c.app.AppID()),
	)

	response, err := c.easypayClient.NotRecordedApi(pageRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get Page response: %v", err)
	}

	return response.PageId, nil
}

func (c *client) GetRecordedExchange(ctx context.Context, requestID string) (*easypay.RecordedExchange, error) {
	if c.easypayClient == nil || c.easypayClient.GetRecorder() == nil {
		return nil, fmt.Errorf("recorder is not configured")
	}

	if requestID == "" {
		return nil, fmt.Errorf("requestID cannot be empty")
	}

	recorder := c.easypayClient.GetRecorder()

	request, err := recorder.GetRequest(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get request for ID %s: %w", requestID, err)
	}

	response, err := recorder.GetResponse(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get response for ID %s: %w", requestID, err)
	}

	exchange := &easypay.RecordedExchange{
		RequestID: requestID,
		Request:   request,
		Response:  response,
		Tags:      c.extractTagsFromRequest(request),
		Timestamp: time.Now(),
	}

	return exchange, nil
}

func (c *client) extractTagsFromRequest(requestData []byte) map[string]string {
	tags := make(map[string]string)

	if len(requestData) == 0 {
		return tags
	}

	var requestObj map[string]interface{}
	if err := json.Unmarshal(requestData, &requestObj); err != nil {
		return tags
	}

	if orderID, ok := requestObj["order_id"].(string); ok && orderID != "" {
		tags["order_id"] = orderID
	}

	if transactionID, ok := requestObj["transaction_id"].(float64); ok {
		tags["transaction_id"] = fmt.Sprintf("%.0f", transactionID)
	}

	if url, ok := requestObj["url"].(string); ok && url != "" {
		tags["url"] = url
	}

	return tags
}

func (c *client) GetExchangesByOrderID(ctx context.Context, orderID string) ([]*easypay.RecordedExchange, error) {
	if c.easypayClient == nil || c.easypayClient.GetRecorder() == nil {
		return nil, fmt.Errorf("recorder is not configured")
	}

	if orderID == "" {
		return nil, fmt.Errorf("orderID cannot be empty")
	}

	recorder := c.easypayClient.GetRecorder()

	// Use the correct tag:value format for FindByTag
	redisKeys, err := recorder.FindByTag(ctx, "order_id:"+orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find exchanges by order ID %s: %w", orderID, err)
	}

	if len(redisKeys) == 0 {
		return []*easypay.RecordedExchange{}, nil
	}

	// Extract unique request IDs from the Redis keys
	requestIDSet := make(map[string]bool)
	for _, key := range redisKeys {
		requestID := c.extractRequestIDFromKey(key)
		if requestID != "" {
			requestIDSet[requestID] = true
		}
	}

	exchanges := make([]*easypay.RecordedExchange, 0, len(requestIDSet))
	var errors []error

	for requestID := range requestIDSet {
		exchange, err := c.GetRecordedExchange(ctx, requestID)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to get exchange %s: %w", requestID, err))
			log.NewLogger("easypay").Error("failed to get exchange", "requestID", requestID, "error", err)
			continue
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}

func (c *client) GetExchangesByTransactionID(ctx context.Context, transactionID string) ([]*easypay.RecordedExchange, error) {
	if c.easypayClient == nil || c.easypayClient.GetRecorder() == nil {
		return nil, fmt.Errorf("recorder is not configured")
	}

	if transactionID == "" {
		return nil, fmt.Errorf("transactionID cannot be empty")
	}

	recorder := c.easypayClient.GetRecorder()

	// Use the correct tag:value format for FindByTag
	redisKeys, err := recorder.FindByTag(ctx, "transaction_id:"+transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find exchanges by transaction ID %s: %w", transactionID, err)
	}

	if len(redisKeys) == 0 {
		return []*easypay.RecordedExchange{}, nil
	}

	// Extract unique request IDs from the Redis keys
	requestIDSet := make(map[string]bool)
	for _, key := range redisKeys {
		requestID := c.extractRequestIDFromKey(key)
		if requestID != "" {
			requestIDSet[requestID] = true
		}
	}

	exchanges := make([]*easypay.RecordedExchange, 0, len(requestIDSet))
	var errors []error

	for requestID := range requestIDSet {
		exchange, err := c.GetRecordedExchange(ctx, requestID)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to get exchange %s: %w", requestID, err))
			log.NewLogger("easypay").Error("failed to get exchange", "requestID", requestID, "error", err)
			continue
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}

func (c *client) extractRequestIDFromKey(key string) string {
	if key == "" {
		return ""
	}

	parts := strings.Split(key, ":")
	if len(parts) == 0 {
		return ""
	}

	lastPart := parts[len(parts)-1]

	if len(lastPart) > 8 && (strings.Contains(lastPart, "-") || len(lastPart) == 36) {
		return lastPart
	}

	return ""
}
