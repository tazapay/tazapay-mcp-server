package registertool

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tazapay/tazapay-mcp-server/tools/tazapay/balance"
	"github.com/tazapay/tazapay-mcp-server/tools/tazapay/beneficiary"
	"github.com/tazapay/tazapay-mcp-server/tools/tazapay/checkout"
	"github.com/tazapay/tazapay-mcp-server/tools/tazapay/customer"
	"github.com/tazapay/tazapay-mcp-server/tools/tazapay/payin"
	"github.com/tazapay/tazapay-mcp-server/tools/tazapay/paymentattempt"
	"github.com/tazapay/tazapay-mcp-server/tools/tazapay/payout"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// RegisterTools registers all tools with the server

// NOTE: All tool constructors (e.g., NewFXTool, NewCreatePayinTool, etc.) must be lightweight.
// They should NOT perform any blocking or heavy operations (network calls, file I/O, etc.).
// Only assign struct fields and log. Any heavy setup should be deferred to the handler or background goroutines.
func RegisterTools(s *server.MCPServer, logger *slog.Logger) {
	tools := []types.Tool{
		balance.NewFXTool(logger),
		balance.NewBalanceTool(logger),
		payout.NewGetPayoutTool(logger),
		payout.NewFundPayoutTool(logger),
		payout.NewCreatePayoutTool(logger),
		payin.NewGetPayinTool(logger),
		payin.NewCreatePayinTool(logger),
		payin.NewUpdatePayinTool(logger),
		payin.NewCancelPayinTool(logger),
		//payin.NewConfirmPayinTool(logger),
		checkout.NewPaymentLinkTool(logger),
		checkout.NewFetchCheckoutTool(logger),
		checkout.NewExpireCheckoutTool(logger),
		beneficiary.NewGetBeneficiaryTool(logger),
		beneficiary.NewCreateBeneficiaryTool(logger),
		beneficiary.NewUpdateBeneficiaryTool(logger),
		paymentattempt.NewGetPaymentAttemptTool(logger),
		customer.NewCreateCustomerTool(logger),
		customer.NewFetchCustomerTool(logger),
	}

	for _, tool := range tools {
		registerTool(s, tool)
	}
}

// registerTool registers a single tool with the server
func registerTool(s *server.MCPServer, tool types.Tool) {
	s.AddTool(tool.Definition(), createHandler(tool))
}

// createHandler creates a handler function for a tool
func createHandler(tool types.Tool) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return tool.Handle(ctx, req)
	}
}
