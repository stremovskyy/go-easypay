package main

import (
	"github.com/google/uuid"

	go_easypay "github.com/stremovskyy/go-easypay"
	"github.com/stremovskyy/go-easypay/internal/utils"
	"github.com/stremovskyy/go-easypay/log"
	"github.com/stremovskyy/go-easypay/private"
)

func main() {
	client := go_easypay.NewDefaultClient()

	merchant := &go_easypay.Merchant{
		Name:            private.MerchantName,
		PartnerKey:      private.PartnerKey,
		ServiceKey:      private.ServiceKey,
		SecretKey:       private.SecretKey,
		SuccessRedirect: private.SuccessRedirect,
		FailRedirect:    private.FailRedirect,
	}

	uuidString := uuid.New().String()

	VerificationRequest := &go_easypay.Request{
		Merchant: merchant,
		PaymentData: &go_easypay.PaymentData{
			EasypayPaymentID: utils.Ref(int64(private.EasypayPaymentID)),
			PaymentID:        utils.Ref(uuidString),
			OrderID:          uuidString,
			Description:      "Verification payment: " + uuidString,
		},
		PersonalData: &go_easypay.PersonalData{
			UserID:    utils.Ref(123),
			FirstName: utils.Ref("John"),
			LastName:  utils.Ref("Doe"),
			TaxID:     utils.Ref("1234567890"),
		},
	}

	client.SetLogLevel(log.LevelDebug)
	VerificationRequest.SetWebhookURL(utils.Ref(private.WebhookURL))

	tokenURL, err := client.VerificationLink(VerificationRequest)
	if err != nil {
		panic(err)
	}

	println(tokenURL.String())
}
