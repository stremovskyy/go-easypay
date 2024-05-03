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
	"github.com/stremovskyy/go-easypay/currency"
	"github.com/stremovskyy/go-easypay/easypay"
)

type Request struct {
	Merchant      *Merchant
	PersonalData  *PersonalData
	PaymentData   *PaymentData
	PaymentMethod *PaymentMethod
}

func (r *Request) GetRedirects() (string, string) {
	if r.Merchant == nil {
		return "", ""
	}

	return r.Merchant.SuccessRedirect, r.Merchant.FailRedirect
}

func (r *Request) GetEasypayPaymentID() int64 {
	if r.PaymentData == nil || r.PaymentData.EasypayPaymentID == nil {
		return 0
	}

	return *r.PaymentData.EasypayPaymentID
}

func (r *Request) GetCardToken() *string {
	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil {
		return nil
	}

	return r.PaymentMethod.Card.Token
}

func (r *Request) GetPaymentID() *string {
	if r.PaymentData == nil {
		return nil
	}

	return r.PaymentData.PaymentID
}

func (r *Request) SetRedirects(successURL string, failURL string) {
	if r.Merchant == nil {
		r.Merchant = &Merchant{}
	}

	r.Merchant.SuccessRedirect = successURL
	r.Merchant.FailRedirect = failURL
}

func (r *Request) GetWebhookURL() *string {
	if r.PaymentData == nil {
		return nil
	}

	return r.PaymentData.WebhookURL
}

func (r *Request) SetWebhookURL(webhookURL *string) {
	if r.PaymentData == nil {
		r.PaymentData = &PaymentData{}
	}

	r.PaymentData.WebhookURL = webhookURL
}

func (r *Request) GetAmount() float64 {
	if r.PaymentData == nil {
		return 0
	}

	return r.PaymentData.Amount

}

func (r *Request) GetDescription() string {
	if r.PaymentData == nil {
		return ""
	}

	return r.PaymentData.Description
}

func (r *Request) GetCurrency() currency.Code {
	if r.PaymentData == nil {
		return ""
	}

	return r.PaymentData.Currency

}

func (r *Request) IsMobile() bool {
	if r.PaymentData == nil {
		return false
	}

	return r.PaymentData.IsMobile || r.PaymentMethod.AppleContainer != nil
}

func (r *Request) GetAppleContainer() *string {
	if r.PaymentMethod == nil || r.PaymentMethod.AppleContainer == nil {
		return nil
	}

	return r.PaymentMethod.AppleContainer
}

func (r *Request) IsApplePay() bool {
	return r.PaymentMethod != nil && r.PaymentMethod.AppleContainer != nil
}

func (r *Request) GetGoogleToken() *string {
	if r.PaymentMethod == nil || r.PaymentMethod.GoogleToken == nil {
		return nil
	}

	return r.PaymentMethod.GoogleToken
}

func (r *Request) GetCardTokenID() string {
	if r.PaymentMethod == nil || r.PaymentMethod.Card == nil {
		return ""
	}

	return r.PaymentMethod.Card.Name
}

func (r *Request) GetBankingDetails() *easypay.BankingDetails {
	return &easypay.BankingDetails{
		Payee: easypay.Payee{
			ID:   r.Merchant.PayeeID,
			Name: r.Merchant.PayeeName,
			Bank: easypay.Bank{
				Account: r.Merchant.PayeeBankAccount,
			},
		},
		Payer: easypay.Payer{
			Name: r.Merchant.PayerName,
		},
		Narrative: easypay.Narrative{
			Name: r.Merchant.PayeeNarative,
		},
	}
}
