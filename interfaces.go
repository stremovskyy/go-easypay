package go_easypay

import (
	"context"
	"net/url"

	"github.com/stremovskyy/go-easypay/easypay"
	"github.com/stremovskyy/go-easypay/log"
)

type Easypay interface {
	VerificationLink(request *Request) (*url.URL, error)
	Status(request *Request) (*easypay.Response, error)
	PaymentURL(invoiceRequest *Request) (*easypay.Response, error)
	Payment(invoiceRequest *Request) (*easypay.Response, error)
	Hold(invoiceRequest *Request) (*easypay.Response, error)
	Capture(invoiceRequest *Request) (*easypay.Response, error)
	Refund(invoiceRequest *Request) (*easypay.Response, error)
	Credit(invoiceRequest *Request) (*easypay.Response, error)
	SetLogLevel(levelDebug log.Level)

	GetRecordedExchange(ctx context.Context, requestID string) (*easypay.RecordedExchange, error)
	GetExchangesByOrderID(ctx context.Context, orderID string) ([]*easypay.RecordedExchange, error)
	GetExchangesByTransactionID(ctx context.Context, transactionID string) ([]*easypay.RecordedExchange, error)
}
