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

## Prerequisites

- Go 1.16 or higher
- Tazapay API credentials

## Configuration

1. Create a `.tazapay-mcp-server.yml` file in the `config` directory with your API credentials:

```yaml
api_key: your_api_key
api_secret: your_api_secret
```

2. The configuration file can be placed in either:
   - The `config` directory in your project root
   - The parent directory's `config` folder
   - The current working directory

## Installation

1. Clone the repository:
```bash
git clone https://github.com/your-org/tazapay-mcp-server.git
cd tazapay-mcp-server
```

2. Install dependencies:
```bash
go mod download
```

## Usage

### CLI

The CLI provides an interactive interface for Tazapay API operations:

```bash
go run cmd/cli/main.go
```

Available options:
1. Check Balance - View your account balances
2. Create Beneficiary - Add a new beneficiary
3. Create Payout - Initiate a payout to a beneficiary
4. Get Exchange Rate - Calculate exchange rates between currencies
5. Create Payment - Create a new payment request
6. Exit - Quit the CLI

### Server

To run the MCP server:

```bash
go run cmd/server/main.go
```

## Development

1. Run tests:
```bash
go test ./...
```

2. Build the project:
```bash
go build ./...
```

3. Run specific tests:
```bash
go test -v ./test/payment_test.go
```

## Troubleshooting

If you encounter the error "Config File not found":
1. Ensure the config file exists in one of the supported locations
2. Verify the file name is exactly `.tazapay-mcp-server.yml`
3. Check file permissions

## License

Proprietary - All rights reserved 