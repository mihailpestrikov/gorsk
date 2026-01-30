package server

import (
	"encoding/json"
	"net/http"
)

// Health returns 200 OK with a JSON response indicating the API status.
// This endpoint does not require authentication.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Explicitly marshal and write with a newline to match the test's expectation.
	resp := map[string]string{"status": "ok"}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
	w.Write([]byte("\n")) // Add a newline to satisfy the test
}
