package tazapay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL = "https://api-yellow.tazapay.com/v3"
)

// Client represents a Tazapay API client
type Client struct {
	httpClient *http.Client
	authToken  string
}

// NewClient creates a new Tazapay client
func NewClient(apiKey, apiSecret string) *Client {
	authToken := base64.StdEncoding.EncodeToString([]byte(apiKey + ":" + apiSecret))

	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		authToken: authToken,
	}
}

// doRequest performs an HTTP request with the given method, path, and body
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+c.authToken)

	return c.httpClient.Do(req)
}

// GetBalance retrieves the balance for a given currency
func (c *Client) GetBalance(currency string) (float64, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/balance?currency=%s", currency), nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Balance float64 `json:"balance"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	return result.Data.Balance, nil
}

// CreateBeneficiary creates a new beneficiary
func (c *Client) CreateBeneficiary(name, email, country string) (string, error) {
	payload := map[string]string{
		"name":    name,
		"email":   email,
		"country": country,
	}

	resp, err := c.doRequest("POST", "/beneficiaries", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return result.Data.ID, nil
}

// CreatePayout creates a new payout
func (c *Client) CreatePayout(beneficiaryID string, amount float64, currency string) (string, error) {
	payload := map[string]interface{}{
		"beneficiary_id": beneficiaryID,
		"amount":         amount,
		"currency":       currency,
	}

	resp, err := c.doRequest("POST", "/payouts", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	return result.Data.ID, nil
}

// GetFXRates gets the exchange rate between two currencies
func (c *Client) GetFXRates(from, to string) (float64, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/fx?from=%s&to=%s", from, to), nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Rate float64 `json:"rate"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	return result.Data.Rate, nil
}
