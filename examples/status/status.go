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

	go_easypay "github.com/stremovskyy/go-easypay"
	"github.com/stremovskyy/go-easypay/internal/utils"
	"github.com/stremovskyy/go-easypay/log"
	"github.com/stremovskyy/go-easypay/private"
	"github.com/stremovskyy/recorder/file_recorder"
)

func main() {
	recorder := file_recorder.NewFileRecorder("records")
	client := go_easypay.NewClientWithRecorder(recorder)

	merchant := &go_easypay.Merchant{
		Name:            private.MerchantName,
		PartnerKey:      private.PartnerKey,
		ServiceKey:      private.ServiceKey,
		SecretKey:       private.SecretKey,
		SuccessRedirect: private.SuccessRedirect,
		FailRedirect:    private.FailRedirect,
	}

	statusRequest := &go_easypay.Request{
		Merchant: merchant,
		PaymentData: &go_easypay.PaymentData{
			EasypayPaymentID: utils.Ref(int64(private.EasypayPaymentID)),
			PaymentID:        utils.Ref(private.EasypayOrderID),
		},
	}

	client.SetLogLevel(log.LevelDebug)
	statusRequest.SetWebhookURL(utils.Ref(private.WebhookURL))

	statusResponse, err := client.Status(statusRequest)
	if err != nil {
		panic(err)
	}

	if statusResponse.GetError() != nil {
		panic(statusResponse.GetError())
	}

	fmt.Printf("Payment status: %s\n", statusResponse.PaymentState)
}
