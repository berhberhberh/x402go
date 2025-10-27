package x402go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is an HTTP client that automatically handles x402 payments
type Client struct {
	httpClient *http.Client

	// PaymentHandler is called when payment is required
	// It should return a Payment object with the transaction details
	PaymentHandler func(requirements *PaymentRequirements) (*Payment, error)

	// MaxRetries is the maximum number of payment retry attempts
	MaxRetries int
}

// NewClient creates a new x402 client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		MaxRetries: 1,
	}
}

// NewClientWithHandler creates a client with a custom payment handler
func NewClientWithHandler(handler func(*PaymentRequirements) (*Payment, error)) *Client {
	return &Client{
		httpClient:     &http.Client{},
		PaymentHandler: handler,
		MaxRetries:     1,
	}
}

// Do executes an HTTP request and handles 402 payment requirements
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// If not 402, return response as-is
	if resp.StatusCode != http.StatusPaymentRequired {
		return resp, nil
	}

	// Handle payment requirement
	return c.handlePaymentRequired(req, resp)
}

// Get performs a GET request
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post performs a POST request
func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

// handlePaymentRequired processes a 402 response and retries with payment
func (c *Client) handlePaymentRequired(originalReq *http.Request, paymentResp *http.Response) (*http.Response, error) {
	defer paymentResp.Body.Close()

	// Extract payment requirements from header
	paymentHeader := paymentResp.Header.Get(HeaderPayment)
	if paymentHeader == "" {
		return nil, fmt.Errorf("402 response missing %s header", HeaderPayment)
	}

	// Parse payment requirements
	var requirements PaymentRequirements
	if err := json.Unmarshal([]byte(paymentHeader), &requirements); err != nil {
		return nil, fmt.Errorf("failed to parse payment requirements: %w", err)
	}

	// Check if we have a payment handler
	if c.PaymentHandler == nil {
		return nil, fmt.Errorf("payment required but no payment handler configured")
	}

	// Call payment handler to make payment
	payment, err := c.PaymentHandler(&requirements)
	if err != nil {
		return nil, fmt.Errorf("payment handler failed: %w", err)
	}

	// Convert payment to JSON
	paymentJSON, err := payment.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to encode payment: %w", err)
	}

	// Clone the original request
	retryReq, err := cloneRequest(originalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to clone request: %w", err)
	}

	// Add payment response header
	retryReq.Header.Set(HeaderPaymentResponse, paymentJSON)

	// Retry request with payment
	return c.httpClient.Do(retryReq)
}

// cloneRequest creates a copy of an HTTP request
func cloneRequest(req *http.Request) (*http.Request, error) {
	clone := req.Clone(req.Context())

	// If original request had a body, we need to copy it
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		// Restore original body
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Set clone body
		clone.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return clone, nil
}
