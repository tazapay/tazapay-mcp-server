# Tazapay MCP Server

This project provides a command-line interface and server implementation for interacting with Tazapay's API services.

## Project Structure

```
tazapay-mcp-server/
├── cmd/                    # Main applications
│   ├── cli/               # Command-line interface
│   │   └── main.go
│   └── server/            # MCP server
│       └── main.go
├── internal/              # Private application code
│   ├── config/           # Configuration handling
│   ├── tazapay/          # Tazapay API client
│   └── tools/            # MCP tools
├── pkg/                   # Public library code
│   ├── api/              # API definitions
│   └── mcp/              # MCP server implementation
├── test/                 # Test files
├── config/               # Configuration files
├── docs/                 # Documentation
└── scripts/             # Build and deployment scripts
```

## Features

- Balance checking
- Beneficiary management
- Payout creation
- Payment processing
- Exchange rate calculation
- Log analysis

## Configuration

Create a `.tazapay-mcp-server.yml` file in the config directory with your API credentials:

```yaml
api_key: your_api_key
api_secret: your_api_secret
```

## Usage

### CLI

```bash
go run cmd/cli/main.go
```

### Server

```bash
go run cmd/server/main.go
```

## Development

1. Clone the repository
2. Set up your configuration file
3. Run tests: `go test ./...`
4. Build: `go build ./...`

## License

Proprietary - All rights reserved 