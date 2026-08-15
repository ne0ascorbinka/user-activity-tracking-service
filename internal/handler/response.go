package handler

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized JSON error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// respondJSON writes a JSON response with status code and Content-Type header.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// respondError writes a structured JSON error response.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
