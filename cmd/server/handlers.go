package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
