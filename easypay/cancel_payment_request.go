package easypay

// CancelPaymentRequest is the structure for canceling a payment
type CancelPaymentRequest struct {
	ServiceKey    string  `json:"serviceKey"`
	OrderID       string  `json:"orderId"`
	TransactionID string  `json:"transactionId"`
	Amount        float64 `json:"amount,omitempty"` // optional for full cancellation
}
