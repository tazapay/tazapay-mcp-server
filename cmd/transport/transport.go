package transport

import (
	"context"
	"log/slog"
	"net/http"
	"os"

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
	viper.AutomaticEnv()

	// base mux for health check - handle root and any non-stream paths
	baseMux := http.NewServeMux()
	baseMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Create streamable HTTP server
	streamServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/stream"),
		server.WithHTTPContextFunc(utils.AuthHeaderHTTPContextFunc),
	)

	// Create a wrapper handler that handles health check for all paths except /stream
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/stream" {
			streamServer.ServeHTTP(w, r)
			return
		}
		baseMux.ServeHTTP(w, r)
	})

	// Create the main HTTP server
	httpServer := &http.Server{
		Handler: mainHandler,
	}

	defer func() {
		streamServer.Shutdown(context.Background())
		httpServer.Shutdown(context.Background())
	}()

	// Viper env then fallback to os environment variable, then default
	addr := viper.GetString("STREAM_SERVER_ADDR")
	if addr == "" {
		addr = os.Getenv("STREAM_SERVER_ADDR")
	}
	if addr == "" {
		addr = ":8081"
	}
	httpServer.Addr = addr

	// Log the server startup with correct spelling and actual address
	logger.InfoContext(context.Background(), "HTTP Server started",
		"running_port_addr", addr,
		"env_var_input", os.Getenv("STREAM_SERVER_ADDR"))

	return httpServer.ListenAndServe()
}

// HandleSseServer starts a custom SSE server with custom session management.
// It does not use server.NewSSEServer from mark3build.

func HandleSseServer(logger *slog.Logger) error {
	viper.AutomaticEnv()

	// Use own session-based SSE server
	sseServer := NewMySSEServer()

	// Base mux for health check
	baseMux := http.NewServeMux()
	baseMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/sse":
			ctx := utils.AuthHeaderHTTPContextFunc(r.Context(), r)
			r = r.WithContext(ctx)
			sseServer.HandleSSE(w, r)
		case "/message":
			w.Header().Set("Access-Control-Allow-Origin", "*")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			ctx := utils.AuthHeaderHTTPContextFunc(r.Context(), r)
			r = r.WithContext(ctx)
			sseServer.HandleMessage(w, r)
		case "/execute":
			ctx := utils.AuthHeaderHTTPContextFunc(r.Context(), r)
			r = r.WithContext(ctx)
			sseServer.HandleDirectExecution(w, r)
		default:
			baseMux.ServeHTTP(w, r)
		}
	})

	httpServer := &http.Server{
		Handler: mainHandler,
	}

	defer func() {
		httpServer.Shutdown(context.Background())
	}()

	addr := viper.GetString("STREAM_SERVER_ADDR")
	if addr == "" {
		addr = os.Getenv("STREAM_SERVER_ADDR")
	}
	if addr == "" {
		addr = ":8081"
	}
	httpServer.Addr = addr

	logger.InfoContext(context.Background(), "Custom SSE Server started",
		"running_port_addr", addr,
		"env_var_input", os.Getenv("STREAM_SERVER_ADDR"))

	return httpServer.ListenAndServe()
}
