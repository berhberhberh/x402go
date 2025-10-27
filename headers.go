package x402go

const (
	// HeaderPayment is the header key for payment requirements (server to client)
	HeaderPayment = "X-Payment"

	// HeaderPaymentResponse is the header key for payment payload (client to server)
	HeaderPaymentResponse = "X-Payment-Response"

	// HeaderWWWAuthenticate is used with 402 status code
	HeaderWWWAuthenticate = "WWW-Authenticate"

	// SchemeExact represents the "exact" payment scheme
	SchemeExact = "exact"
)
