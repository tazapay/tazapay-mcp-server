# Tazapay MCP Server

The Tazapay MCP Server is a Model Context Protocol (MCP) server designed to bridge LLM-based tools with Tazapay's Payments API. It supports structured tool invocation for seamless payment automation and foreign exchange capabilities.

## Features

* ✅ Payment Link Generation via Tazapay API
* 🌍 Real-Time FX Rate Conversion
* 🧩 Modular Tool Architecture
* 🔗 Fully Compatible with Anthropic Claude, GitHub Copilot, Cursor IDE
* 📝 Roadmap: Global Payout Tools, Refund Tools.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go |
| Framework | Anthropic MCP Plugin (mark3labs) |
| Containerization | Docker |

## Tools Overview

#### 1. `generate_payment_link_tool`
* **Input:**
   * `invoice_currency` (string)
   * `payment_amount` (number)
   * `customer_name` (string)
   * `customer_email` (string)
   * `customer_country` (string)
   * `transaction_description` (string)
* **Output:** Shareable Tazapay payment link

#### 2. `tazapay_fetch_fx_tool`
* **Input:**
   * `from_currency` (string)
   * `to_currency` (string)
   * `amount` (number)
* **Output:** FX rate and converted amount

#### 3. `tazapay_fetch_balance_tool`
* **Input:**
  * `currency`(optional string) – If specified, returns the balance in the given currency.
* **Output:** Returns the current available balance in the merchant’s account.

## Prerequisites

Ensure the following tools are installed before setup:

* Go 1.24+
* An IDE or tool that supports MCP (e.g., Claude Desktop, GitHub Copilot, Cursor IDE)

## Get started with Docker and Claude integration

- Ensure Docker is installed on your desktop
  
- Pull the Following Docker Image from dockerhub into your local machine
  ```bash
  docker pull tazapay/tazapay-mcp-server:latest
  ```

- Add the following to your `claude_desktop_config.json`:
   ```json
   {
     "mcpServers": {
       "Tazapay-Docker-Server": {
         "command": "docker",
         "description": "MCP server to integrate Tazapay API's and payments solutions.",
         "args": [
           "run","--rm","-i",
           "-e","TAZAPAY_API_KEY",
           "-e","TAZAPAY_API_SECRET",
           "tazapay/tazapay-mcp-server:latest"
         ],
         "env": {
           "TAZAPAY_API_KEY": "your_tazapay_api_key",
           "TAZAPAY_API_SECRET": "your_tazapay_api_secret"
         }
       }
     }
   }
   ```
- Now you are ready to interact with LLM to take care of operations with your Tazapay account.

## Setting Up Locally 

* Add a `.tazapay-mcp-server.yaml` config file in your home directory with the following content:

   ```yaml
   TAZAPAY_API_KEY: "your_key"
   TAZAPAY_API_SECRET: "your_secret"
   ```
   
- Verify that the file '.tazapay-mcp-server.yaml' is added to your home directory. If not add the file there.
  ```bash
  [ -f "$HOME/.tazapay-mcp-server.yaml" ] && echo "Config file found." || echo "Config file missing at $HOME/.tazapay-mcp-server.yaml"
  ```
- Clone Repo and Build Locally

   ```bash
   git clone https://github.com/tazapay/tazapay-mcp-server.git
   cd tazapay-mcp-server
   go build -o tazapay-mcp-server ./cmd/server
   ```
   The binary `tazapay-mcp-server` will be available post build.

- In Claude Desktop, Add the following to your `claude_desktop_config.json`:
   
   ```json
   {
     "mcpServers": {
       "Tazapay-mcp-server": {
         "command": "/absolute/path/to/repo/tazapay-mcp-server",
         "description": "MCP server to integrate Tazapay API's and payments solutions.",
       }
     }
   }
   ```
- Now you are ready to interact with LLM to take care of operations with your Tazapay account.

## Integration With other popular IDE 

### GitHub Copilot Chat in VS code
* Navigate to: **Settings > Features > Chat > Enable MCP**
* Now add the given json file in **settings.json**
  ```json
  {
    "mcp": {
      "inputs": [],
      "servers": {
        "Tazapay-mcp-server": {
          "command": "/absolute/path/to/repo/tazapay-mcp-server",
          "description": "MCP server to integrate Tazapay APIs and payment solutions."
        }
      }
    }
  }
  ```
* Configure MCP tools via the **gear icon** in the Chat tab.

### Cursor IDE
* Go to: **Settings > MCP > Add New Global MCP Server**
* Paste the JSON configuration from above and tools are ready to use within the chat. 

## Server Modes and Endpoints

The server can run in two modes:

### 1. Standard I/O Mode (Default)
Used for direct integration with tools like Claude Desktop and VS Code.

### 2. Streamable HTTP Mode
When running in HTTP mode (set via `TRANSPORT_TYPE=streamablehttp`), the server exposes these endpoints:

#### Available Endpoints

All endpoints are available on port 8081:

* **MCP Stream**: `http://localhost:8081/stream` - Main MCP protocol endpoint
* **Health Check**: `http://localhost:8081/healthz` - Server health status endpoint

#### Running in HTTP Mode

```bash
# Run with streamable HTTP server mode
docker run --rm -p 8081:8081 \
  -e TAZAPAY_API_KEY=your_api_key \
  -e TAZAPAY_API_SECRET=your_api_secret \
  -e SERVER_TYPE=streamablehttp \
  tazapay/tazapay-mcp-server:latest
```

#### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_TYPE` | Server mode (`stdio` or `streamablehttp`) | `stdio` | No |
| `STREAM_SERVER_ADDR` | HTTP server address | `:8081` | No |
| `TAZAPAY_API_KEY` | Your Tazapay API key | - | Yes |
| `TAZAPAY_API_SECRET` | Your Tazapay API secret | - | Yes |

## License

This project is licensed under the MIT license. Refer to LICENSE for details.

[tazapay-License](https://github.com/tazapay/tazapay-mcp-server/blob/feat-balance-tool-addition/LICENSE)
