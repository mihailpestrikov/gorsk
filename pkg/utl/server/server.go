package server

import (
	"encoding/json"
	"net/http"
)

// Health returns a simple status "ok" for health checks.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Add any other existing server utility functions or structs below this line.
// Example:
// type Server struct {
// 	Logger *zap.Logger
// }
// func NewServer() *Server { return &Server{} }