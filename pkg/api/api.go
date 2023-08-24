package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

type APIServer struct {
	port       string
	certFolder string
	server     *http.Server
	router     *mux.Router
}

var (
	Server *APIServer
)

func init() {
	Server = &APIServer{}
}

func AddRoute(path string, handler func(http.ResponseWriter, *http.Request)) {
}

// startServer starts the HTTP/2 server in a goroutine.
func (s *APIServer) Initialize(port string, certFolder string) {
	s.port = port
	s.certFolder = certFolder

	log.Println("Configuring router")
	s.router = mux.NewRouter()

	// Heartbeat Endpoint (/):
	// This endpoint is designed to check the health of the service.
	// Details:
	// - Purpose: To check the health of the service.
	// - Functionality: The endpoint returns a simple "OK" response.
	// - Output: A string indicating that the service is alive.
	// - Use Case: This endpoint can be used to check the health of the service.
	//
	// TODO: Add Version Information, Database Connection Status, Persistent Storage Status, etc.
	s.router.HandleFunc("/", HeartbeatHandler).Methods("GET")

	s.server = &http.Server{
		Addr:    port,
		Handler: s.router,
	}

}

func (s *APIServer) Start() {

	defer cleanup()

	// Start the server in a goroutine so that it doesn't block.
	go func() {
		// Usage in your server setup:
		log.Println("Loading TLS certificates")
		certPath, keyPath := getCertPaths(s.certFolder)

		log.Printf("http/2 API Server listening on %s\n", s.server.Addr)
		if err := s.server.ListenAndServeTLS(certPath, keyPath); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown()
}

// stopServer stops the HTTP/2 server.
func (s *APIServer) Stop() {
	//if we have a server, stop it
	if s.server != nil {
		log.Default().Println("Stopping server")
		s.server.Close()
	}
}
