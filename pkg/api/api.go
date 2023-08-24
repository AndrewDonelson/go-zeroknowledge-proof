package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/gorilla/mux"

	gzp "github.com/AndrewDonelson/go-zeroknowledge-proof/pkg/zeroknowledge"
)

// Mock node ID for the current node (in a real-world scenario, this would be dynamically assigned or configured)
const currentNodeID = "node_123"

// Define a struct for the broadcast payload
type BroadcastPayload struct {
	Data               string `json:"data"`
	Proof              string `json:"proof"`
	OriginatingNodeURL string `json:"originatingNodeURL"`
}

// Define a struct for the acknowledgment payload
type AcknowledgmentPayload struct {
	DataID string `json:"dataId"`
	NodeID string `json:"nodeId"`
	Status string `json:"status"`
}

// authorizedNodes is a List of authorized nodes (for simplicity, using hardcoded URLs; in a real-world scenario, use a dynamic registry or DB)
var authorizedNodes = []string{
	"https://node1.example.com",
	"https://node2.example.com",
	// ... add other node URLs here ...
}

// nodeDataStorage is a struct to represent the data storage (for simplicity, using a map; in a real-world scenario, you'd use a database or other persistent storage)
var nodeDataStorage = map[string]string{} // key: data, value: proof

// server is the HTTP/2 server.
var server *http.Server

func RunAPI(port string, certFolder string) {
	// Start the HTTP/2 server
	startServer(port, certFolder)

	// Wait for shutdown signal
	waitForShutdown()

	// Cleanup
	cleanup()
}

// startServer starts the HTTP/2 server in a goroutine.
func startServer(port string, certFolder string) {
	log.Println("Configuring router")
	r := mux.NewRouter()

	// Heartbeat Endpoint (/):
	// This endpoint is designed to check the health of the service.
	// Details:
	// - Purpose: To check the health of the service.
	// - Functionality: The endpoint returns a simple "OK" response.
	// - Output: A string indicating that the service is alive.
	// - Use Case: This endpoint can be used to check the health of the service.
	//
	// TODO: Add Version Information, Database Connection Status, Persistent Storage Status, etc.
	r.HandleFunc("/", HeartbeatHandler).Methods("GET")

	// Prove Endpoint (/prove):
	// This endpoint is designed to provide a cryptographic hash of a piece of data to demonstrate its possession without revealing the actual data.
	// Details:
	// - Purpose: When a client wants to prove that they have a specific piece of data without disclosing the data itself, they can use this endpoint.
	// - Functionality: The endpoint takes a piece of data (in our example, it's hardcoded as "secret_data" for simplicity) and computes its SHA-256 hash. This hash is then returned to the client.
	// - Output: A string representing the SHA-256 hash of the data.
	// - Use Case: Think of this as a "show of proof" without exposure. For instance, you might want to prove you have a specific document or secret, but you don't want to disclose its content.
	r.HandleFunc("/prove", ProveHandler).Methods("POST")

	// This endpoint is designed to compare two cryptographic hashes and determine if they match, effectively verifying if two pieces of data are the same without exposing the data itself.
	// Details:
	// - Purpose: To verify that a given hash matches an expected hash, confirming the authenticity or correctness of data without seeing the original data.
	// - Functionality: The endpoint takes in two hashes: providedHash (the hash received from a client or another source) and expectedHash (the known correct hash). It then compares these hashes for equality.
	// - Output: A boolean value indicating whether the two hashes match (true) or not (false).
	// - Use Case: Suppose a client has previously used the "prove" endpoint and provided their hash to another party. Later, they can use this hash in the "verify" endpoint to confirm that the data they have matches the expected data.
	r.HandleFunc("/verify", VerifyHandler).Methods("POST")

	// This endpoint is designed to broadcast a piece of data to all authorized nodes.
	// Details:
	// - Purpose: To broadcast a piece of data to all authorized nodes.
	// - Functionality: The endpoint takes in a piece of data and its proof (hash) and broadcasts it to all authorized nodes.
	// - Output: A string indicating that the broadcast has been initiated.
	// - Use Case: Suppose a client has a piece of data that they want to broadcast to all authorized nodes. They can use this endpoint to do so.
	r.HandleFunc("/broadcast", BroadcastHandler).Methods("POST")

	// This endpoint is designed to receive a piece of data from another node.
	// Details:
	// - Purpose: To receive a piece of data from another node.
	// - Functionality: The endpoint takes in a piece of data and its proof (hash) and stores it in the node's data storage.
	// - Output: A string indicating that the data has been accepted and stored.
	// - Use Case: Suppose a node receives a piece of data from another node. They can use this endpoint to store the data in their data storage.
	r.HandleFunc("/receive", ReceiveHandler).Methods("POST")

	// This endpoint is designed to receive an acknowledgment from another node.
	// Details:
	// - Purpose: To receive an acknowledgment from another node.
	// - Functionality: The endpoint takes in an acknowledgment payload and stores it or processes it as needed.
	// - Output: A string indicating that the acknowledgment has been received.
	// - Use Case: Suppose a node sends a piece of data to another node. The receiving node can use this endpoint to receive an acknowledgment from the other node.
	r.HandleFunc("/acknowledge", AcknowledgeHandler).Methods("POST")

	server = &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Start the server in a goroutine so that it doesn't block.
	go func() {
		// Usage in your server setup:
		log.Println("Loading TLS certificates")
		certPath, keyPath := getCertPaths(certFolder)

		log.Printf("http/2 API Server listening on %s\n", server.Addr)
		if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
			log.Fatal(err)
		}
	}()
}

// stopServer stops the HTTP/2 server.
func stopServer() {
	//if we have a server, stop it
	if server != nil {
		log.Default().Println("Stopping server")
		server.Close()
	}
}

// HANDLER FUNCTIONS (PUBLIC)

// HeartbeatHandler checks the health of the service. For simplicity, we'll just send an "OK" response.
func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service is alive!")
}

// ProveHandler handles the /prove endpoint.
func ProveHandler(w http.ResponseWriter, r *http.Request) {
	// Here, you'd extract the data from the request, for simplicity, we'll use a static value.
	data := gzp.Data{Content: "secret_data"}
	hashToSend := gzp.Prove(data)

	w.Write([]byte(hashToSend))
}

// VerifyHandler handles the /verify endpoint.
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Extract providedHash and expectedHash from the request
	// For simplicity, we'll use static values.
	providedHash := "some_hash_value"
	expectedHash := "some_expected_hash_value"

	isVerified := gzp.Verify(providedHash, expectedHash)

	w.Write([]byte(fmt.Sprintf("Verification result: %v", isVerified)))
}

// BroadcastHandler handles the /broadcast endpoint
func BroadcastHandler(w http.ResponseWriter, r *http.Request) {
	var payload BroadcastPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Compute the hash (proof) of the data
	computedProof := gzp.Prove(gzp.Data{Content: payload.Data})
	if computedProof != payload.Proof {
		http.Error(w, "Proof verification failed", http.StatusUnauthorized)
		return
	}

	// Broadcast the data to all authorized nodes
	for _, nodeURL := range authorizedNodes {
		go broadcastToNode(nodeURL, payload) // broadcasting in goroutines for concurrency
	}

	w.Write([]byte("Broadcast initiated"))
}

// ReceiveHandler handles the data broadcasted from other nodes
func ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	var payload BroadcastPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Step 1: Compute the hash (proof) of the received data
	computedProof := gzp.Prove(gzp.Data{Content: payload.Data})

	// Step 2: Compare the computed hash with the received hash
	isVerified := gzp.Verify(computedProof, payload.Proof)
	if !isVerified {
		http.Error(w, "Data verification failed", http.StatusUnauthorized)
		return
	}

	// Step 3: If the hashes match, store the data
	nodeDataStorage[payload.Data] = computedProof

	// Step 4: Send acknowledgment back to the originating node (assuming we know the originating node's URL; for simplicity, using a hardcoded one)
	//w.Write([]byte("Data accepted and stored"))
	// If data is verified
	if isVerified {
		// Store data
		nodeDataStorage[payload.Data] = computedProof

		// Send acknowledgment back to the originating node
		go sendAcknowledgmentToNode(payload.OriginatingNodeURL, computedProof, "accepted")
		w.Write([]byte("Data accepted and stored"))
	} else {
		// Optionally, send a rejection acknowledgment
		go sendAcknowledgmentToNode(payload.OriginatingNodeURL, computedProof, "rejected")
		http.Error(w, "Data verification failed", http.StatusUnauthorized)
	}
}

// AcknowledgeHandler handles the /acknowledge endpoint for other nodes to acknowledge data acceptance
func AcknowledgeHandler(w http.ResponseWriter, r *http.Request) {
	var payload AcknowledgmentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Store or process the acknowledgment as needed. For now, we'll just log it.
	log.Printf("Received acknowledgment from Node %s regarding data %s with status: %s", payload.NodeID, payload.DataID, payload.Status)

	// Respond with a simple acknowledgment receipt message
	w.Write([]byte("Acknowledgment received"))
}

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
	stopServer()
	log.Println("Server stopped")
}
