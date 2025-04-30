package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/tazapay/tazapay-mcp-server/pkg/tazapay"
)

func initConfig() (string, string, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("error getting working directory: %w", err)
	}

	// Try both .yml and .yaml extensions
	configPaths := []string{
		filepath.Join(wd, ".tazapay-mcp-server.yml"),
		filepath.Join(wd, ".tazapay-mcp-server.yaml"),
		filepath.Join(os.Getenv("HOME"), ".tazapay-mcp-server.yml"),
		filepath.Join(os.Getenv("HOME"), ".tazapay-mcp-server.yaml"),
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile == "" {
		return "", "", fmt.Errorf("config file not found. Tried: %v", configPaths)
	}

	fmt.Printf("Using config file: %s\n", configFile)

	// Read the config file directly
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return "", "", fmt.Errorf("error reading config file: %w", err)
	}

	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	if apiKey == "" || apiSecret == "" {
		return "", "", fmt.Errorf("api_key and api_secret must be set in config file")
	}

	fmt.Println("Successfully loaded configuration")
	return apiKey, apiSecret, nil
}

func main() {
	// Initialize configuration
	apiKey, apiSecret, err := initConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize client
	client := tazapay.NewClient(apiKey, apiSecret)

	// Interactive CLI
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nTazapay CLI")
		fmt.Println("1. Check Balance")
		fmt.Println("2. Create Beneficiary")
		fmt.Println("3. Create Payout")
		fmt.Println("4. Get Exchange Rate")
		fmt.Println("5. Exit")
		fmt.Print("Select an option: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			balance, err := client.GetBalance()
			if err != nil {
				fmt.Printf("Error checking balance: %v\n", err)
				continue
			}
			fmt.Printf("Balance: %.2f %s\n", balance.Data.Balance, balance.Data.Currency)

		case "2":
			fmt.Print("Enter beneficiary name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Enter beneficiary email: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)

			fmt.Print("Enter beneficiary country: ")
			country, _ := reader.ReadString('\n')
			country = strings.TrimSpace(country)

			beneficiary := &tazapay.Beneficiary{
				Name:    name,
				Email:   email,
				Country: country,
			}

			response, err := client.CreateBeneficiary(beneficiary)
			if err != nil {
				fmt.Printf("Error creating beneficiary: %v\n", err)
				continue
			}
			fmt.Printf("Beneficiary created with ID: %s\n", response.Data.ID)

		case "3":
			fmt.Print("Enter beneficiary ID: ")
			beneficiaryID, _ := reader.ReadString('\n')
			beneficiaryID = strings.TrimSpace(beneficiaryID)

			fmt.Print("Enter amount: ")
			var amount float64
			fmt.Scanf("%f", &amount)

			fmt.Print("Enter currency: ")
			currency, _ := reader.ReadString('\n')
			currency = strings.TrimSpace(currency)

			payout := &tazapay.Payout{
				BeneficiaryID: beneficiaryID,
				Amount:        amount,
				Currency:      currency,
			}

			response, err := client.CreatePayout(payout)
			if err != nil {
				fmt.Printf("Error creating payout: %v\n", err)
				continue
			}
			fmt.Printf("Payout created with ID: %s\n", response.Data.ID)

		case "4":
			fmt.Print("Enter from currency: ")
			fromCurrency, _ := reader.ReadString('\n')
			fromCurrency = strings.TrimSpace(fromCurrency)

			fmt.Print("Enter to currency: ")
			toCurrency, _ := reader.ReadString('\n')
			toCurrency = strings.TrimSpace(toCurrency)

			rate, err := client.GetExchangeRate(fromCurrency, toCurrency)
			if err != nil {
				fmt.Printf("Error getting exchange rate: %v\n", err)
				continue
			}
			fmt.Printf("Exchange rate: %.6f\n", rate.Data.Rate)

		case "5":
			fmt.Println("Exiting...")
			os.Exit(0)

		default:
			fmt.Println("Invalid option")
		}
	}
}
