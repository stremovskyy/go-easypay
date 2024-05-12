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

package main

import (
	"fmt"

	"github.com/google/uuid"

	go_easypay "github.com/stremovskyy/go-easypay"
	"github.com/stremovskyy/go-easypay/currency"
	"github.com/stremovskyy/go-easypay/internal/utils"
	"github.com/stremovskyy/go-easypay/log"
	"github.com/stremovskyy/go-easypay/private"
)

func main() {
	client := go_easypay.NewDefaultClient()

	merchant := &go_easypay.Merchant{
		Name:             private.MerchantName,
		PartnerKey:       private.PartnerKey,
		ServiceKey:       private.ServiceKey,
		SecretKey:        private.SecretKey,
		SuccessRedirect:  private.SuccessRedirect,
		FailRedirect:     private.FailRedirect,
		PayeeID:          private.PayeeID,
		PayeeName:        private.PayeeName,
		PayeeBankAccount: private.PayeeBankAccount,
		PayeeNarative:    private.PayeeNarative,
		PayerName:        private.PayerName,
	}

	uuidString := uuid.New().String()

	paymentRequest := &go_easypay.Request{
		Merchant: merchant,
		PaymentMethod: &go_easypay.PaymentMethod{
			AppleContainer: utils.Ref(private.AppleContainer),
		},
		PaymentData: &go_easypay.PaymentData{
			PaymentID:   utils.Ref(uuidString),
			Amount:      1.0,
			Currency:    currency.UAH,
			OrderID:     uuidString,
			Description: "Test payment: " + uuidString,
			IsMobile:    true,
		},
		PersonalData: &go_easypay.PersonalData{
			UserID:    utils.Ref(123),
			FirstName: utils.Ref("John"),
			LastName:  utils.Ref("Doe"),
			TaxID:     utils.Ref("1234567890"),
		},
	}

	client.SetLogLevel(log.LevelDebug)
	paymentRequest.SetWebhookURL(utils.Ref(private.WebhookURL))

	paymentResponse, err := client.Payment(paymentRequest)
	if err != nil {
		panic(err)
	}

	if paymentResponse.GetError() != nil {
		panic(paymentResponse.GetError())
	}

	fmt.Printf("Payment: %s is %s", uuidString, paymentResponse.PaymentState)
}
