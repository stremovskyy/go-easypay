package easypay

// Request represents the request body for creating an order
type Request struct {
	UserInfo              *UserInfo              `json:"userInfo,omitempty"`
	Order                 *Order                 `json:"order,omitempty"`
	URLs                  *URLs                  `json:"urls,omitempty"`
	BankingDetailsID      *string                `json:"bankingDetailsId,omitempty"`
	BankingDetails        *BankingDetails        `json:"bankingDetails,omitempty"`
	Recurrent             *Recurrent             `json:"reccurent,omitempty"`
	Splitting             *Splitting             `json:"splitting,omitempty"`
	UserPaymentInstrument *UserPaymentInstrument `json:"userPaymentInstrument,omitempty"`
	PartnerInfo           *PartnerInfo           `json:"partnerInfo,omitempty"`
	BrowserInfo           *BrowserInfo           `json:"browserInfo,omitempty"`
	ServiceKey            *string                `json:"serviceKey,omitempty"`
	OrderID               *string                `json:"orderId,omitempty"`
	TransactionID         *string                `json:"transactionId,omitempty"`
	Amount                *float64               `json:"amount,omitempty"` // optional for full cancellation
	Phone                 *string                `json:"phone,omitempty"`

	Url       string            `json:"-"`
	Headers   map[string]string `json:"-"`
	SecretKey string            `json:"-"`
}

// UserInfo holds user-specific information
type UserInfo struct {
	Phone *string `json:"phone"`
}

// Order contains details about the order
type Order struct {
	ServiceKey       *string            `json:"serviceKey,omitempty"`
	OrderID          *string            `json:"orderId,omitempty"`
	Description      *string            `json:"description,omitempty"`
	Amount           *float64           `json:"amount,omitempty"`
	PaymentOperation *string            `json:"paymentOperation,omitempty"`
	AdditionalItems  *map[string]string `json:"additionalItems,omitempty"`
	Expire           *string            `json:"expire,omitempty"`
	IsOneTimePay     *bool              `json:"isOneTimePay,omitempty"`
	Fields           *[]Field           `json:"fields,omitempty"`
}

// Field represents a single custom field in the order
type Field struct {
	FieldName  string `json:"fieldName,omitempty"`
	FieldValue string `json:"fieldValue,omitempty"`
	FieldKey   string `json:"fieldKey,omitempty"`
}

// URLs struct to hold success and failed URLs
type URLs struct {
	Success *string `json:"success,omitempty"`
	Failed  *string `json:"failed,omitempty"`
	Notify  *string `json:"notify,omitempty"`
}

// BankingDetails holds banking information
type BankingDetails struct {
	Payee     *Payee     `json:"payee,omitempty"`
	Payer     *Payer     `json:"payer,omitempty"`
	Narrative *Narrative `json:"narrative,omitempty"`
}

// Payee details about the payee or payer
type Payee struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Bank *Bank  `json:"bank,omitempty"`
}

// Payer details about the payer
type Payer struct {
	Name string `json:"name"`
}

// Narrative for the transaction
type Narrative struct {
	Name string `json:"name"`
}

// Bank information
type Bank struct {
	Name    string `json:"name,omitempty"`
	MFO     string `json:"mfo,omitempty"`
	Account string `json:"account,omitempty"`
}

// Recurrent details for recurring payments
type Recurrent struct {
	CronRule   string              `json:"cronRule"`
	DateExpire string              `json:"dateExpire"`
	DateRun    string              `json:"dateRun"`
	Properties RecurrentProperties `json:"properties"`
}

// RecurrentProperties additional settings for recurrence
type RecurrentProperties struct {
	FailedCount int     `json:"failedCount"`
	FailedRule  string  `json:"failedRule"`
	Amount      float64 `json:"amount"`
	UrlNotify   string  `json:"UrlNotify"`
}

// Splitting information for dividing the payment
type Splitting struct {
	Items []SplitItem `json:"items"`
}

// SplitItem details for each split
type SplitItem struct {
	ServiceKey       string         `json:"serviceKey"`
	OrderID          string         `json:"orderId"`
	BankingDetailsID string         `json:"bankingDetailsId"`
	BankingDetails   BankingDetails `json:"bankingDetails"`
	Unit             string         `json:"unit"`
	Value            float64        `json:"value"`
	WithCommission   bool           `json:"withCommission"`
}

// UserPaymentInstrument payment method details
type UserPaymentInstrument struct {
	InstrumentType *string `json:"instrumentType"`
	CardGuid       *string `json:"cardGuid,omitempty"`
	Pan            *string `json:"pan,omitempty"`
	Expire         *string `json:"expire,omitempty"`
	CVV            *string `json:"cvv,omitempty"`
}

// PartnerInfo information about the partner
type PartnerInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
}

// BrowserInfo details for 3DSecure and browser-based authentication
type BrowserInfo struct {
	ColorDepth        string `json:"colorDepth"`
	ScreenHeight      string `json:"screenHeight"`
	ScreenWidth       string `json:"screenWidth"`
	Language          string `json:"language"`
	JavaEnabled       string `json:"javaEnabled"`
	JavaScriptEnabled string `json:"javascriptEnabled"`
	TimeZone          string `json:"timeZone"`
}
