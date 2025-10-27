package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/berhberhberh/x402go"
)

func main() {
	// Create a payment handler that simulates making a blockchain transaction
	paymentHandler := func(requirements *x402go.PaymentRequirements) (*x402go.Payment, error) {
		fmt.Printf("Payment Required:\n")
		fmt.Printf("  Amount: %s\n", requirements.Amount)
		fmt.Printf("  Token: %s\n", requirements.Token)
		fmt.Printf("  Chain: %s\n", requirements.Chain)
		fmt.Printf("  Recipient: %s\n", requirements.Recipient)
		fmt.Printf("\nProcessing payment...\n")

		// In a real implementation, you would:
		// 1. Connect to a wallet
		// 2. Sign a transaction
		// 3. Submit to blockchain
		// 4. Wait for confirmation
		// 5. Return the transaction hash

		// For demo purposes, we'll simulate a transaction
		simulatedTxHash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"

		payment := &x402go.Payment{
			Scheme:    requirements.Scheme,
			TxHash:    simulatedTxHash,
			Chain:     requirements.Chain,
			Token:     requirements.Token,
			Amount:    requirements.Amount,
			Sender:    "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
			Recipient: requirements.Recipient,
			Nonce:     requirements.Nonce,
			Timestamp: time.Now().Unix(),
		}

		fmt.Printf("Payment sent! TxHash: %s\n\n", simulatedTxHash)
		return payment, nil
	}

	// Create client with payment handler
	client := x402go.NewClientWithHandler(paymentHandler)

	// Make a request to a protected endpoint
	fmt.Println("Requesting protected resource...")
	resp, err := client.Get("http://localhost:8080/premium")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(body))
}
