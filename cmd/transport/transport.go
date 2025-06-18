package transport

import (
	"context"
	"log/slog"
	"net/http"

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

	// Create base mux for health check
	baseMux := http.NewServeMux()
	baseMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Create streamable HTTP server
	streamServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/stream"),
		server.WithHTTPContextFunc(utils.AuthHeaderHTTPContextFunc),
	)

	// Create a wrapper handler that handles both health check and stream
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			baseMux.ServeHTTP(w, r)
			return
		}
		streamServer.ServeHTTP(w, r)
	})

	// Create the main HTTP server
	httpServer := &http.Server{
		Handler: mainHandler,
	}

	defer func() {
		streamServer.Shutdown(context.Background())
		httpServer.Shutdown(context.Background())
	}()

	addr := viper.GetString("STREAM_SERVER_ADDR")
	if addr == "" {
		addr = ":8081"
	}
	httpServer.Addr = addr

	return httpServer.ListenAndServe()
}
