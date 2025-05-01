package tazapay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	baseURL = "https://api-yellow.tazapay.com/v3"
)

// Response structs
type BalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Available []struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"available"`
		Object    string `json:"object"`
		UpdatedAt string `json:"updated_at"`
	} `json:"data"`
}

type BankDetails struct {
	Country   string `json:"country"`
	Currency  string `json:"currency"`
	BankCodes struct {
		SwiftCode string `json:"swift_code"`
	} `json:"bank_codes"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
}

type DestinationDetails struct {
	Type string      `json:"type"`
	Bank BankDetails `json:"bank"`
}

type Phone struct {
	Number      string `json:"number"`
	CallingCode string `json:"calling_code"`
}

type Beneficiary struct {
	AccountID                    string             `json:"account_id"`
	Name                         string             `json:"name"`
	Type                         string             `json:"type"`
	Email                        string             `json:"email"`
	NationalIdentificationNumber string             `json:"national_identification_number"`
	TaxID                        string             `json:"tax_id"`
	DestinationDetails           DestinationDetails `json:"destination_details"`
	Phone                        Phone              `json:"phone"`
}

type BeneficiaryResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID string `json:"id"`
	} `json:"data"`
}

// BeneficiaryDetails represents the beneficiary details for a payout
type BeneficiaryDetails struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	DestinationDetails struct {
		Type string      `json:"type"`
		Bank BankDetails `json:"bank"`
	} `json:"destination_details"`
}

// Payout represents a payout request
type Payout struct {
	BeneficiaryDetails  BeneficiaryDetails `json:"beneficiary_details"`
	Purpose             string             `json:"purpose"`
	Amount              int64              `json:"amount"`
	Currency            string             `json:"currency"`
	Beneficiary         string             `json:"beneficiary"`
	HoldingCurrency     string             `json:"holding_currency"`
	Type                string             `json:"type"`
	ChargeType          string             `json:"charge_type"`
	StatementDescriptor string             `json:"statement_descriptor"`
}

// PayoutResponse represents the response from creating a payout
type PayoutResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID          string `json:"id"`
		Status      string `json:"status"`
		Amount      int64  `json:"amount"`
		Currency    string `json:"currency"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Purpose     string `json:"purpose"`
		Type        string `json:"type"`
		ChargeType  string `json:"charge_type"`
		Beneficiary string `json:"beneficiary"`
	} `json:"data"`
}

type ExchangeRateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Amount          int64   `json:"amount"`
		ConvertedAmount int64   `json:"converted_amount"`
		ExchangeRate    float64 `json:"exchange_rate"`
		FinalCurrency   string  `json:"final_currency"`
		InitialCurrency string  `json:"initial_currency"`
		Timestamp       string  `json:"timestamp"`
	} `json:"data"`
}

// PaymentRequest represents a payment creation request
type PaymentRequest struct {
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	InvoiceCurrency string  `json:"invoice_currency"`
	Description     string  `json:"description"`
	TransactionDesc string  `json:"transaction_description"`
	SuccessURL      string  `json:"success_url"`
	CancelURL       string  `json:"cancel_url"`
	CustomerEmail   string  `json:"customer_email"`
	CustomerName    string  `json:"customer_name"`
	CustomerDetails struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Phone struct {
			Number      string `json:"number"`
			CallingCode string `json:"calling_code"`
		} `json:"phone"`
		Address string `json:"address"`
		Country string `json:"country"`
	} `json:"customer_details"`
	PaymentMethods []string `json:"payment_methods,omitempty"`
}

// PaymentResponse represents the response from creating a payment
type PaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID        string  `json:"id"`
		Status    string  `json:"status"`
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		URL       string  `json:"url"`
		CreatedAt string  `json:"created_at"`
	} `json:"data"`
}

// Client represents a Tazapay API client
type Client struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Tazapay API client
func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
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

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers to match the working curl request
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("expires", "0")
	req.Header.Set("tz-account-id", "acc_d00eij4qqkhc9e5ats4g")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.apiKey, c.apiSecret)))))

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

	// Handle specific HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK:
		return respBody, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication failed: invalid API key or secret")
	case http.StatusForbidden:
		return nil, fmt.Errorf("access forbidden: %s", string(respBody))
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("rate limit exceeded: %s", string(respBody))
	case http.StatusBadRequest:
		return nil, fmt.Errorf("bad request: %s", string(respBody))
	default:
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
}

// GetBalance retrieves the account balance
func (c *Client) GetBalance() (*BalanceResponse, error) {
	respBody, err := c.doRequest("GET", "/balance", nil)
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %w", err)
	}

	var balanceResp BalanceResponse
	if err := json.Unmarshal(respBody, &balanceResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &balanceResp, nil
}

// CreateBeneficiary creates a new beneficiary
func (c *Client) CreateBeneficiary(beneficiary *Beneficiary) (*BeneficiaryResponse, error) {
	respBody, err := c.doRequest("POST", "/beneficiary", beneficiary)
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
	respBody, err := c.doRequest("POST", "/payout", payout)
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
func (c *Client) GetExchangeRate(fromCurrency, toCurrency string, amount float64) (*ExchangeRateResponse, error) {
	// Convert amount to integer (cents)
	amountInCents := int64(amount * 100)
	// Convert currency codes to uppercase
	fromCurrency = strings.ToUpper(fromCurrency)
	toCurrency = strings.ToUpper(toCurrency)
	path := fmt.Sprintf("/fx/payout?initial_currency=%s&final_currency=%s&amount=%d", fromCurrency, toCurrency, amountInCents)
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

// GetAuthToken returns the base64 encoded authentication token
func (c *Client) GetAuthToken() string {
	auth := fmt.Sprintf("%s:%s", c.apiKey, c.apiSecret)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// CreatePayment creates a new payment
func (c *Client) CreatePayment(payment *PaymentRequest) (*PaymentResponse, error) {
	respBody, err := c.doRequest("POST", "/checkout", payment)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	var response PaymentResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &response, nil
}

// GetLogs retrieves recent logs from the Tazapay API
func (c *Client) GetLogs(count int) ([]interface{}, error) {
	path := fmt.Sprintf("/logs?count=%d", count)
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting logs: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	logs, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	return logs, nil
}
