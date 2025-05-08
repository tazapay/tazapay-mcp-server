package registertool

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tazapay/tazapay-mcp-server/tools/tazapay"
	"github.com/tazapay/tazapay-mcp-server/types"
)

// RegisterTools registers all tools with the server
func RegisterTools(s *server.MCPServer) {
	// Register all tools
	registerTool(s, tazapay.NewFXTool())
	registerTool(s, tazapay.NewPaymentLinkTool())
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
