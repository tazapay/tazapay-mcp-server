package tazapay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://api.tazapay.com/v3"
)

// Response structs
type BalanceResponse struct {
	Data struct {
		Balance  float64 `json:"balance"`
		Currency string  `json:"currency"`
	} `json:"data"`
}

type Beneficiary struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Country string `json:"country"`
}

type BeneficiaryResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type Payout struct {
	BeneficiaryID string  `json:"beneficiary_id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
}

type PayoutResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type ExchangeRateResponse struct {
	Data struct {
		Rate float64 `json:"rate"`
	} `json:"data"`
}

// Client represents a Tazapay API client
type Client struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

// NewClient creates a new Tazapay API client
func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getAuthHeader returns the Authorization header value
func (c *Client) getAuthHeader() string {
	auth := fmt.Sprintf("%s:%s", c.apiKey, c.apiSecret)
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s%s", baseURL, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader())
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("Making request to: %s\n", url) // Debug log

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Log the response status and body for debugging
	fmt.Printf("Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetBalance retrieves the account balance
func (c *Client) GetBalance() (*BalanceResponse, error) {
	respBody, err := c.doRequest("GET", "/balance", nil)
	if err != nil {
		return nil, fmt.Errorf("error checking balance: %w", err)
	}

	var balance BalanceResponse
	if err := json.Unmarshal(respBody, &balance); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &balance, nil
}

// CreateBeneficiary creates a new beneficiary
func (c *Client) CreateBeneficiary(beneficiary *Beneficiary) (*BeneficiaryResponse, error) {
	respBody, err := c.doRequest("POST", "/beneficiaries", beneficiary)
	if err != nil {
		return nil, fmt.Errorf("error creating beneficiary: %w", err)
	}

	var response BeneficiaryResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// CreatePayout creates a new payout
func (c *Client) CreatePayout(payout *Payout) (*PayoutResponse, error) {
	respBody, err := c.doRequest("POST", "/payouts", payout)
	if err != nil {
		return nil, fmt.Errorf("error creating payout: %w", err)
	}

	var response PayoutResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetExchangeRate gets the exchange rate between two currencies
func (c *Client) GetExchangeRate(fromCurrency, toCurrency string) (*ExchangeRateResponse, error) {
	path := fmt.Sprintf("/exchange-rate?from=%s&to=%s", fromCurrency, toCurrency)
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting exchange rate: %w", err)
	}

	var response ExchangeRateResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}
