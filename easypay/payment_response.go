package easypay

type OrderResponse struct {
	PaymentState            string                  `json:"paymentState"`
	ActionType              string                  `json:"actionType"`
	Action                  string                  `json:"action"`
	AlternativeRedirectUrl  string                  `json:"alternativeRedirectUrl,omitempty"`
	TransactionId           int64                   `json:"transactionId,omitempty"`
	RetrievalReferenceNo    string                  `json:"retrievalReferenceNo,omitempty"`
	ForwardUrl              string                  `json:"forwardUrl,omitempty"`
	ActionContent           string                  `json:"actionContent,omitempty"`
	Error                   *Error                  `json:"error"`
	ResponseItems           ResponseItems           `json:"responseItems,omitempty"`
	PaymentInstrumentsTypes []PaymentInstrumentType `json:"paymentInstrumentsTypes,omitempty"`
}

// ResponseItems could contain various details specific to the transaction
type ResponseItems struct {
	SessionId         string `json:"sessionId,omitempty"`
	MerchantOperation string `json:"merchantOperation,omitempty"`
	Operation         string `json:"operation,omitempty"`
	BankingDetails    string `json:"bankingDetails,omitempty"`
}

// PaymentInstrumentType provides details about available payment methods and conditions
type PaymentInstrumentType struct {
	InstrumentType         string                         `json:"instrumentType"`
	Commission             float64                        `json:"commission"`
	AmountMin              float64                        `json:"amountMin"`
	AmountMax              float64                        `json:"amountMax"`
	UserPaymentInstruments []UserPaymentInstrumentDetails `json:"userPaymentInstruments"`
}

// UserPaymentInstrumentDetails provides details about user-specific payment instruments, like saved cards
type UserPaymentInstrumentDetails struct {
	InstrumentId      int64                  `json:"instrumentId"`
	InstrumentType    string                 `json:"instrumentType"`
	InstrumentValue   string                 `json:"instrumentValue,omitempty"`
	Alias             string                 `json:"alias,omitempty"`
	Commission        float64                `json:"commission"`
	LoyaltyCommission float64                `json:"loyaltyCommission,omitempty"`
	ActionsKeys       []string               `json:"actionsKeys,omitempty"`
	PriorityIndex     int                    `json:"priorityIndex"`
	AdditionalParams  map[string]interface{} `json:"additionalParams,omitempty"`
}

// Error structure to encapsulate API error details
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
