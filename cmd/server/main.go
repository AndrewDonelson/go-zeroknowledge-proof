package main

import (
	api "github.com/AndrewDonelson/go-http2-api/pkg"
)

func main() {
	// Initialize the HTTP/2 server
	api.Server.Initialize(":8080", "../certs/")

	// Add routes here - https://localhost:8080/

	// Zero Proof Prove Route: https://localhost:8080/v1/zeroproof/prove
	api.Server.AddRoute(&api.APIRoute{Version: 1, SubRoute: "zeroproof", Name: "prove", Method: "GET", Handler: ProveHandler})

	// Zero Proof Verify Route: https://localhost:8080/v1/zeroproof/verify
	api.Server.AddRoute(&api.APIRoute{Version: 1, SubRoute: "zeroproof", Name: "verify", Method: "GET", Handler: VerifyHandler})

	// Zero Proof Broadcast Route: https://localhost:8080/v1/zeroproof/broadcast
	api.Server.AddRoute(&api.APIRoute{Version: 1, SubRoute: "zeroproof", Name: "broadcast", Method: "POST", Handler: BroadcastHandler})

	// Zero Proof Acknowledge Route: https://localhost:8080/v1/zeroproof/acknowledge
	api.Server.AddRoute(&api.APIRoute{Version: 1, SubRoute: "zeroproof", Name: "acknowledge", Method: "POST", Handler: AcknowledgeHandler})

	// Start the HTTP/2 server
	api.Server.Start()
}
