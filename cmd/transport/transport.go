package transport

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
	"github.com/tazapay/tazapay-mcp-server/pkg/utils"
)

// HandleStdioServer starts the MCP server using stdio transport.
// It logs the server start and delegates to server.ServeStdio.
func HandleStdioServer(s *server.MCPServer, logger *slog.Logger) error {
	// Only log on actual start
	logger.InfoContext(context.Background(), "Stdio server started")
	return server.ServeStdio(s)
}

// HandleStreamableHTTPServer starts the MCP server using a streamable HTTP transport.
// It sets up the HTTP server with endpoint path and authentication context, logs the start,
// and listens on the configured address (default :8081).
func HandleStreamableHTTPServer(s *server.MCPServer, logger *slog.Logger) error {
	// Only log on actual start
	logger.InfoContext(context.Background(), "Streamable HTTP server started")
	streamServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/stream"),
		server.WithHTTPContextFunc(utils.AuthHeaderHTTPContextFunc),
	)
	defer streamServer.Shutdown(context.Background())
	addr := viper.GetString("STREAM_SERVER_ADDR")
	if addr == "" {
		addr = ":8081"
	}
	return streamServer.Start(addr)
}
