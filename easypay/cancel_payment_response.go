package easypay

// CancelPaymentResponse is the structure for the response from cancelling a payment
type CancelPaymentResponse struct {
	MerchantKey         string  `json:"merchantKey"`
	TransactionID       int64   `json:"transactionId"`
	RefundTransactionID int64   `json:"refundTransactionId"`
	OrderID             string  `json:"orderId"`
	Amount              float64 `json:"amount"`
	PaymentState        string  `json:"paymentState"`
	Error               *Error  `json:"error"`
}
