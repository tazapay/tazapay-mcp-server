package tools

import (
	"context"
	"fmt"

	"tazapay-mcp-server/internal/tazapay"

	"github.com/mark3labs/mcp-go/mcp"
)

// Handle balance check
func (tm *ToolManager) handleBalanceTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	balance, err := tm.client.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("error checking balance: %w", err)
	}

	// Format balance information
	var balanceText string
	for _, bal := range balance.Data.Available {
		balanceText += fmt.Sprintf("%s %s\n", bal.Amount, bal.Currency)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Available Balance:\n%s", balanceText),
			},
		},
	}, nil
}

// Handle FX rate check
func (tm *ToolManager) handleFXTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	fromCurrency := args["from_currency"].(string)
	toCurrency := args["to_currency"].(string)
	amount := args["amount"].(float64)

	rate, err := tm.client.GetExchangeRate(fromCurrency, toCurrency, amount)
	if err != nil {
		return nil, fmt.Errorf("error getting exchange rate: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Exchange Rate Details:\nFrom: %s\nTo: %s\nRate: %.6f\nAmount: %.2f %s\nConverted: %.2f %s",
					rate.Data.InitialCurrency,
					rate.Data.FinalCurrency,
					rate.Data.ExchangeRate,
					float64(rate.Data.Amount)/100,
					rate.Data.InitialCurrency,
					float64(rate.Data.ConvertedAmount)/100,
					rate.Data.FinalCurrency,
				),
			},
		},
	}, nil
}

// Handle payment creation
func (tm *ToolManager) handlePaymentTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments

	payment := &tazapay.PaymentRequest{
		Amount:          args["amount"].(float64),
		Currency:        args["currency"].(string),
		InvoiceCurrency: args["currency"].(string),
		Description:     args["description"].(string),
		TransactionDesc: args["description"].(string),
		SuccessURL:      args["success_url"].(string),
		CancelURL:       args["cancel_url"].(string),
		CustomerEmail:   args["customer_email"].(string),
		CustomerName:    args["customer_name"].(string),
		CustomerDetails: struct {
			Email string `json:"email"`
			Name  string `json:"name"`
			Phone struct {
				Number      string `json:"number"`
				CallingCode string `json:"calling_code"`
			} `json:"phone"`
			Address string `json:"address"`
			Country string `json:"country"`
		}{
			Email:   args["customer_email"].(string),
			Name:    args["customer_name"].(string),
			Country: "BR", // Default to Brazil
			Phone: struct {
				Number      string `json:"number"`
				CallingCode string `json:"calling_code"`
			}{
				CallingCode: "1",
				Number:      args["customer_phone"].(string),
			},
			Address: args["customer_address"].(string),
		},
		PaymentMethods: []string{"card", "bank_transfer"},
	}

	response, err := tm.client.CreatePayment(payment)
	if err != nil {
		return nil, fmt.Errorf("error creating payment: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Payment created successfully!\n\nPayment Details:\nID: %s\nAmount: %.2f %s\nStatus: %s\n\nCheckout URL: %s",
					response.Data.ID,
					response.Data.Amount,
					response.Data.Currency,
					response.Data.Status,
					response.Data.URL),
			},
		},
	}, nil
}

// Handle payout creation
func (tm *ToolManager) handlePayoutTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	amount := int64(args["amount"].(float64))

	payout := &tazapay.Payout{
		Beneficiary:     args["beneficiary_id"].(string),
		Amount:          amount,
		Currency:        args["currency"].(string),
		HoldingCurrency: args["holding_currency"].(string),
		Type:            "local",
		ChargeType:      "shared",
		Purpose:         "PYR002",
	}

	response, err := tm.client.CreatePayout(payout)
	if err != nil {
		return nil, fmt.Errorf("error creating payout: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Payout created successfully!\n\nPayout Details:\nID: %s\nAmount: %d %s\nStatus: %s\nCreated: %s",
					response.Data.ID,
					response.Data.Amount,
					response.Data.Currency,
					response.Data.Status,
					response.Data.CreatedAt),
			},
		},
	}, nil
}

// Handle beneficiary creation
func (tm *ToolManager) handleBeneficiaryTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	beneficiary := &tazapay.Beneficiary{
		AccountID: "acc_d00eij4qqkhc9e5ats4g",
		Name:      args["name"].(string),
		Type:      args["type"].(string),
		Email:     args["email"].(string),
		DestinationDetails: tazapay.DestinationDetails{
			Type: "bank",
			Bank: tazapay.BankDetails{
				Country:       args["bank_country"].(string),
				Currency:      args["bank_currency"].(string),
				BankName:      args["bank_name"].(string),
				AccountNumber: args["account_number"].(string),
				BankCodes: struct {
					SwiftCode string `json:"swift_code"`
				}{
					SwiftCode: args["swift_code"].(string),
				},
			},
		},
		Phone: tazapay.Phone{
			CallingCode: "84", // Default to Vietnam
		},
	}

	response, err := tm.client.CreateBeneficiary(beneficiary)
	if err != nil {
		return nil, fmt.Errorf("error creating beneficiary: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Beneficiary created successfully!\n\nBeneficiary ID: %s", response.Data.ID),
			},
		},
	}, nil
}
