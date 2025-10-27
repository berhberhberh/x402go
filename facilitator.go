package x402go

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Facilitator provides blockchain verification and settlement services
type Facilitator interface {
	// Verify checks if a payment transaction is valid on the blockchain
	Verify(req *VerifyRequest) (*VerifyResponse, error)

	// Settle processes and settles a payment
	Settle(req *SettleRequest) (*SettleResponse, error)
}

// FacilitatorServer wraps a Facilitator with HTTP handlers
type FacilitatorServer struct {
	facilitator Facilitator
	mux         *http.ServeMux
}

// NewFacilitatorServer creates a new facilitator server
func NewFacilitatorServer(facilitator Facilitator) *FacilitatorServer {
	fs := &FacilitatorServer{
		facilitator: facilitator,
		mux:         http.NewServeMux(),
	}

	// Register handlers
	fs.mux.HandleFunc("/verify", fs.handleVerify)
	fs.mux.HandleFunc("/settle", fs.handleSettle)
	fs.mux.HandleFunc("/health", fs.handleHealth)

	return fs
}

// ServeHTTP implements http.Handler
func (fs *FacilitatorServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fs.mux.ServeHTTP(w, r)
}

// handleVerify handles POST /verify requests
func (fs *FacilitatorServer) handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := fs.facilitator.Verify(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleSettle handles POST /settle requests
func (fs *FacilitatorServer) handleSettle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SettleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := fs.facilitator.Settle(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleHealth handles GET /health requests
func (fs *FacilitatorServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// FacilitatorClient is a client for communicating with a facilitator server
type FacilitatorClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewFacilitatorClient creates a new facilitator client
func NewFacilitatorClient(baseURL string) *FacilitatorClient {
	return &FacilitatorClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// Verify verifies a payment transaction
func (fc *FacilitatorClient) Verify(req *VerifyRequest) (*VerifyResponse, error) {
	resp, err := fc.httpClient.Post(
		fc.baseURL+"/verify",
		"application/json",
		toReader(req),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var verifyResp VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, err
	}

	return &verifyResp, nil
}

// Settle settles a payment
func (fc *FacilitatorClient) Settle(req *SettleRequest) (*SettleResponse, error) {
	resp, err := fc.httpClient.Post(
		fc.baseURL+"/settle",
		"application/json",
		toReader(req),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var settleResp SettleResponse
	if err := json.NewDecoder(resp.Body).Decode(&settleResp); err != nil {
		return nil, err
	}

	return &settleResp, nil
}

// toReader converts any value to an io.Reader containing its JSON representation
func toReader(v interface{}) *bytes.Reader {
	data, _ := json.Marshal(v)
	return bytes.NewReader(data)
}
