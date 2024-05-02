package easypay

// PageResponse represents the response from the createPage endpoint
type PageResponse struct {
	LogoPath           string `json:"logoPath"`
	HintImagesPath     string `json:"hintImagesPath"`
	ApiVersion         string `json:"apiVersion"`
	AppID              string `json:"appId"`
	PageID             string `json:"pageId"`
	RequestedSessionId string `json:"requestedSessionId"`
	Error              *Error `json:"error"`
}
