package easypay

import (
	"encoding/json"
	"fmt"
	"html"
	"strings"
)

type Status string

const (
	StatusConfirmed         Status = "Confirmed"
	StatusRejected          Status = "Rejected"
	StatusWaitVerify        Status = "WaitVerify"
	StatusPending           Status = "Pending"
	StatusRefunded          Status = "Refunded"
	StatusWaitConfirm       Status = "WaitConfirm"
	StatusCancelingAccepted Status = "accepted"
	StatusCancelingDeclined Status = "declined"
)

type Response struct {
	PaymentState            Status                  `json:"paymentState"`
	ActionType              string                  `json:"actionType"`
	Action                  string                  `json:"action"`
	AlternativeRedirectUrl  string                  `json:"alternativeRedirectUrl,omitempty"`
	TransactionId           int64                   `json:"transactionId,omitempty"`
	RetrievalReferenceNo    string                  `json:"retrievalReferenceNo,omitempty"`
	ForwardUrl              string                  `json:"forwardUrl,omitempty"`
	ActionContent           string                  `json:"actionContent,omitempty"`
	Error                   *Error                  `json:"error,omitempty"`
	ResponseItems           ResponseItems           `json:"responseItems,omitempty"`
	PaymentInstrumentsTypes []PaymentInstrumentType `json:"paymentInstrumentsTypes,omitempty"`
	LogoPath                *string                 `json:"logoPath,omitempty"`
	HintImagesPath          *string                 `json:"hintImagesPath,omitempty"`
	ApiVersion              *string                 `json:"apiVersion,omitempty"`
	AppId                   *string                 `json:"appId,omitempty"`
	PageId                  *string                 `json:"pageId,omitempty"`
	RequestedSessionId      string                  `json:"requestedSessionId"`
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
	ErrorCode       *string `json:"errorCode"`
	ClientErrorCode *string `json:"clientErrorCode"`
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	ErrorMessage    *string `json:"errorMessage"`
	FieldErrors     []struct {
		FieldName    string      `json:"fieldName"`
		ErrorCode    interface{} `json:"errorCode"`
		ErrorMessage string      `json:"errorMessage"`
	} `json:"fieldErrors"`
}

// GetError returns a constructed error based on the response error details
// CustomError is a custom type for formatting errors from the Response struct.
type CustomError struct {
	Resp *Response
}

func (ce *CustomError) Error() string {
	if ce.Resp.Error == nil {
		return "no error present"
	}
	e := ce.Resp.Error
	var sb strings.Builder

	sb.WriteString("API Error: ")
	if e.ErrorCode != nil {
		sb.WriteString(fmt.Sprintf("Code: %s, ", *e.ErrorCode))
	}
	if e.ClientErrorCode != nil {
		sb.WriteString(fmt.Sprintf("Client Code: %s, ", *e.ClientErrorCode))
	}
	if e.Title != nil {
		sb.WriteString(fmt.Sprintf("Title: %s, ", html.UnescapeString(*e.Title)))
	}
	if e.Description != nil {
		sb.WriteString(fmt.Sprintf("Description: %s, ", html.UnescapeString(*e.Description)))
	}
	if e.ErrorMessage != nil {
		sb.WriteString(fmt.Sprintf("Message: %s, ", html.UnescapeString(*e.ErrorMessage)))
	}
	if len(e.FieldErrors) > 0 {
		sb.WriteString("Field Errors: ")
		for _, fe := range e.FieldErrors {
			sb.WriteString(fmt.Sprintf("[%s: Code: %v, Message: %s], ", fe.FieldName, fe.ErrorCode, html.UnescapeString(fe.ErrorMessage)))
		}
	}
	return strings.TrimRight(sb.String(), ", ")
}

func (r *Response) GetError() error {
	if r.Error != nil {
		return &CustomError{Resp: r}
	}
	return nil
}

func (r *Response) App() *App {
	return NewApp(r.LogoPath, r.ApiVersion, r.AppId)
}

func UnmarshalJSONResponse(data []byte) (*Response, error) {
	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}
	return &resp, nil
}
