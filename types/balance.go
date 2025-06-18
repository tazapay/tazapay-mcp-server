package types

type BalanceRequest struct {
	Currency string `json:"currency,omitempty"` // Optional currency filter
}

type BalanceResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    BalanceDataBlock `json:"data"`
}

type BalanceDataBlock struct {
	Object    string    `json:"object"`
	UpdatedAt string    `json:"updated_at"`
	Available []Balance `json:"available"`
}

type Balance struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}
