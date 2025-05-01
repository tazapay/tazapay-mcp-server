package chat

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"tazapay-mcp-server/internal/tazapay"
)

// ChatInterface handles user interactions and API calls
type ChatInterface struct {
	client *tazapay.Client
	reader *bufio.Reader
}

// NewChatInterface creates a new chat interface
func NewChatInterface(client *tazapay.Client) *ChatInterface {
	return &ChatInterface{
		client: client,
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run starts the chat interface
func (c *ChatInterface) Run() {
	fmt.Println("Welcome to Tazapay Chat Interface!")
	fmt.Println("Available commands:")
	fmt.Println("1. Check balance")
	fmt.Println("2. Create beneficiary")
	fmt.Println("3. Create payout")
	fmt.Println("4. Get FX rates")
	fmt.Println("5. Exit")

	for {
		fmt.Print("\nEnter command (1-5): ")
		input, _ := c.reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			c.handleCheckBalance()
		case "2":
			c.handleCreateBeneficiary()
		case "3":
			c.handleCreatePayout()
		case "4":
			c.handleGetFXRates()
		case "5":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid command. Please enter a number between 1 and 5.")
		}
	}
}

func (c *ChatInterface) handleCheckBalance() {
	balance, err := c.client.GetBalance()
	if err != nil {
		fmt.Printf("Error checking balance: %v\n", err)
		return
	}

	fmt.Println("Your balances:")
	for _, bal := range balance.Data.Available {
		fmt.Printf("- %s %s\n", bal.Amount, bal.Currency)
	}
	fmt.Printf("Last updated: %s\n", balance.Data.UpdatedAt)
}

func (c *ChatInterface) handleCreateBeneficiary() {
	fmt.Print("Enter beneficiary name: ")
	name, _ := c.reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter email: ")
	email, _ := c.reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter type (individual/business): ")
	beneficiaryType, _ := c.reader.ReadString('\n')
	beneficiaryType = strings.TrimSpace(beneficiaryType)

	fmt.Print("Enter bank country code (e.g., VN): ")
	bankCountry, _ := c.reader.ReadString('\n')
	bankCountry = strings.TrimSpace(bankCountry)

	fmt.Print("Enter bank currency (e.g., VND): ")
	bankCurrency, _ := c.reader.ReadString('\n')
	bankCurrency = strings.TrimSpace(bankCurrency)

	fmt.Print("Enter bank name: ")
	bankName, _ := c.reader.ReadString('\n')
	bankName = strings.TrimSpace(bankName)

	fmt.Print("Enter account number: ")
	accountNumber, _ := c.reader.ReadString('\n')
	accountNumber = strings.TrimSpace(accountNumber)

	fmt.Print("Enter SWIFT code: ")
	swiftCode, _ := c.reader.ReadString('\n')
	swiftCode = strings.TrimSpace(swiftCode)

	beneficiary := &tazapay.Beneficiary{
		AccountID: "acc_d00eij4qqkhc9e5ats4g",
		Name:      name,
		Type:      beneficiaryType,
		Email:     email,
		DestinationDetails: tazapay.DestinationDetails{
			Type: "bank",
			Bank: tazapay.BankDetails{
				Country:       bankCountry,
				Currency:      bankCurrency,
				BankName:      bankName,
				AccountNumber: accountNumber,
				BankCodes: struct {
					SwiftCode string `json:"swift_code"`
				}{
					SwiftCode: swiftCode,
				},
			},
		},
		Phone: tazapay.Phone{
			CallingCode: "84", // Default to Vietnam
		},
	}

	response, err := c.client.CreateBeneficiary(beneficiary)
	if err != nil {
		fmt.Printf("Error creating beneficiary: %v\n", err)
		return
	}

	fmt.Printf("Beneficiary created successfully! ID: %s\n", response.Data.ID)
}

func (c *ChatInterface) handleCreatePayout() {
	fmt.Println("\nCreate Payout")
	fmt.Println("-------------")

	fmt.Print("Enter beneficiary ID: ")
	beneficiaryID := c.readInput()

	fmt.Print("Enter amount (in smallest currency unit, e.g., cents): ")
	amountStr := c.readInput()
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		fmt.Printf("Error: Invalid amount: %v\n", err)
		return
	}

	fmt.Print("Enter currency code: ")
	currency := strings.ToUpper(c.readInput())

	fmt.Print("Enter holding currency (e.g., INR): ")
	holdingCurrency := strings.ToUpper(c.readInput())

	fmt.Print("Enter statement descriptor: ")
	statementDescriptor := c.readInput()

	payout := &tazapay.Payout{
		BeneficiaryDetails: tazapay.BeneficiaryDetails{
			DestinationDetails: struct {
				Type string              `json:"type"`
				Bank tazapay.BankDetails `json:"bank"`
			}{
				Type: "bank",
			},
		},
		Purpose:             "PYR002",
		Amount:              amount,
		Currency:            currency,
		Beneficiary:         beneficiaryID,
		HoldingCurrency:     holdingCurrency,
		Type:                "local",
		ChargeType:          "shared",
		StatementDescriptor: statementDescriptor,
	}

	response, err := c.client.CreatePayout(payout)
	if err != nil {
		fmt.Printf("Error creating payout: %v\n", err)
		return
	}

	fmt.Printf("\nPayout created successfully!\n")
	fmt.Printf("Status: %s\n", response.Status)
	if response.Message != "" {
		fmt.Printf("Message: %s\n", response.Message)
	}
	fmt.Printf("Payout Details:\n")
	fmt.Printf("- ID: %s\n", response.Data.ID)
	fmt.Printf("- Status: %s\n", response.Data.Status)
	fmt.Printf("- Amount: %d %s\n", response.Data.Amount, response.Data.Currency)
	fmt.Printf("- Created At: %s\n", response.Data.CreatedAt)
	fmt.Printf("- Updated At: %s\n", response.Data.UpdatedAt)
	fmt.Printf("- Purpose: %s\n", response.Data.Purpose)
	fmt.Printf("- Type: %s\n", response.Data.Type)
	fmt.Printf("- Charge Type: %s\n", response.Data.ChargeType)
	fmt.Printf("- Beneficiary: %s\n", response.Data.Beneficiary)
}

func (c *ChatInterface) handleGetFXRates() {
	fmt.Println("\nGet Exchange Rate")
	fmt.Print("Enter source currency (e.g., USD): ")
	fromCurrency := c.readInput()
	fmt.Print("Enter target currency (e.g., EUR): ")
	toCurrency := c.readInput()
	fmt.Print("Enter amount: ")
	amountStr := c.readInput()
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Printf("Error: Invalid amount: %v\n", err)
		return
	}

	rate, err := c.client.GetExchangeRate(fromCurrency, toCurrency, amount)
	if err != nil {
		fmt.Printf("Error getting exchange rate: %v\n", err)
		return
	}

	fmt.Printf("\nExchange Rate Details:\n")
	fmt.Printf("From Currency: %s\n", rate.Data.InitialCurrency)
	fmt.Printf("To Currency: %s\n", rate.Data.FinalCurrency)
	fmt.Printf("Amount: %.2f %s\n", float64(rate.Data.Amount)/100, rate.Data.InitialCurrency)
	fmt.Printf("Exchange Rate: %.6f\n", rate.Data.ExchangeRate)
	fmt.Printf("Converted Amount: %.2f %s\n", float64(rate.Data.ConvertedAmount)/100, rate.Data.FinalCurrency)
	fmt.Printf("Timestamp: %s\n", rate.Data.Timestamp)
}

func (c *ChatInterface) readInput() string {
	input, _ := c.reader.ReadString('\n')
	return strings.TrimSpace(input)
}
