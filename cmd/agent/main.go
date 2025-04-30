package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

const (
	baseURL = "https://api-yellow.tazapay.com/v3"
)

type Agent struct {
	client    *http.Client
	authToken string
	reader    *bufio.Reader
}

func NewAgent() (*Agent, error) {
	// Initialize config
	viper.SetConfigName(".tazapay-mcp-server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	// Get API credentials
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	// Generate auth token
	authToken := base64.StdEncoding.EncodeToString([]byte(apiKey + ":" + apiSecret))

	return &Agent{
		client:    &http.Client{},
		authToken: authToken,
		reader:    bufio.NewReader(os.Stdin),
	}, nil
}

func (a *Agent) Run() {
	fmt.Println("Welcome to Tazapay AI Agent!")
	fmt.Println("Available commands:")
	fmt.Println("1. Check balance")
	fmt.Println("2. Create beneficiary")
	fmt.Println("3. Create payout")
	fmt.Println("4. Get FX rates")
	fmt.Println("5. Exit")

	for {
		fmt.Print("\nEnter command (1-5): ")
		input, _ := a.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			a.checkBalance()
		case "2":
			a.createBeneficiary()
		case "3":
			a.createPayout()
		case "4":
			a.getFXRates()
		case "5":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid command. Please enter a number between 1 and 5.")
		}
	}
}

func (a *Agent) checkBalance() {
	fmt.Print("Enter currency (e.g., USD): ")
	currency, _ := a.reader.ReadString('\n')
	currency = strings.TrimSpace(currency)

	url := fmt.Sprintf("%s/balance?currency=%s", baseURL, currency)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+a.authToken)

	resp, err := a.client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %v\n", result["message"])
		return
	}

	balance := result["data"].(map[string]interface{})["balance"].(float64)
	fmt.Printf("Your %s balance is %.2f\n", currency, balance)
}

func (a *Agent) createBeneficiary() {
	fmt.Print("Enter beneficiary name: ")
	name, _ := a.reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter email: ")
	email, _ := a.reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter country: ")
	country, _ := a.reader.ReadString('\n')
	country = strings.TrimSpace(country)

	payload := map[string]interface{}{
		"name":    name,
		"email":   email,
		"country": country,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error creating payload: %v\n", err)
		return
	}

	url := fmt.Sprintf("%s/beneficiaries", baseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+a.authToken)

	resp, err := a.client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %v\n", result["message"])
		return
	}

	fmt.Printf("Beneficiary created successfully! ID: %v\n", result["data"].(map[string]interface{})["id"])
}

func (a *Agent) createPayout() {
	fmt.Print("Enter beneficiary ID: ")
	beneficiaryID, _ := a.reader.ReadString('\n')
	beneficiaryID = strings.TrimSpace(beneficiaryID)

	fmt.Print("Enter amount: ")
	amountStr, _ := a.reader.ReadString('\n')
	amount, _ := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)

	fmt.Print("Enter currency: ")
	currency, _ := a.reader.ReadString('\n')
	currency = strings.TrimSpace(currency)

	payload := map[string]interface{}{
		"beneficiary_id": beneficiaryID,
		"amount":         amount,
		"currency":       currency,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error creating payload: %v\n", err)
		return
	}

	url := fmt.Sprintf("%s/payouts", baseURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+a.authToken)

	resp, err := a.client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %v\n", result["message"])
		return
	}

	fmt.Printf("Payout created successfully! ID: %v\n", result["data"].(map[string]interface{})["id"])
}

func (a *Agent) getFXRates() {
	fmt.Print("Enter source currency (e.g., USD): ")
	from, _ := a.reader.ReadString('\n')
	from = strings.TrimSpace(from)

	fmt.Print("Enter target currency (e.g., EUR): ")
	to, _ := a.reader.ReadString('\n')
	to = strings.TrimSpace(to)

	url := fmt.Sprintf("%s/fx?from=%s&to=%s", baseURL, from, to)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+a.authToken)

	resp, err := a.client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %v\n", result["message"])
		return
	}

	rate := result["data"].(map[string]interface{})["rate"].(float64)
	fmt.Printf("Exchange rate from %s to %s: %.4f\n", from, to, rate)
}

func main() {
	agent, err := NewAgent()
	if err != nil {
		fmt.Printf("Error initializing agent: %v\n", err)
		os.Exit(1)
	}

	agent.Run()
}
