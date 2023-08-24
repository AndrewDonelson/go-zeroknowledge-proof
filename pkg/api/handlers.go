package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	gzp "github.com/AndrewDonelson/go-zeroknowledge-proof/pkg/zeroknowledge"
)

// HANDLER FUNCTIONS (PUBLIC) - Add your handlers here

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
