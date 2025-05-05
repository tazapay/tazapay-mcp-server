# Tazapay MCP Server

## Overview

This repository contains a working prototype of the MCP Server tailored for calling tools integrating LLMs with Tazapay's Payments API. The server interacts with tools that accept structured queries and generates responses based on natural language input.

The architecture is designed to support future enhancements such as:

- Support for various payment methods (e.g., credit card, bank transfer).
- Dynamic prompting via LLM.
- Modular tool integration (e.g., country-currency compatibility).

## Tech Stack

- **Language**: Go
- **Framework**: Custom Plugin Server built for Anthropic MCP by mark3labs
- **Containerization**: Docker
- **Tooling**: LLM-based input handling and prompt chaining
- **Logging**: Structured JSON logging (planned)

## Prerequisites

- Docker installed
- MCP-compatible environment (e.g., Anthropic setup, GitHub Copilot, etc.)
- Basic understanding of Go and Docker

## Getting Started

This section provides a step-by-step guide to set up the MCP server locally for development and testing purposes.

### Clone the Repository

```bash
git clone https://github.com/tazapay/tazapay-mcp-server.git
cd tazapay-mcp-server
```

### Build the Go Application

```bash
go build cmd/server/main.go
```

### Add Environment Variables to Your Home Directory

Add Tazapay API credentials to your home directory. Create a file named `.tazapay-mcp-server.yaml` in your home directory and add the following variables:

```yaml
TAZAPAY_API_SECRET: "YOUR_TAZAPAY_API_SECRET"
TAZAPAY_API_KEY: "YOUR_TAZAPAY_API_KEY"
```

### Add the Server to Your Preferred MCP Environment

#### For Anthropic MCP Claude Desktop

1. Go to **Settings > Developer > Edit Config**.  
   ![Claude Config Page](/assets/readme/claude_config_page.png)
2. Open the config file and add the following JSON configuration:

```json
"Tazapay-mcp-server": {
    "command": "EXACT/PATH/TO/CLONED/REPO/tazapay-mcp-server/main",
    "description": "This MCP server is responsible for handling payment link generation via Tazapay's checkout API. It exposes one tool, 'tazapaymentlinktool', which can be invoked by an agent or external interface to create a shareable payment URL that allows a customer to make a payment in a specific currency.",
    "args": [],
    "tools": [
        {
            "name": "TazaPaymentLinkTool",
            "description": "This tool interacts with the Tazapay checkout API to generate a payment link. The payment link is associated with a specific transaction defined by the provided customer and payment details. This link can be sent to the customer so they can complete the payment through Tazapay's hosted checkout.",
            "parameters": {
                "invoice_currency": {
                    "type": "string",
                    "description": "The 3-letter ISO 4217 currency code in which the payment should be made. Examples include 'USD' for US Dollars, 'INR' for Indian Rupees, and 'SGD' for Singapore Dollars. This determines the currency displayed on the payment page."
                },
                "payment_amount": {
                    "type": "number",
                    "description": "The amount to be paid, represented in the smallest currency unit (e.g., cents for USD, paise for INR). To convert a standard float value like 100.00 to the required integer format, multiply it by 100 (i.e., 100.00 becomes 10000). This ensures precise currency handling without floating point issues."
                },
                "customer_name": {
                    "type": "string",
                    "description": "The full legal name of the customer who will make the payment. This name will appear on the Tazapay payment interface and may be used for invoicing and compliance checks."
                },
                "customer_email": {
                    "type": "string",
                    "description": "The email address of the customer. This is used to send the payment link and any related transaction notifications. Must be a valid email format (e.g., 'user@example.com')."
                },
                "customer_country": {
                    "type": "string",
                    "description": "The 2-letter ISO 3166-1 alpha-2 country code representing the customer's country of residence. Examples include 'US' for the United States, 'IN' for India, 'SG' for Singapore. This is important for compliance and payment routing."
                },
                "transaction_description": {
                    "type": "string",
                    "description": "A short text that describes the purpose of the transaction or what the payment is for. This appears on the payment page and helps both the sender and receiver understand the context (e.g., 'Consulting Fees for April 2025')."
                }
            }
        },
        {
            "name": "TazaFXTool",
            "description": "This tool interacts with the Tazapay FX API to get the latest exchange rates. It allows conversion of an amount from one currency to another based on real-time rates.",
            "parameters": {
                "from_currency": {
                    "type": "string",
                    "description": "The 3-letter ISO 4217 currency code of the currency to convert from. Examples include 'USD' for US Dollars, 'EUR' for Euros, and 'JPY' for Japanese Yen."
                },
                "to_currency": {
                    "type": "string",
                    "description": "The 3-letter ISO 4217 currency code of the currency to convert to. Examples include 'USD' for US Dollars, 'EUR' for Euros, and 'JPY' for Japanese Yen."
                },
                "amount": {
                    "type": "number",
                    "description": "The amount to be converted, represented in the smallest currency unit (e.g., cents for USD, paise for INR). To convert a standard float value like 100.00 to the required integer format, multiply it by 100 (i.e., 100.00 becomes 10000)."
                }
            }
        }
    ]
}
```

#### For GitHub Copilot Chat in VSCode

1. Go to **Settings > Features > Chat > Enable MCP**.
2. In the Chat tab, click on the gear icon, select the tools to be used, and click **OK**.  
   ![Chat Tab Gear Icon](/assets/readme/chat_gears.png)
3. Now you can use the Tazapay MCP server in your GitHub Copilot Chat.  
   ![Sample Execution in Go Server](/assets/readme/sample_copilot_chat.png)

#### For Cursor IDE

1. Go to **Settings > MCP > Add New global MCP Server**.
2. Add the above JSON configuration in the file.
3. In the chat, select **Agent**, and Cursor is ready for running the Tazapay API from the MCP server.

## Usage and Tools

The MCP server exposes two tools:

- **TazapaymentLinkTool**: This tool can be invoked to generate a payment link.  
  It accepts the following parameters:
  - `invoice_currency`: The 3-letter ISO 4217 currency code in which the payment should be made.
  - `payment_amount`: The amount to be paid, represented in the smallest currency unit (e.g., cents for USD, paise for INR).
  - `customer_name`: The full legal name of the customer who will make the payment.
  - `customer_email`: The email address of the customer.
  - `customer_country`: The 2-letter ISO 3166-1 alpha-2 country code representing the customer's country of residence.
  - `transaction_description`: A short text that describes the purpose of the transaction or what the payment is for.

- **TazaFXTool**: This tool interacts with the Tazapay FX API to get the latest exchange rates.  
  It accepts the following parameters:
  - `from_currency`: The 3-letter ISO 4217 currency code of the currency to convert from.
  - `to_currency`: The 3-letter ISO 4217 currency code of the currency to convert to.
  - `amount`: The amount to be converted, represented in the smallest currency unit (e.g., cents for USD, paise for INR).

Thank you for checking out our repository! Happy coding!!!
