package main

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"

	"github.com/tazapay/tazapay-mcp-server/cmd/transport"
	"github.com/tazapay/tazapay-mcp-server/constants"
	"github.com/tazapay/tazapay-mcp-server/pkg/log"
	tools "github.com/tazapay/tazapay-mcp-server/tools/register"
)

func main() {
	transportType := os.Getenv("TRANSPORT_TYPE")
	if transportType == "" {
		transportType = constants.TransportTypeStreamableHTTP
	}

	// set log configs
	var logConfig = log.Config{
		Format:   "json",
		FilePath: viper.GetString("LOG_FILE_PATH"),
	}

	// create logger
	logger, _, logErr := log.New(logConfig)
	if logErr != nil {
		os.Exit(1)
	}

	//create server and register tools
	s := server.NewMCPServer("tazapay", "0.1.2")
	tools.RegisterTools(s, logger)

	// Only keep this high-level log
	logger.InfoContext(context.Background(), "Tazapay MCP Server started", "Transport type", transportType)

	// based on server type start server and handle accordingly
	switch transportType {
	case constants.TransportTypeStdio:
		if err := transport.HandleStdioServer(s, logger); err != nil {
			logger.ErrorContext(context.Background(), "server exited with error", "error", err)
			os.Exit(1)
		}
	default:
		if err := transport.HandleStreamableHTTPServer(s, logger); err != nil {
			logger.ErrorContext(context.Background(), "server exited with error", "error", err)
			os.Exit(1)
		}
	}
}
