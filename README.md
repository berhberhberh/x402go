# x402go

A Go implementation of the x402 payment protocol for accepting digital payments over HTTP.

## What is x402?

x402 is an open-standard payment protocol built on HTTP that enables accepting digital payments with minimal friction. It uses the HTTP 402 "Payment Required" status code to implement a standardized payment flow.

## Features

- 🚀 Simple HTTP middleware for requiring payments
- 💰 Support for blockchain-based micropayments
- 🔌 Chain and token agnostic
- 🛡️ Trust-minimizing design
- 📦 Easy integration with existing Go HTTP servers

## Installation

```bash
go get github.com/berhberhberh/x402go
```

## Quick Start

### Server Example

```go
package main

import (
    "net/http"
    "github.com/berhberhberh/x402go"
)

func main() {
    // Create payment requirements
    requirements := &x402go.PaymentRequirements{
        Scheme: "exact",
        Amount: "1000000", // 1 USDC (6 decimals)
        Token: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", // USDC
        Chain: "8453", // Base
        Recipient: "0xYourAddress",
    }

    // Wrap your handler with payment middleware
    handler := x402go.RequirePayment(requirements, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Premium content!"))
    }))

    http.ListenAndServe(":8080", handler)
}
```

### Client Example

```go
package main

import (
    "net/http"
    "github.com/berhberhberh/x402go"
)

func main() {
    client := x402go.NewClient()

    resp, err := client.Get("http://localhost:8080/premium")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Client automatically handles 402 responses and makes payment
}
```

## License

MIT
