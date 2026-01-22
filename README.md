# Tazapay MCP Server

The Tazapay MCP Server is a Model Context Protocol (MCP) server designed to bridge LLM-based tools with Tazapay's Payments API. It supports structured tool invocation for seamless payment automation and foreign exchange capabilities.

## Features

* Payment Link Generation via Tazapay API
* Real-Time FX Rate Conversion
* Account Balance Retrieval
* Complete Payin/Payout Management
* Beneficiary Management
* Customer Management
* Checkout Session Management
* Modular Tool Architecture
* Fully Compatible with Anthropic Claude, GitHub Copilot, Cursor IDE
* Supports both Standard I/O and Streamable HTTP modes

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go |
| Framework | Anthropic MCP Plugin (mark3labs/mcp-go) |
| Containerization | Docker |

---

## Quick Start - Remote Server (Recommended)

The easiest way to use Tazapay MCP Server is via our hosted remote server.

**Remote Server URL:** `https://mcp.tazapay.com/stream`

### Claude Desktop Configuration (Remote)

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "Tazapay-Remote-Server": {
      "type": "streamable-http",
      "url": "https://mcp.tazapay.com/stream",
      "headers": {
        "Authorization": "Basic <base64_encoded_api_key:api_secret>"
      }
    }
  }
}
```

> **Note:** Replace `<base64_encoded_api_key:api_secret>` with your Base64-encoded `TAZAPAY_API_KEY:TAZAPAY_API_SECRET` string.

### VS Code / GitHub Copilot Configuration (Remote)

Add to your `settings.json`:

```json
{
  "mcp": {
    "servers": {
      "Tazapay-Remote-Server": {
        "type": "streamable-http",
        "url": "https://mcp.tazapay.com/stream",
        "headers": {
          "Authorization": "Basic <base64_encoded_api_key:api_secret>"
        }
      }
    }
  }
}
```

### Cursor IDE Configuration (Remote)

Go to **Settings > MCP > Add New Global MCP Server** and use:

```json
{
  "Tazapay-Remote-Server": {
    "type": "streamable-http",
    "url": "https://mcp.tazapay.com/stream",
    "headers": {
      "Authorization": "Basic <base64_encoded_api_key:api_secret>"
    }
  }
}
```

---

## Docker Setup

### Pull and Run with Docker

```bash
# Pull the Docker image
docker pull tazapay/tazapay-mcp-server:latest

# Run in Standard I/O mode (for Claude Desktop)
docker run --rm -i \
  -e TAZAPAY_API_KEY=your_api_key \
  -e TAZAPAY_API_SECRET=your_api_secret \
  tazapay/tazapay-mcp-server:latest

# Run in Streamable HTTP mode
docker run --rm -p 8081:8081 \
  -e TAZAPAY_API_KEY=your_api_key \
  -e TAZAPAY_API_SECRET=your_api_secret \
  -e SERVER_TYPE=streamablehttp \
  tazapay/tazapay-mcp-server:latest
```

### Claude Desktop Configuration (Docker)

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "Tazapay-Docker-Server": {
      "command": "docker",
      "description": "MCP server to integrate Tazapay API's and payments solutions.",
      "args": [
        "run", "--rm", "-i",
        "-e", "TAZAPAY_API_KEY",
        "-e", "TAZAPAY_API_SECRET",
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

---

## Local Development Setup

### Prerequisites

* Go 1.24+
* An IDE or tool that supports MCP (e.g., Claude Desktop, GitHub Copilot, Cursor IDE)

### Step 1: Create Configuration File

Add a `.tazapay-mcp-server.yaml` config file in your home directory:

```yaml
TAZAPAY_API_KEY: "your_key"
TAZAPAY_API_SECRET: "your_secret"
```

Verify the configuration file exists:

```bash
[ -f "$HOME/.tazapay-mcp-server.yaml" ] && echo "Config file found." || echo "Config file missing at $HOME/.tazapay-mcp-server.yaml"
```

### Step 2: Clone and Build

```bash
git clone https://github.com/tazapay/tazapay-mcp-server.git
cd tazapay-mcp-server
go build -o tazapay-mcp-server ./cmd/server
```

### Step 3: Configure Claude Desktop (Local Binary)

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "Tazapay-mcp-server": {
      "command": "/absolute/path/to/repo/tazapay-mcp-server",
      "description": "MCP server to integrate Tazapay API's and payments solutions."
    }
  }
}
```

---

## Integration with IDEs

### GitHub Copilot Chat in VS Code

1. Navigate to: **Settings > Features > Chat > Enable MCP**
2. Add to your `settings.json`:

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

3. Configure MCP tools via the **gear icon** in the Chat tab.

### Cursor IDE

1. Go to: **Settings > MCP > Add New Global MCP Server**
2. Paste the JSON configuration and tools are ready to use within the chat.

---

## Server Modes and Endpoints

### 1. Standard I/O Mode (Default)

Used for direct integration with tools like Claude Desktop and VS Code.

### 2. Streamable HTTP Mode

When running in HTTP mode (set via `SERVER_TYPE=streamablehttp`), the server exposes:

| Endpoint | Description |
|----------|-------------|
| `/stream` | Main MCP protocol endpoint |
| `/healthz` | Server health status endpoint |

Default port: `8081`

---

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_TYPE` | Server mode (`stdio` or `streamablehttp`) | `stdio` | No |
| `STREAM_SERVER_ADDR` | HTTP server address | `:8081` | No |
| `TAZAPAY_API_KEY` | Your Tazapay API key | - | Yes |
| `TAZAPAY_API_SECRET` | Your Tazapay API secret | - | Yes |

---

## Tools Reference

### Balance & FX Tools

#### `tazapay_fetch_fx_tool`

Get real-time foreign exchange rates and convert amounts.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `from_currency` | string | Yes | Source currency (ISO-4217, e.g., USD, EUR) |
| `to_currency` | string | Yes | Target currency (ISO-4217) |
| `amount` | number | Yes | Amount to convert |

**Output:** Exchange rate and converted amount

---

#### `tazapay_fetch_balance_tool`

Fetch your Tazapay account balance.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `currency` | string | No | Filter by currency (ISO-4217). If empty, returns all balances |

**Output:** Account balance(s)

---

### Payment Link / Checkout Tools

#### `generate_payment_link_tool`

Generate a shareable payment link for collecting payments.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `invoice_currency` | string | Yes | Currency for the invoice (ISO-4217) |
| `payment_amount` | number | Yes | Payment amount (decimal format, e.g., 10.50) |
| `customer_name` | string | Yes | Customer's full name |
| `customer_email` | string | Yes | Customer's email address |
| `customer_country` | string | Yes | Customer's country (ISO 3166-1 alpha-2, e.g., US, SG, IN) |
| `transaction_description` | string | Yes | Description of the transaction |

**Output:** Payment link URL and checkout session ID

---

#### `fetch_checkout_tool`

Fetch details of an existing checkout session.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Checkout session ID |

**Output:** Checkout session details including status, amount, and customer info

---

#### `expire_checkout_tool`

Expire a checkout session to prevent further payments.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Checkout session ID to expire |

**Output:** Confirmation of expiration with updated status

---

### Payin Tools

#### `create_payin_tool`

Create a new payin (payment collection) request.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `invoice_currency` | string | Yes | Invoice currency (ISO-4217) |
| `amount` | number | Yes | Payment amount |
| `customer_details` | object | Yes | Customer info (name, email, country, phone) |
| `success_url` | string | Yes | Redirect URL on successful payment |
| `cancel_url` | string | Yes | Redirect URL on cancelled payment |
| `transaction_description` | string | Yes | Transaction description |
| `customer` | string | No | Existing customer ID (cus_xxx) |
| `webhook_url` | string | No | URL for webhook notifications |
| `shipping_details` | object | No | Shipping address details |
| `billing_details` | object | No | Billing address details |
| `reference_id` | string | No | Your internal reference ID |
| `confirm` | boolean | No | Auto-confirm the payin |
| `metadata` | object | No | Custom key-value pairs |

**Output:** Payin ID

---

#### `get_payin_tool`

Fetch details of an existing payin.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Payin ID (must start with `pay_`) |

**Output:** Full payin details including status, amount, and payment method

---

#### `update_payin_tool`

Update an existing payin before confirmation.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Payin ID to update |
| `success_url` | string | Yes | Updated success URL |
| `cancel_url` | string | Yes | Updated cancel URL |
| `payment_method_details` | object | Yes | Payment method configuration |
| `customer_details` | object | No | Updated customer details |
| `shipping_details` | object | No | Updated shipping details |
| `billing_details` | object | No | Updated billing details |
| `reference_id` | string | No | Updated reference ID |

**Output:** Updated payin status

---

#### `cancel_payin_tool`

Cancel an existing payin.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Payin ID to cancel |

**Output:** Cancellation confirmation with status

---

### Payout Tools

#### `create_payout_tool`

Create a new payout to a beneficiary.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `amount` | number | Yes | Payout amount (decimal format) |
| `currency` | string | Yes | Payout currency (ISO-4217) |
| `purpose` | string | Yes | Purpose code (PYR001 to PYR028) |
| `transaction_description` | string | Yes | Description of the payout |
| `beneficiary` | string | No* | Existing beneficiary ID (bnf_xxx) |
| `beneficiary_details` | object | No* | Inline beneficiary details |
| `reference_id` | string | No | Your internal reference ID |
| `holding_currency` | string | No | Currency to fund from |
| `type` | string | No | Payout type: `local`, `swift`, or `wallet` |
| `charge_type` | string | No | For wire: `shared` or `ours` |
| `documents` | array | No | Supporting documents |
| `metadata` | object | No | Custom key-value pairs |

> **Note:** Either `beneficiary` OR `beneficiary_details` is required, but not both.

**Output:** Payout ID

---

#### `get_payout_tool`

Fetch details of an existing payout.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Payout ID (must start with `pot_`) |

**Output:** Full payout details including status, transactions, and beneficiary info

---

#### `fund_payout_tool`

Fund a payout that is in `requires_funding` state.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Payout ID to fund (must start with `pot_`) |

**Output:** Funding confirmation with updated status and amount

---

### Beneficiary Tools

#### `create_beneficiary_tool`

Create a new beneficiary for payouts.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Beneficiary's full legal name |
| `type` | string | Yes | `individual` or `business` |
| `destination_details` | object | Yes | Payment destination (bank, wallet, or local_payment_network) |
| `email` | string | No | Beneficiary's email |
| `phone` | object | No | Phone number (calling_code, number) |
| `address` | object | No | Physical address |
| `tax_id` | string | No | Tax identification number |
| `national_identification_number` | string | No | National ID or passport number |
| `document` | object | No | Supporting document (type, url) |

**Destination Details Structure:**
- `type`: `bank`, `wallet`, or `local_payment_network`
- `bank`: account_number/iban, bank_codes (swift_code, ifsc_code, etc.), country, currency
- `wallet`: deposit_address, type, currency
- `local_payment_network`: type, deposit_key_type, deposit_key

**Output:** Beneficiary ID and destination ID

---

#### `get_beneficiary_tool`

Fetch details of an existing beneficiary.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Beneficiary ID (must start with `bnf_`) |

**Output:** Full beneficiary details

---

#### `update_beneficiary_tool`

Update an existing beneficiary.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Beneficiary ID to update |
| `name` | string | No | Updated name |
| `type` | string | No | Updated type |
| `address` | object | No | Updated address |
| `phone` | object | No | Updated phone |
| `destination_details` | object | No | Updated payment destination |
| `tax_id` | string | No | Updated tax ID |

**Output:** Updated beneficiary details

---

### Customer Tools

#### `tazapay_create_customer_tool`

Create a new customer record.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Customer's name |
| `email` | string | Yes | Customer's email |
| `country` | string | Yes | Country (ISO 3166-1 alpha-2) |
| `reference_id` | string | No | Your internal customer reference |
| `phone` | object | No | Phone details |
| `billing` | array | No | Billing addresses |
| `shipping` | array | No | Shipping addresses |
| `metadata` | object | No | Custom key-value pairs |

**Output:** Customer ID and name

---

#### `tazapay_fetch_customer_tool`

Fetch details of an existing customer.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Customer ID (must start with `cus_`) |

**Output:** Full customer details

---

### Payment Attempt Tools

#### `get_payment_attempt_tool`

Fetch details of a specific payment attempt.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | Yes | Payment attempt ID (must start with `pat_`) |

**Output:** Payment attempt details including status and amount

---

## ID Prefixes Reference

| Entity | Prefix | Example |
|--------|--------|---------|
| Payin | `pay_` | `pay_abc123` |
| Payout | `pot_` | `pot_xyz789` |
| Beneficiary | `bnf_` | `bnf_def456` |
| Customer | `cus_` | `cus_ghi012` |
| Payment Attempt | `pat_` | `pat_jkl345` |
| Checkout | `chk_` | `chk_mno678` |

---

## License

This project is licensed under the MIT license. Refer to [LICENSE](https://github.com/tazapay/tazapay-mcp-server/blob/main/LICENSE) for details.
