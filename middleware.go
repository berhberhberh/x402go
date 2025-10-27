package x402go

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

// PaymentContextKey is the context key for storing payment information
type contextKey string

const (
	paymentContextKey contextKey = "x402_payment"
)

// PaymentVerifier is an interface for verifying payments
type PaymentVerifier interface {
	Verify(payment *Payment, requirements *PaymentRequirements) (bool, error)
}

// MiddlewareConfig holds configuration for payment middleware
type MiddlewareConfig struct {
	// Requirements defines the payment requirements
	Requirements *PaymentRequirements

	// Verifier is used to verify payments (optional, uses default if nil)
	Verifier PaymentVerifier

	// OnPaymentVerified is called when a payment is successfully verified
	OnPaymentVerified func(payment *Payment, r *http.Request)

	// OnPaymentRequired is called when payment is required but not provided
	OnPaymentRequired func(w http.ResponseWriter, r *http.Request)

	// NonceGenerator generates unique nonces (optional)
	NonceGenerator func() string

	// ExpiryDuration sets how long payment requirements are valid (default: 5 minutes)
	ExpiryDuration time.Duration
}

// RequirePayment creates HTTP middleware that requires payment before processing requests
func RequirePayment(requirements *PaymentRequirements, next http.Handler) http.Handler {
	config := &MiddlewareConfig{
		Requirements: requirements,
	}
	return RequirePaymentWithConfig(config, next)
}

// RequirePaymentWithConfig creates HTTP middleware with custom configuration
func RequirePaymentWithConfig(config *MiddlewareConfig, next http.Handler) http.Handler {
	if config.Requirements == nil {
		panic("payment requirements cannot be nil")
	}

	// Set defaults
	if config.NonceGenerator == nil {
		config.NonceGenerator = generateNonce
	}
	if config.ExpiryDuration == 0 {
		config.ExpiryDuration = 5 * time.Minute
	}
	if config.Verifier == nil {
		config.Verifier = &DefaultVerifier{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if payment response header is present
		paymentHeader := r.Header.Get(HeaderPaymentResponse)

		if paymentHeader == "" {
			// No payment provided, return 402 with requirements
			sendPaymentRequired(w, config)
			return
		}

		// Parse payment from header
		var payment Payment
		if err := json.Unmarshal([]byte(paymentHeader), &payment); err != nil {
			http.Error(w, "Invalid payment format", http.StatusBadRequest)
			return
		}

		// Verify payment
		valid, err := config.Verifier.Verify(&payment, config.Requirements)
		if err != nil {
			http.Error(w, "Payment verification failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		if !valid {
			http.Error(w, "Payment verification failed", http.StatusPaymentRequired)
			return
		}

		// Payment verified, call callback if provided
		if config.OnPaymentVerified != nil {
			config.OnPaymentVerified(&payment, r)
		}

		// Add payment to context
		ctx := context.WithValue(r.Context(), paymentContextKey, &PaymentContext{
			Payment:    payment,
			Verified:   true,
			VerifiedAt: time.Now(),
		})

		// Continue to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// sendPaymentRequired sends a 402 Payment Required response with payment requirements
func sendPaymentRequired(w http.ResponseWriter, config *MiddlewareConfig) {
	if config.OnPaymentRequired != nil {
		config.OnPaymentRequired(w, nil)
		return
	}

	// Set nonce and expiry if not already set
	requirements := *config.Requirements // Copy
	if requirements.Nonce == "" {
		requirements.Nonce = config.NonceGenerator()
	}
	if requirements.Expiry == 0 {
		requirements.Expiry = time.Now().Add(config.ExpiryDuration).Unix()
	}

	// Convert requirements to JSON
	reqJSON, err := requirements.ToJSON()
	if err != nil {
		http.Error(w, "Failed to generate payment requirements", http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set(HeaderPayment, reqJSON)
	w.Header().Set(HeaderWWWAuthenticate, "X-Payment")
	w.Header().Set("Content-Type", "application/json")

	// Send 402 response
	w.WriteHeader(http.StatusPaymentRequired)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "Payment Required",
		"payment": requirements,
	})
}

// GetPayment retrieves payment information from request context
func GetPayment(r *http.Request) (*PaymentContext, bool) {
	ctx := r.Context().Value(paymentContextKey)
	if ctx == nil {
		return nil, false
	}
	payment, ok := ctx.(*PaymentContext)
	return payment, ok
}

// generateNonce generates a random nonce
func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// DefaultVerifier is a basic payment verifier (you should implement proper blockchain verification)
type DefaultVerifier struct{}

// Verify implements basic validation (should be enhanced with actual blockchain verification)
func (v *DefaultVerifier) Verify(payment *Payment, requirements *PaymentRequirements) (bool, error) {
	// Basic validation
	if payment.Chain != requirements.Chain {
		return false, nil
	}
	if payment.Token != requirements.Token {
		return false, nil
	}
	if payment.Amount != requirements.Amount {
		return false, nil
	}
	if payment.Recipient != requirements.Recipient {
		return false, nil
	}

	// Note: In production, you should verify the transaction on-chain
	// This is where you would check the blockchain to confirm the transaction exists
	// and has the correct parameters

	return true, nil
}
