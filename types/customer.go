package types

// Customer represents the full customer object returned by the API
// This should match the structure of the Tazapay customer response.
type Customer struct {
	Phone     *Phone            `json:"phone,omitempty"`
	Metadata  map[string]any    `json:"metadata"`
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Country   string            `json:"country"`
	CreatedAt string            `json:"created_at"`
	Object    string            `json:"object"`
	Billing   []CustomerContact `json:"billing,omitempty"`
	Shipping  []CustomerContact `json:"shipping,omitempty"`
}

type CustomerContact struct {
	Address *Address `json:"address,omitempty"`
	Phone   *Phone   `json:"phone,omitempty"`
	Name    string   `json:"name"`
}
