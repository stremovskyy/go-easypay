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

package easypay

import (
	"github.com/stremovskyy/go-easypay/consts"
	"github.com/stremovskyy/go-easypay/internal/utils"
)

func NewRequest(url string, options ...func(request *Request)) *Request {
	rw := &Request{
		Url:     url,
		Headers: map[string]string{},
	}

	for _, option := range options {
		option(rw)
	}

	return rw
}

func WithAmount(a float64) func(*Request) {
	return func(rw *Request) {
		if rw.Order == nil {
			rw.Order = &Order{}
		}

		rw.Order.Amount = &a
	}
}
func WithRootAmount(a float64) func(*Request) {
	return func(rw *Request) {
		rw.Amount = &a
	}
}

func WithPartnerKeyHeader(key string) func(request *Request) {
	return func(rw *Request) {
		rw.Headers["PartnerKey"] = key
	}
}

func WithAppIDHeader(id string) func(request *Request) {
	return func(rw *Request) {
		rw.Headers["AppId"] = id
	}
}

func WithPageIDHeader(id *string) func(request *Request) {
	return func(rw *Request) {
		if id != nil {
			rw.Headers["PageId"] = *id
		}
	}
}

func WithSecretKey(key string) func(request *Request) {
	return func(rw *Request) {
		rw.SecretKey = key
	}
}

func WithPhone(s *string) func(request *Request) {
	return func(rw *Request) {
		rw.Phone = s
	}
}

func WithRedirects(successURL string, failURL string) func(request *Request) {
	return func(rw *Request) {
		if rw.URLs == nil {
			rw.URLs = &URLs{}
		}

		rw.URLs.Success = &successURL
		rw.URLs.Failed = &failURL
	}
}

func WithWebhook(url *string) func(request *Request) {
	return func(rw *Request) {
		if rw.URLs == nil {
			rw.URLs = &URLs{}
		}
		rw.URLs.Notify = url
	}
}

func WithAdditionalWebhook(url *string) func(request *Request) {
	return func(rw *Request) {
		if rw.Order == nil {
			rw.Order = &Order{
				AdditionalItems: utils.Ref(make(map[string]string)),
			}
		}

		if rw.Order.AdditionalItems == nil {
			rw.Order.AdditionalItems = utils.Ref(make(map[string]string))
		}

		(*rw.Order.AdditionalItems)["Merchant.UrlNotify"] = *url
	}
}

func WithOrderID(id *string) func(request *Request) {
	return func(rw *Request) {
		if rw.Order == nil {
			rw.Order = &Order{}
		}

		rw.Order.OrderID = id
	}
}

func WithRootOrderID(id *string) func(request *Request) {
	return func(rw *Request) {
		rw.OrderID = id
	}
}

func WithDescription(description string) func(request *Request) {
	return func(rw *Request) {
		if rw.Order == nil {
			rw.Order = &Order{}
		}

		rw.Order.Description = &description
	}
}

func WithServiceKey(key string) func(request *Request) {
	return func(rw *Request) {
		if rw.Order == nil {
			rw.Order = &Order{}
		}

		rw.Order.ServiceKey = &key
	}
}

func WithRootServiceKey(key string) func(request *Request) {
	return func(rw *Request) {
		rw.ServiceKey = &key
	}
}

func WithCardToken(token *string) func(request *Request) {
	return func(rw *Request) {
		if rw.UserPaymentInstrument == nil {
			rw.UserPaymentInstrument = &UserPaymentInstrument{}
		}

		rw.UserPaymentInstrument.InstrumentType = utils.Ref("Card")
		rw.UserPaymentInstrument.CardGuid = token
	}
}

func WithOneTimePayment(b bool) func(request *Request) {
	return func(rw *Request) {
		if rw.Order == nil {
			rw.Order = &Order{}
		}

		rw.Order.IsOneTimePay = &b
	}
}

func WithPaymentOperation(operation consts.PaymentOperation) func(request *Request) {
	return func(rw *Request) {
		if rw.Order == nil {
			rw.Order = &Order{}
		}

		rw.Order.PaymentOperation = utils.Ref(string(operation))
	}
}

func WithCardTokenID(id string) func(request *Request) {
	return func(rw *Request) {
		if rw.UserInfo == nil {
			rw.UserInfo = &UserInfo{}
		}

		rw.UserInfo.Phone = &id
	}
}

func WithBankingDetails(details *BankingDetails) func(request *Request) {
	return func(rw *Request) {
		if rw.UserPaymentInstrument == nil {
			rw.UserPaymentInstrument = &UserPaymentInstrument{}
		}

		rw.BankingDetails = details
	}
}

func WithTransactionID(transactionID *int64) func(request *Request) {
	return func(rw *Request) {
		rw.TransactionID = transactionID
	}
}

// setPaymentInstrument sets the payment instrument type and token on the Request object.
func setPaymentInstrument(rw *Request, instrumentType string, token *string) {
	if rw.UserPaymentInstrument == nil {
		rw.UserPaymentInstrument = &UserPaymentInstrument{}
	}

	rw.UserPaymentInstrument.InstrumentType = utils.Ref(instrumentType)
	rw.UserPaymentInstrument.Token = token
}

func WithGooglePayToken(token *string) func(*Request) {
	return func(rw *Request) {
		setPaymentInstrument(rw, "GooglePay", token)
	}
}

func WithApplePayContainer(container *string) func(*Request) {
	return func(rw *Request) {
		setPaymentInstrument(rw, "ApplePay", container)
	}
}

func WithPaymentInstrumentMerchantID(id *string) func(*Request) {
	return func(rw *Request) {
		if rw.UserPaymentInstrument == nil {
			rw.UserPaymentInstrument = &UserPaymentInstrument{}
		}

		rw.UserPaymentInstrument.GatewayMerchantId = id
	}
}
