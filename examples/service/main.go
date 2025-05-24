package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/githonllc/arbiterid"
)

// Server represents the ID generation service
type Server struct {
	node *arbiterid.Node
	port string
}

// GenerateRequest represents the request payload for ID generation
type GenerateRequest struct {
	IDType *int `json:"id_type,omitempty"` // Optional ID type, defaults to 0
	Count  *int `json:"count,omitempty"`   // Optional count, defaults to 1
}

// GenerateResponse represents the response payload for ID generation
type GenerateResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// IDData represents individual ID data in the response
type IDData struct {
	ID       string `json:"id"`        // Base58 encoded ID
	IDInt64  int64  `json:"id_int64"`  // Raw int64 value
	IDBase64 string `json:"id_base64"` // Base64 encoded ID
	IDHex    string `json:"id_hex"`    // Hexadecimal representation
	Type     int    `json:"type"`      // ID type
	Time     string `json:"time"`      // ISO timestamp
	Node     int64  `json:"node"`      // Node ID
	Sequence int64  `json:"sequence"`  // Sequence number
}

// NewServer creates a new ID generation server
func NewServer(nodeID int, port string) (*Server, error) {
	// Use quiet mode for production service
	node, err := arbiterid.NewNode(nodeID,
		arbiterid.WithStrictMonotonicityCheck(true),
		arbiterid.WithQuietMode(true))
	if err != nil {
		return nil, fmt.Errorf("failed to create arbiterid node: %w", err)
	}

	return &Server{
		node: node,
		port: port,
	}, nil
}

// generateHandler handles POST /generate requests
func (s *Server) generateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		s.sendError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// If JSON decode fails, assume it's a simple request with query parameters
		req = GenerateRequest{}
	}

	// Get ID type from JSON body or query parameters
	idType := 0 // default
	if req.IDType != nil {
		idType = *req.IDType
	} else if idTypeStr := r.URL.Query().Get("type"); idTypeStr != "" {
		if parsed, err := strconv.Atoi(idTypeStr); err == nil {
			idType = parsed
		}
	}

	// Get count from JSON body or query parameters
	count := 1 // default
	if req.Count != nil && *req.Count > 0 {
		count = *req.Count
	} else if countStr := r.URL.Query().Get("count"); countStr != "" {
		if parsed, err := strconv.Atoi(countStr); err == nil && parsed > 0 {
			count = parsed
		}
	}

	// Limit count to prevent abuse
	if count > 100 {
		s.sendError(w, http.StatusBadRequest, "Count cannot exceed 100")
		return
	}

	// Validate ID type
	if idType < 0 || idType > 1023 {
		s.sendError(w, http.StatusBadRequest, "ID type must be between 0 and 1023")
		return
	}

	// Generate IDs
	var results []IDData
	for i := 0; i < count; i++ {
		id, err := s.node.Generate(arbiterid.IDType(idType))
		if err != nil {
			s.sendError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate ID: %v", err))
			return
		}

		idType, _, node, seq := id.Components()
		results = append(results, IDData{
			ID:       id.Base58(),
			IDInt64:  id.Int64(),
			IDBase64: id.Base64(),
			IDHex:    fmt.Sprintf("%x", id.Int64()),
			Type:     int(idType),
			Time:     id.TimeISO(),
			Node:     node,
			Sequence: seq,
		})
	}

	// Send response
	var data interface{}
	if count == 1 {
		data = results[0]
	} else {
		data = results
	}

	s.sendSuccess(w, data)
}

// healthHandler handles GET /health requests
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	// Generate a test ID to verify the service is working
	testID, err := s.node.Generate(0)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Service unhealthy: failed to generate test ID")
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"node_id":   testID.Node(),
		"last_id":   s.node.LastID().String(),
	}

	s.sendSuccess(w, response)
}

// infoHandler handles GET /info requests
func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		s.sendError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	response := map[string]interface{}{
		"service":     "ArbiterID Generation Service",
		"version":     "1.0.0",
		"description": "Distributed unique ID generation service using Snowflake-inspired algorithm",
		"node_id":     s.node.LastID().Node(),
		"epoch":       "2025-01-01T00:00:00.000Z",
		"bit_layout": map[string]interface{}{
			"type":      "10 bits (0-1023)",
			"timestamp": "41 bits (milliseconds since epoch)",
			"node":      "2 bits (0-3)",
			"sequence":  "10 bits (0-1023)",
		},
		"endpoints": map[string]string{
			"POST /generate": "Generate new ID(s)",
			"GET /health":    "Health check",
			"GET /info":      "Service information",
		},
	}

	s.sendSuccess(w, response)
}

// sendSuccess sends a successful JSON response
func (s *Server) sendSuccess(w http.ResponseWriter, data interface{}) {
	response := GenerateResponse{
		Success: true,
		Data:    data,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// sendError sends an error JSON response
func (s *Server) sendError(w http.ResponseWriter, statusCode int, message string) {
	response := GenerateResponse{
		Success: false,
		Error:   message,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// setupRoutes sets up HTTP routes
func (s *Server) setupRoutes() {
	http.HandleFunc("/generate", s.generateHandler)
	http.HandleFunc("/health", s.healthHandler)
	http.HandleFunc("/info", s.infoHandler)

	// Root handler provides basic info
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			s.sendError(w, http.StatusNotFound, "Endpoint not found")
			return
		}
		s.infoHandler(w, r)
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.setupRoutes()

	log.Printf("Starting ArbiterID service on port %s", s.port)
	log.Printf("Node ID: %d", s.node.LastID().Node())
	log.Println("Available endpoints:")
	log.Println("  POST /generate - Generate new ID(s)")
	log.Println("  GET  /health   - Health check")
	log.Println("  GET  /info     - Service information")
	log.Println("  GET  /         - Service information")

	return http.ListenAndServe(":"+s.port, nil)
}

func main() {
	// Get configuration from environment variables with defaults
	nodeIDStr := os.Getenv("NODE_ID")
	if nodeIDStr == "" {
		nodeIDStr = "0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Parse node ID
	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil || nodeID < 0 || nodeID > 3 {
		log.Fatalf("Invalid NODE_ID: %s (must be 0-3)", nodeIDStr)
	}

	// Create and start server
	server, err := NewServer(nodeID, port)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Graceful shutdown would be nice, but keeping this simple for the example
	log.Fatal(server.Start())
}
