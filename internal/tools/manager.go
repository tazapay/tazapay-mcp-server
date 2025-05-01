package tools

import (
	"tazapay-mcp-server/internal/logs"
	"tazapay-mcp-server/internal/tazapay"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolManager handles registration and management of all tools
type ToolManager struct {
	client *tazapay.Client
	server *server.MCPServer
}

// NewToolManager creates a new tool manager
func NewToolManager(client *tazapay.Client, server *server.MCPServer) *ToolManager {
	return &ToolManager{
		client: client,
		server: server,
	}
}

// RegisterAllTools registers all available tools with the server
func (tm *ToolManager) RegisterAllTools() {
	// Register balance tool
	tm.server.AddTool(
		mcp.NewTool("check_balance",
			mcp.WithDescription("Check account balance"),
		),
		server.ToolHandlerFunc(tm.handleBalanceTool),
	)

	// Register FX tool
	tm.server.AddTool(
		mcp.NewTool("get_fx_rate",
			mcp.WithDescription("Get exchange rate between currencies"),
			mcp.WithString("from_currency",
				mcp.Required(),
				mcp.Description("Source currency code (e.g., USD)"),
			),
			mcp.WithString("to_currency",
				mcp.Required(),
				mcp.Description("Target currency code (e.g., EUR)"),
			),
			mcp.WithNumber("amount",
				mcp.Required(),
				mcp.Description("Amount to convert"),
			),
		),
		server.ToolHandlerFunc(tm.handleFXTool),
	)

	// Register payment tool
	tm.server.AddTool(
		mcp.NewTool("create_payment",
			mcp.WithDescription("Create a new payment"),
			mcp.WithNumber("amount",
				mcp.Required(),
				mcp.Description("Payment amount"),
			),
			mcp.WithString("currency",
				mcp.Required(),
				mcp.Description("Payment currency (e.g., USD)"),
			),
			mcp.WithString("description",
				mcp.Required(),
				mcp.Description("Payment description"),
			),
			mcp.WithString("success_url",
				mcp.Required(),
				mcp.Description("URL to redirect on success"),
			),
			mcp.WithString("cancel_url",
				mcp.Required(),
				mcp.Description("URL to redirect on cancel"),
			),
			mcp.WithString("customer_email",
				mcp.Required(),
				mcp.Description("Customer email"),
			),
			mcp.WithString("customer_name",
				mcp.Required(),
				mcp.Description("Customer name"),
			),
		),
		server.ToolHandlerFunc(tm.handlePaymentTool),
	)

	// Register payout tool
	tm.server.AddTool(
		mcp.NewTool("create_payout",
			mcp.WithDescription("Create a new payout"),
			mcp.WithString("beneficiary_id",
				mcp.Required(),
				mcp.Description("Beneficiary ID"),
			),
			mcp.WithNumber("amount",
				mcp.Required(),
				mcp.Description("Payout amount"),
			),
			mcp.WithString("currency",
				mcp.Required(),
				mcp.Description("Payout currency (e.g., USD)"),
			),
			mcp.WithString("holding_currency",
				mcp.Required(),
				mcp.Description("Holding currency (e.g., INR)"),
			),
		),
		server.ToolHandlerFunc(tm.handlePayoutTool),
	)

	// Register beneficiary tool
	tm.server.AddTool(
		mcp.NewTool("create_beneficiary",
			mcp.WithDescription("Create a new beneficiary"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Beneficiary name"),
			),
			mcp.WithString("type",
				mcp.Required(),
				mcp.Description("Beneficiary type (individual/business)"),
			),
			mcp.WithString("email",
				mcp.Required(),
				mcp.Description("Beneficiary email"),
			),
			mcp.WithString("bank_country",
				mcp.Required(),
				mcp.Description("Bank country code (e.g., VN)"),
			),
			mcp.WithString("bank_currency",
				mcp.Required(),
				mcp.Description("Bank currency code (e.g., VND)"),
			),
			mcp.WithString("bank_name",
				mcp.Required(),
				mcp.Description("Bank name"),
			),
			mcp.WithString("account_number",
				mcp.Required(),
				mcp.Description("Bank account number"),
			),
			mcp.WithString("swift_code",
				mcp.Required(),
				mcp.Description("SWIFT code"),
			),
		),
		server.ToolHandlerFunc(tm.handleBeneficiaryTool),
	)

	// Register log analysis tool
	tm.server.AddTool(
		mcp.NewTool("analyze_logs",
			mcp.WithDescription("Fetch and analyze recent logs"),
			mcp.WithNumber("count",
				mcp.Description("Number of logs to fetch (default: 9)"),
			),
		),
		server.ToolHandlerFunc(logs.HandleLogAnalysisTool(tm.client)),
	)
}
