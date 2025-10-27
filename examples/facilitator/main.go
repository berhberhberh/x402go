package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/berhberhberh/x402go"
)

// MockFacilitator is a simple facilitator implementation for demonstration
type MockFacilitator struct{}

// Verify checks if a payment transaction is valid
func (f *MockFacilitator) Verify(req *x402go.VerifyRequest) (*x402go.VerifyResponse, error) {
	fmt.Printf("Verifying transaction: %s on chain %s\n", req.TxHash, req.Chain)

	// In a real implementation, you would:
	// 1. Connect to the blockchain RPC
	// 2. Fetch the transaction details
	// 3. Verify the transaction exists and is confirmed
	// 4. Extract payment details (amount, token, sender, recipient)

	// For demo purposes, we'll simulate successful verification
	return &x402go.VerifyResponse{
		Valid:     true,
		TxHash:    req.TxHash,
		Chain:     req.Chain,
		Token:     "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Amount:    "1000000",
		Sender:    "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		Recipient: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
	}, nil
}

// Settle processes and settles a payment
func (f *MockFacilitator) Settle(req *x402go.SettleRequest) (*x402go.SettleResponse, error) {
	fmt.Printf("Settling payment: %s\n", req.Payment.TxHash)

	// In a real implementation, you would:
	// 1. Verify the payment hasn't been settled already
	// 2. Update your settlement database
	// 3. Potentially trigger additional actions (notifications, accounting, etc.)

	return &x402go.SettleResponse{
		Settled:   true,
		TxHash:    req.Payment.TxHash,
		Timestamp: time.Now().Unix(),
	}, nil
}

func main() {
	// Create facilitator
	facilitator := &MockFacilitator{}

	// Create facilitator server
	server := x402go.NewFacilitatorServer(facilitator)

	fmt.Println("Facilitator server starting on :8081")
	fmt.Println("Endpoints:")
	fmt.Println("  POST /verify - Verify a transaction")
	fmt.Println("  POST /settle - Settle a payment")
	fmt.Println("  GET  /health - Health check")

	log.Fatal(http.ListenAndServe(":8081", server))
}
