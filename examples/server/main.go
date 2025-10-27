package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shdfhasdfdsgkofngd/x402go"
)

func main() {
	// Define payment requirements
	requirements := &x402go.PaymentRequirements{
		Scheme:    x402go.SchemeExact,
		Amount:    "1000000", // 1 USDC (6 decimals)
		Token:     "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Chain:     "8453", // Base
		Recipient: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
	}

	// Create a protected endpoint
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get payment info from context
		payment, ok := x402go.GetPayment(r)
		if !ok {
			http.Error(w, "Payment context not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Premium content!", "payment": {"txHash": "%s", "amount": "%s"}}`,
			payment.Payment.TxHash, payment.Payment.Amount)
	})

	// Wrap with payment middleware
	http.Handle("/premium", x402go.RequirePayment(requirements, protectedHandler))

	// Free endpoint for comparison
	http.HandleFunc("/free", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "Free content!"}`)
	})

	log.Println("Server starting on :8080")
	log.Println("Try: curl http://localhost:8080/free")
	log.Println("Try: curl http://localhost:8080/premium")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
