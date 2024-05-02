package easypay

// PaymentDetail holds information about individual payments in a split payment or detailed transaction list
type PaymentDetail struct {
	MerchantKey         string  `json:"merchantKey"`
	TransactionID       int64   `json:"transactionId"`
	OrderID             string  `json:"orderId"`
	Amount              float64 `json:"amount"`
	PaymentState        string  `json:"paymentState"`
	RefundTransactionID int64   `json:"refundTransactionId,omitempty"`
	Date                string  `json:"date,omitempty"`
	Error               *Error  `json:"error,omitempty"`
}

// PaymentStatusResponse is the structure for the payment status response
type PaymentStatusResponse struct {
	MerchantKey         string          `json:"merchantKey"`
	TransactionID       int64           `json:"transactionId"`
	OrderID             string          `json:"orderId"`
	Amount              float64         `json:"amount"`
	PaymentState        string          `json:"paymentState"`
	RefundTransactionID int64           `json:"refundTransactionId,omitempty"`
	Error               *Error          `json:"error"`
	PaymentsList        []PaymentDetail `json:"paymentsList,omitempty"`
}
