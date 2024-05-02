package easypay

// AppResponse represents the response from the createApp endpoint
type AppResponse struct {
	LogoPath       string `json:"logoPath"`
	HintImagesPath string `json:"hintImagesPath"`
	ApiVersion     string `json:"apiVersion"`
	AppID          string `json:"appId"`
	PageID         string `json:"pageId"`
	Error          *Error `json:"error"`
}
