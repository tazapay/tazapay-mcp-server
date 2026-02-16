package types

// PayoutRequest represents the payload for creating a payout
// Reuses beneficiary types for nested fields
// Add/modify fields as needed for your payout API

type PayoutRequest struct {
	LogisticsTrackingDetails *LogisticsTrackingDetails `json:"logistics_tracking_details,omitempty"`
	StatementDescription     string                    `json:"statement_description,omitempty"`
	TransactionDescription   string                    `json:"transaction_description,omitempty"`
	ReferenceID              string                    `json:"reference_id"`
	Beneficiary              string                    `json:"beneficiary"`
	Purpose                  string                    `json:"purpose"`
	ChargeType               string                    `json:"charge_type,omitempty"`
	Type                     string                    `json:"type,omitempty"`
	HoldingCurrency          string                    `json:"holding_currency,omitempty"`
	OnBehalfOf               string                    `json:"on_behalf_of,omitempty"`
	Metadata                 string                    `json:"metadata,omitempty"`
	Currency                 string                    `json:"currency"`
	BeneficiaryDetails       Beneficiary               `json:"beneficiary_details"`
	Amount                   int64                     `json:"amount"`
}

// LogisticsTrackingDetails represents logistics tracking info for a payout
// (Add/modify fields as needed)
type LogisticsTrackingDetails struct {
	LogisticsProviderName string `json:"logistics_provider_name,omitempty"`
	LogisticsProviderCode string `json:"logistics_provider_code,omitempty"`
	TrackingNumber        string `json:"tracking_number,omitempty"`
}
