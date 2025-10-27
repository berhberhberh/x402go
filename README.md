# x402go

A streamlined Go infra folder for enabling x402 architecture dynamically and easily, complete payments to agents using the outdated 402-HTTP request type.

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
