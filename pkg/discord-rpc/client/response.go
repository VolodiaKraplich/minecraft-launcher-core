package client

// Response represents a Discord RPC response
type Response struct {
	Cmd   string `json:"cmd"`
	Data  any    `json:"data"`
	Event string `json:"evt,omitempty"`
	Nonce string `json:"nonce,omitempty"`
}

// HandshakeResponse represents the response to a handshake request
type HandshakeResponse struct {
	Version int    `json:"v"`
	Config  Config `json:"config"`
	User    User   `json:"user"`
}

// Config represents Discord configuration data
type Config struct {
	CDNHost     string `json:"cdn_host"`
	APIEndpoint string `json:"api_endpoint"`
	Environment string `json:"environment"`
}

// User represents Discord user data
// PremiumType represents the user's Discord Nitro subscription type
type PremiumType int

const (
	PremiumNone    PremiumType = 0
	PremiumClassic PremiumType = 1
	PremiumNitro   PremiumType = 2
)

type User struct {
	ID            string      `json:"id"`
	Username      string      `json:"username"`
	Discriminator string      `json:"discriminator"`
	Avatar        string      `json:"avatar"`
	Bot           bool        `json:"bot"`
	Flags         int         `json:"flags"`
	PremiumType   PremiumType `json:"premium_type"`
}
