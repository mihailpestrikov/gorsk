package server

import (
	"encoding/json"
	"net/http"
)

// Health returns a simple health check status.
// It responds with a 200 OK status and a JSON body {"status": "ok"}.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
