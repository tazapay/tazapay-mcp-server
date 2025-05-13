// types/tool.go
package types

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// Tool defines an interface that all tools must implement
type Tool interface {
	// Definition returns the tool definition
	Definition() mcp.Tool

	// Handle processes the tool call
	Handle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
