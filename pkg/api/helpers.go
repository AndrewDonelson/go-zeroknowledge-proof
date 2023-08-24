package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

// HELPER FUNCTIONS (PRIVATE)

// Helper function to broadcast data to a specific node
func broadcastToNode(nodeURL string, payload BroadcastPayload) {
	data, _ := json.Marshal(payload)
	resp, err := http.Post(nodeURL+"/receive", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to broadcast to node %s: %v", nodeURL, err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Node %s response: %s", nodeURL, body)
}

// Helper function to send acknowledgment to a specific node
func sendAcknowledgmentToNode(nodeURL, dataID, status string) {
	payload := AcknowledgmentPayload{
		DataID: dataID,
		NodeID: currentNodeID,
		Status: status,
	}
	data, _ := json.Marshal(payload)
	resp, err := http.Post(nodeURL+"/acknowledge", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to send acknowledgment to node %s: %v", nodeURL, err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Node %s acknowledgment response: %s", nodeURL, body)
}

// getCertPaths returns the paths to the TLS certificates.
func getCertPaths(certFolder string) (string, string) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	certPath := filepath.Join(basepath, certFolder+"cert.pem")
	keyPath := filepath.Join(basepath, certFolder+"key.pem")
	return certPath, keyPath
}

func waitForShutdown() {
	// Create a channel to listen for OS signals
	signals := make(chan os.Signal, 1)

	// Notify the signals channel for SIGINT (CTRL+C) and SIGTERM (termination request sent to the program)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive the signal
	<-signals

	log.Println("Shutdown signal received")
}

// cleanup closes the server when the application exits.
func cleanup() {
	Server.Stop()
	log.Println("Server stopped")
}
