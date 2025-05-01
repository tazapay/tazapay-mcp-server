package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"tazapay-mcp-server/internal/config"
	"tazapay-mcp-server/internal/tazapay"
)

func main() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	// Try to find the config file in the config directory first
	configPath := filepath.Join(wd, "config", ".tazapay-mcp-server.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try the parent directory's config folder
		configPath = filepath.Join(filepath.Dir(wd), "config", ".tazapay-mcp-server.yml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("Error: Config file not found in %s or config directory\n", wd)
			os.Exit(1)
		}
	}

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize client
	client := tazapay.NewClient(cfg.APIKey, cfg.APISecret)

	// Interactive CLI
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nTazapay CLI")
		fmt.Println("1. Check Balance")
		fmt.Println("2. Create Beneficiary")
		fmt.Println("3. Create Payout")
		fmt.Println("4. Get Exchange Rate")
		fmt.Println("5. Create Payment")
		fmt.Println("6. Exit")
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
			fmt.Println("Your balances:")
			for _, bal := range balance.Data.Available {
				fmt.Printf("- %s %s\n", bal.Amount, bal.Currency)
			}
			fmt.Printf("Last updated: %s\n", balance.Data.UpdatedAt)

		case "2":
			fmt.Print("Enter beneficiary name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Enter email: ")
			email, _ := reader.ReadString('\n')
			email = strings.TrimSpace(email)

			fmt.Print("Enter type (individual/business): ")
			beneficiaryType, _ := reader.ReadString('\n')
			beneficiaryType = strings.TrimSpace(beneficiaryType)

			fmt.Print("Enter bank country code (e.g., VN): ")
			bankCountry, _ := reader.ReadString('\n')
			bankCountry = strings.TrimSpace(bankCountry)

			fmt.Print("Enter bank currency (e.g., VND): ")
			bankCurrency, _ := reader.ReadString('\n')
			bankCurrency = strings.TrimSpace(bankCurrency)

			fmt.Print("Enter bank name: ")
			bankName, _ := reader.ReadString('\n')
			bankName = strings.TrimSpace(bankName)

			fmt.Print("Enter account number: ")
			accountNumber, _ := reader.ReadString('\n')
			accountNumber = strings.TrimSpace(accountNumber)

			fmt.Print("Enter SWIFT code: ")
			swiftCode, _ := reader.ReadString('\n')
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

			response, err := client.CreateBeneficiary(beneficiary)
			if err != nil {
				fmt.Printf("Error creating beneficiary: %v\n", err)
				continue
			}
			fmt.Printf("Beneficiary created successfully! ID: %s\n", response.Data.ID)

		case "3":
			fmt.Print("Enter beneficiary ID: ")
			beneficiaryID, _ := reader.ReadString('\n')
			beneficiaryID = strings.TrimSpace(beneficiaryID)

			fmt.Print("Enter amount (in smallest currency unit, e.g., cents): ")
			amountStr, _ := reader.ReadString('\n')
			amount, err := strconv.ParseInt(strings.TrimSpace(amountStr), 10, 64)
			if err != nil {
				fmt.Printf("Error: Invalid amount: %v\n", err)
				continue
			}

			fmt.Print("Enter currency code: ")
			currency, _ := reader.ReadString('\n')
			currency = strings.ToUpper(strings.TrimSpace(currency))

			fmt.Print("Enter holding currency (e.g., INR): ")
			holdingCurrency, _ := reader.ReadString('\n')
			holdingCurrency = strings.ToUpper(strings.TrimSpace(holdingCurrency))

			payout := &tazapay.Payout{
				Beneficiary:     beneficiaryID,
				Amount:          amount,
				Currency:        currency,
				HoldingCurrency: holdingCurrency,
				Type:            "local",
				ChargeType:      "shared",
				Purpose:         "PYR002",
			}

			response, err := client.CreatePayout(payout)
			if err != nil {
				fmt.Printf("Error creating payout: %v\n", err)
				continue
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

		case "4":
			fmt.Print("Enter from currency: ")
			fromCurrency, _ := reader.ReadString('\n')
			fromCurrency = strings.TrimSpace(fromCurrency)

			fmt.Print("Enter to currency: ")
			toCurrency, _ := reader.ReadString('\n')
			toCurrency = strings.TrimSpace(toCurrency)

			fmt.Print("Enter amount: ")
			amountStr, _ := reader.ReadString('\n')
			amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
			if err != nil {
				fmt.Printf("Error: Invalid amount: %v\n", err)
				continue
			}

			rate, err := client.GetExchangeRate(fromCurrency, toCurrency, amount)
			if err != nil {
				fmt.Printf("Error getting exchange rate: %v\n", err)
				continue
			}

			fmt.Printf("\nExchange Rate Details:\n")
			fmt.Printf("From: %s\n", rate.Data.InitialCurrency)
			fmt.Printf("To: %s\n", rate.Data.FinalCurrency)
			fmt.Printf("Rate: %.6f\n", rate.Data.ExchangeRate)
			fmt.Printf("Amount: %.2f %s\n", float64(rate.Data.Amount)/100, rate.Data.InitialCurrency)
			fmt.Printf("Converted: %.2f %s\n", float64(rate.Data.ConvertedAmount)/100, rate.Data.FinalCurrency)
			fmt.Printf("Timestamp: %s\n", rate.Data.Timestamp)

		case "5":
			fmt.Print("Enter amount: ")
			amountStr, _ := reader.ReadString('\n')
			amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
			if err != nil {
				fmt.Printf("Error: Invalid amount: %v\n", err)
				continue
			}

			fmt.Print("Enter currency (e.g., USD): ")
			currency, _ := reader.ReadString('\n')
			currency = strings.ToUpper(strings.TrimSpace(currency))

			fmt.Print("Enter description: ")
			description, _ := reader.ReadString('\n')
			description = strings.TrimSpace(description)

			fmt.Print("Enter success URL: ")
			successURL, _ := reader.ReadString('\n')
			successURL = strings.TrimSpace(successURL)

			fmt.Print("Enter cancel URL: ")
			cancelURL, _ := reader.ReadString('\n')
			cancelURL = strings.TrimSpace(cancelURL)

			fmt.Print("Enter customer email: ")
			customerEmail, _ := reader.ReadString('\n')
			customerEmail = strings.TrimSpace(customerEmail)

			fmt.Print("Enter customer name: ")
			customerName, _ := reader.ReadString('\n')
			customerName = strings.TrimSpace(customerName)

			payment := &tazapay.PaymentRequest{
				Amount:          amount,
				Currency:        currency,
				InvoiceCurrency: currency,
				Description:     description,
				TransactionDesc: description,
				SuccessURL:      successURL,
				CancelURL:       cancelURL,
				CustomerEmail:   customerEmail,
				CustomerName:    customerName,
			}

			// Set customer details
			payment.CustomerDetails.Email = customerEmail
			payment.CustomerDetails.Name = customerName
			payment.CustomerDetails.Country = "US"          // Default to US
			payment.CustomerDetails.Phone.CallingCode = "1" // Default to US

			response, err := client.CreatePayment(payment)
			if err != nil {
				fmt.Printf("Error creating payment: %v\n", err)
				continue
			}

			fmt.Printf("\nPayment created successfully!\n")
			fmt.Printf("Payment Details:\n")
			fmt.Printf("- ID: %s\n", response.Data.ID)
			fmt.Printf("- Amount: %.2f %s\n", response.Data.Amount, response.Data.Currency)
			fmt.Printf("- Status: %s\n", response.Data.Status)
			fmt.Printf("\nCheckout URL: %s\n", response.Data.URL)

		case "6":
			fmt.Println("Goodbye!")
			os.Exit(0)

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
