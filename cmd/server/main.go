package main

import (
	api "github.com/AndrewDonelson/go-zeroknowledge-proof/pkg/api"
)

func main() {
	// Initialize the HTTP/2 server
	api.Server.Initialize(":8080", "../../certs/")

	// Start the HTTP/2 server
	api.Server.Start()
}
