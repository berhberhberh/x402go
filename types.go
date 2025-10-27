package x402go

import (
	"encoding/json"
	"time"
)

// PaymentRequirements defines the payment details required by a server
type PaymentRequirements struct {
	// Scheme defines the payment scheme (e.g., "exact")
	Scheme string `json:"scheme"`

	// Amount is the payment amount in the smallest unit of the token
	Amount string `json:"amount"`

	// Token is the contract address of the payment token
	Token string `json:"token"`

	// Chain is the blockchain network ID
	Chain string `json:"chain"`

	// Recipient is the address that should receive the payment
	Recipient string `json:"recipient"`

	// Nonce is a unique identifier for this payment request
	Nonce string `json:"nonce,omitempty"`

	// Expiry is when this payment requirement expires (Unix timestamp)
	Expiry int64 `json:"expiry,omitempty"`

	// Facilitator is the URL of the facilitator server (optional)
	Facilitator string `json:"facilitator,omitempty"`
}

// ToJSON converts PaymentRequirements to JSON string
func (pr *PaymentRequirements) ToJSON() (string, error) {
	data, err := json.Marshal(pr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Payment represents a payment made by a client
type Payment struct {
	// Scheme is the payment scheme used
	Scheme string `json:"scheme"`

	// TxHash is the blockchain transaction hash
	TxHash string `json:"txHash"`

	// Chain is the blockchain network ID
	Chain string `json:"chain"`

	// Token is the contract address of the payment token
	Token string `json:"token"`

	// Amount is the payment amount in the smallest unit of the token
	Amount string `json:"amount"`

	// Sender is the address that sent the payment
	Sender string `json:"sender"`

	// Recipient is the address that received the payment
	Recipient string `json:"recipient"`

	// Nonce is the nonce from the payment requirement
	Nonce string `json:"nonce,omitempty"`

	// Timestamp is when the payment was made
	Timestamp int64 `json:"timestamp,omitempty"`
}

// ToJSON converts Payment to JSON string
func (p *Payment) ToJSON() (string, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON parses Payment from JSON string
func (p *Payment) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), p)
}

// VerifyRequest represents a request to verify a payment
type VerifyRequest struct {
	TxHash string `json:"txHash"`
	Chain  string `json:"chain"`
}

// VerifyResponse represents a response from payment verification
type VerifyResponse struct {
	Valid     bool   `json:"valid"`
	TxHash    string `json:"txHash"`
	Chain     string `json:"chain"`
	Token     string `json:"token"`
	Amount    string `json:"amount"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Error     string `json:"error,omitempty"`
}

// SettleRequest represents a request to settle a payment
type SettleRequest struct {
	Payment Payment `json:"payment"`
}

// SettleResponse represents a response from payment settlement
type SettleResponse struct {
	Settled   bool   `json:"settled"`
	TxHash    string `json:"txHash,omitempty"`
	Error     string `json:"error,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// PaymentContext holds information about a verified payment
type PaymentContext struct {
	Payment   Payment
	Verified  bool
	VerifiedAt time.Time
}
