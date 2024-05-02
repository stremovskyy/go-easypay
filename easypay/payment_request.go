package easypay

// OrderRequest represents the request body for creating an order
type OrderRequest struct {
	UserInfo              UserInfo              `json:"userInfo"`
	Order                 Order                 `json:"order"`
	URLs                  URLs                  `json:"urls"`
	BankingDetailsID      string                `json:"bankingDetailsId,omitempty"`
	BankingDetails        BankingDetails        `json:"bankingDetails,omitempty"`
	Recurrent             Recurrent             `json:"reccurent,omitempty"`
	Splitting             Splitting             `json:"splitting,omitempty"`
	UserPaymentInstrument UserPaymentInstrument `json:"userPaymentInstrument"`
	PartnerInfo           PartnerInfo           `json:"partnerInfo,omitempty"`
	BrowserInfo           BrowserInfo           `json:"browserInfo,omitempty"`
}

// UserInfo holds user-specific information
type UserInfo struct {
	Phone string `json:"phone"`
}

// Order contains details about the order
type Order struct {
	ServiceKey       string            `json:"serviceKey"`
	OrderID          string            `json:"orderId"`
	Description      string            `json:"description"`
	Amount           float64           `json:"amount"`
	PaymentOperation string            `json:"paymentOperation"`
	AdditionalItems  map[string]string `json:"additionalItems,omitempty"`
	Expire           string            `json:"expire"`
	IsOneTimePay     bool              `json:"isOneTimePay"`
	Fields           []Field           `json:"fields,omitempty"`
}

// Field represents a single custom field in the order
type Field struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
	FieldKey   string `json:"fieldKey,omitempty"`
}

// URLs struct to hold success and failed URLs
type URLs struct {
	Success string `json:"success"`
	Failed  string `json:"failed"`
}

// BankingDetails holds banking information
type BankingDetails struct {
	Payee     BankingParty `json:"payee"`
	Payer     Payer        `json:"payer"`
	Narrative Narrative    `json:"narrative"`
}

// BankingParty details about the payee or payer
type BankingParty struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bank Bank   `json:"bank"`
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
	Name    string `json:"name"`
	MFO     string `json:"mfo"`
	Account string `json:"account"`
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
	InstrumentType string `json:"instrumentType"`
	CardGuid       string `json:"cardGuid,omitempty"`
	Pan            string `json:"pan,omitempty"`
	Expire         string `json:"expire,omitempty"`
	CVV            string `json:"cvv,omitempty"`
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
