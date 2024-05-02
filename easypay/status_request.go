package easypay

// PaymentStatusRequest is the structure for the payment status check request
type PaymentStatusRequest struct {
	ServiceKey    string `json:"serviceKey"`
	OrderID       string `json:"orderId,omitempty"`
	TransactionID string `json:"transactionId,omitempty"`
}
