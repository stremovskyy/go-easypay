package easypay

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is the interface for the EasyPay client
type Client interface {
	CreateApp(ctx context.Context) (*AppResponse, error)
	CreatePage(ctx context.Context, appId string) (*Response, error)
	CreateOrder(ctx context.Context, req Request, appId, pageId string) (*Response, error)
	CancelPayment(ctx context.Context, req Request) (*CancelPaymentResponse, error)
	CheckPaymentStatus(ctx context.Context, req Request) (*PaymentStatusResponse, error)
}

// API client holds the configuration for the API client
type apiClient struct {
	httpClient *http.Client
	baseURL    string
	partnerKey string
	locale     string
	secretKey  string
	appID      string
	pageID     string
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

type Config struct {
	BaseURL    string
	PartnerKey string
	Locale     string
	SecretKey  string
	AppID      string
	PageID     string
}

func NewClient(config Config, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = newHTTPClient()
	}
	return &apiClient{
		httpClient: httpClient,
		baseURL:    config.BaseURL,
		partnerKey: config.PartnerKey,
		locale:     config.Locale,
		secretKey:  config.SecretKey,
		appID:      config.AppID,
		pageID:     config.PageID,
	}
}

// CreateApp handles the creation of an application session
func (c *apiClient) CreateApp(ctx context.Context) (*AppResponse, error) {
	url := fmt.Sprintf("%s/api/system/createApp", c.baseURL)
	reqBody, err := json.Marshal(
		map[string]string{
			"partnerKey": c.partnerKey,
			"locale":     c.locale,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	var appResponse AppResponse
	if err := json.NewDecoder(response.Body).Decode(&appResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if appResponse.Error != nil {
		return nil, fmt.Errorf("API error: %s", appResponse.Error)
	}

	return &appResponse, nil
}

// CreatePage creates a new session for the user and returns the PageId.
func (c *apiClient) CreatePage(ctx context.Context, appId string) (*Response, error) {
	url := fmt.Sprintf("%s/api/system/createPage", c.baseURL)
	reqBody, err := json.Marshal(
		map[string]string{
			"partnerKey": c.partnerKey,
			"locale":     c.locale,
			"AppId":      appId,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	var pageResponse Response
	if err := json.NewDecoder(response.Body).Decode(&pageResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if pageResponse.Error != nil {
		return nil, fmt.Errorf("API error: %s", pageResponse.Error)
	}

	return &pageResponse, nil
}

// CreateOrder sends a request to create an order
func (c *apiClient) CreateOrder(ctx context.Context, req Request, appId, pageId string) (*Response, error) {
	url := fmt.Sprintf("%s/api/merchant/createOrder", c.baseURL)
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	signature := generateSignature(requestBody, c.secretKey)

	// Add necessary headers
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("AppId", c.appID)
	httpRequest.Header.Add("PageId", c.pageID)
	httpRequest.Header.Add("PartnerKey", c.partnerKey)
	httpRequest.Header.Add("locale", c.locale)
	httpRequest.Header.Add("Sign", signature)

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	var orderResponse Response
	if err := json.NewDecoder(response.Body).Decode(&orderResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &orderResponse, nil
}

// CheckPaymentStatus checks the status of a payment by transaction or order ID
func (c *apiClient) CheckPaymentStatus(ctx context.Context, req Request) (*PaymentStatusResponse, error) {
	url := fmt.Sprintf("%s/api/merchant/orderState", c.baseURL)
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	signature := generateSignature(requestBody, c.secretKey)

	httpRequest, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("AppId", c.appID)
	httpRequest.Header.Add("PageId", c.pageID)
	httpRequest.Header.Add("PartnerKey", c.partnerKey)
	httpRequest.Header.Add("locale", c.locale)
	httpRequest.Header.Add("Sign", signature)

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	var statusResponse PaymentStatusResponse
	if err := json.NewDecoder(response.Body).Decode(&statusResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &statusResponse, nil
}

// CancelPayment cancels an accepted payment
func (c *apiClient) CancelPayment(ctx context.Context, req Request) (*CancelPaymentResponse, error) {
	url := fmt.Sprintf("%s/api/merchant/cancelOrder", c.baseURL)
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	signature := generateSignature(requestBody, c.secretKey)

	httpRequest, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("AppId", c.appID)
	httpRequest.Header.Add("PageId", c.pageID)
	httpRequest.Header.Add("PartnerKey", c.partnerKey)
	httpRequest.Header.Add("locale", c.locale)
	httpRequest.Header.Add("Sign", signature)

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", response.StatusCode)
	}

	var cancelResponse CancelPaymentResponse
	if err := json.NewDecoder(response.Body).Decode(&cancelResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &cancelResponse, nil
}

// generateSignature creates a HMAC-SHA256 signature using the specified secret key and data.
func generateSignature(data []byte, secretKey string) string {
	hmac := hmac.New(sha256.New, []byte(secretKey))
	hmac.Write(data)
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}
