package easypay

import (
	"encoding/json"
	"fmt"
	"time"
)

type Webhook struct {
	PartnerKey      string           `json:"partnerKey"`
	Phone           string           `json:"phone"`
	CardGuid        string           `json:"cardGuid"`
	Pan             string           `json:"pan"`
	Expire          string           `json:"expire"`
	DatePost        string           `json:"datePost"`
	CodeType        string           `json:"codeType"`
	CardLabel       *string          `json:"cardLabel"`
	ExistingTokens  interface{}      `json:"existingTokens"`
	Action          string           `json:"action"`
	MerchantId      int              `json:"merchant_id"`
	OrderId         string           `json:"order_id"`
	Version         string           `json:"version"`
	Date            time.Time        `json:"date"`
	Details         *WebhookDetails  `json:"details"`
	Additionalitems *Additionalitems `json:"additionalitems"`
}

type WebhookDetails struct {
	Amount      int         `json:"amount"`
	Desc        string      `json:"desc"`
	PaymentId   int         `json:"payment_id"`
	RecurrentId interface{} `json:"recurrent_id"`
}

type Additionalitems struct {
	CardPan                     string `json:"Card.Pan"`
	MerchantKey                 string `json:"MerchantKey"`
	CommissionClient            string `json:"Commission.Client"`
	BankName                    string `json:"BankName"`
	MerchantCardTokenizationKey string `json:"Merchant.CardTokenization.Key"`
	MerchantIsOneTimePay        string `json:"Merchant.IsOneTimePay"`
	MerchantOrderId             string `json:"Merchant.OrderId"`
	MerchantUrlNotify           string `json:"Merchant.UrlNotify"`
	AcquirerMerchantId          string `json:"Acquirer.MerchantId"`
	AcquirerTerminalId          string `json:"Acquirer.TerminalId"`
	CardBrandType               string `json:"Card.BrandType"`
	Rrn                         string `json:"Rrn"`
	AuthCode                    string `json:"AuthCode"`
}

func ParseWebhook(body []byte) (*Webhook, error) {
	w := &Webhook{}
	err := json.Unmarshal(body, w)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal webhook: %v", err)
	}

	return w, nil
}
