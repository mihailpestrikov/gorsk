package server

import (
	"encoding/json"
	"net/http"
	// Potentially other imports like log, config, etc.
)

// Server provides common server-level utilities and handlers.
// It can hold configuration, logger, etc. if needed for other handlers.
type Server struct {
	// Example: config *config.Config
	// Example: logger *zlog.Logger
}

// NewServer creates a new Server instance.
// It initializes the Server with its dependencies.
func NewServer() *Server {
	return &Server{}
}

// Health is an HTTP handler function that responds with a 200 OK and
// a JSON payload {"status": "ok"}.
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Using json.NewEncoder for efficiency and error handling
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		// In a real application, you might log this error.
		// For a simple health check, it's often omitted for brevity unless detailed logging is critical.
	}
}
