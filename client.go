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
	"fmt"
	"net/url"

	"github.com/stremovskyy/go-easypay/consts"
	"github.com/stremovskyy/go-easypay/easypay"
	"github.com/stremovskyy/go-easypay/internal/http"
	"github.com/stremovskyy/go-easypay/log"
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
		consts.CardTokenCreate,
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
	panic("implement me")
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

	paymentRequest := easypay.NewRequest(
		consts.CreateOrderURL,
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
		easypay.WithCardToken(request.GetCardToken()),
		easypay.WithCardTokenID(request.GetCardTokenID()),
	)

	apiResponse, err := c.easypayClient.Api(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("error while creating payment: %v", err)
	}

	return apiResponse, nil
}

func (c *client) Hold(request *Request) (*easypay.Response, error) {
	panic("implement me")
}

func (c *client) Capture(invoiceRequest *Request) (*easypay.Response, error) {
	panic("implement me")
}

func (c *client) Refund(request *Request) (*easypay.Response, error) {
	panic("implement me")
}

func (c *client) Credit(request *Request) (*easypay.Response, error) {
	panic("implement me")
}

func (c *client) createApp(merchant *Merchant) error {
	createAppRequest := easypay.NewRequest(
		consts.CreateAppURL,
		easypay.WithPartnerKeyHeader(merchant.getPartnerKey()),
	)

	response, err := c.easypayClient.Api(createAppRequest)
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

	response, err := c.easypayClient.Api(pageRequest)
	if err != nil {
		return nil, fmt.Errorf("cannot get Page response: %v", err)
	}

	return response.PageId, nil
}
