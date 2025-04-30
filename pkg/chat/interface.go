package chat

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tazapay/tazapay-mcp-server/pkg/tazapay"
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

	fmt.Printf("Your balance is %.2f %s\n", balance.Data.Balance, balance.Data.Currency)
}

func (c *ChatInterface) handleCreateBeneficiary() {
	fmt.Print("Enter beneficiary name: ")
	name, _ := c.reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter email: ")
	email, _ := c.reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Enter country: ")
	country, _ := c.reader.ReadString('\n')
	country = strings.TrimSpace(country)

	beneficiary := &tazapay.Beneficiary{
		Name:    name,
		Email:   email,
		Country: country,
	}
	response, err := c.client.CreateBeneficiary(beneficiary)
	if err != nil {
		fmt.Printf("Error creating beneficiary: %v\n", err)
		return
	}

	fmt.Printf("Beneficiary created successfully! ID: %s\n", response.Data.ID)
}

func (c *ChatInterface) handleCreatePayout() {
	fmt.Print("Enter beneficiary ID: ")
	beneficiaryID, _ := c.reader.ReadString('\n')
	beneficiaryID = strings.TrimSpace(beneficiaryID)

	fmt.Print("Enter amount: ")
	amountStr, _ := c.reader.ReadString('\n')
	amount, _ := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)

	fmt.Print("Enter currency: ")
	currency, _ := c.reader.ReadString('\n')
	currency = strings.TrimSpace(currency)

	payout := &tazapay.Payout{
		BeneficiaryID: beneficiaryID,
		Amount:        amount,
		Currency:      currency,
	}
	response, err := c.client.CreatePayout(payout)
	if err != nil {
		fmt.Printf("Error creating payout: %v\n", err)
		return
	}

	fmt.Printf("Payout created successfully! ID: %s\n", response.Data.ID)
}

func (c *ChatInterface) handleGetFXRates() {
	fmt.Print("Enter source currency (e.g., USD): ")
	from, _ := c.reader.ReadString('\n')
	from = strings.TrimSpace(from)

	fmt.Print("Enter target currency (e.g., EUR): ")
	to, _ := c.reader.ReadString('\n')
	to = strings.TrimSpace(to)

	rate, err := c.client.GetExchangeRate(from, to)
	if err != nil {
		fmt.Printf("Error getting exchange rate: %v\n", err)
		return
	}

	fmt.Printf("Exchange rate from %s to %s: %.4f\n", from, to, rate.Data.Rate)
}
