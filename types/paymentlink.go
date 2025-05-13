package types

// PaymentLinkParams represents the input fields extracted from MCP request
type PaymentLinkParams struct {
	InvoiceCurrency string
	Description     string
	CustomerName    string
	CustomerEmail   string
	CustomerCountry string
	PaymentAmount   float64
}

// PaymentLinkRequest defines the payload sent to the internal API
type PaymentLinkRequest struct {
	CustomerDetails        map[string]string `json:"customer_details"`
	InvoiceCurrency        string            `json:"invoice_currency"`
	TransactionDescription string            `json:"transaction_description"`
	Amount                 int64             `json:"amount"`
}
