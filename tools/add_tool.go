package tools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AddAddTool(s *server.MCPServer) {
	tool := mcp.NewTool("add",
		mcp.WithDescription("Add two numbers"),
		mcp.WithNumber("a",
			mcp.Required(),
			mcp.Description("First number to add"),
		),
		mcp.WithNumber("b",
			mcp.Required(),
			mcp.Description("Second number to add"),
		),
	)
	s.AddTool(tool, handleAddTool)
}

func handleAddTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	a, ok1 := arguments["a"].(float64)
	b, ok2 := arguments["b"].(float64)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("invalid number arguments")
	}
	sum := a + b
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("The sum of %f and %f is %f.", a, b, sum),
			},
		},
	}, nil
}
