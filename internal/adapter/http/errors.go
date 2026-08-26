// Package implementation for tenant-isolated indexing and full-text search.
package httpadapter

import (
	"encoding/json"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	RequestID string    `json:"request_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func respondError(w http.ResponseWriter, status int, code, message, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Code: code, Message: message, RequestID: requestID, Timestamp: time.Now().UTC()})
}
func methodAllowed(w http.ResponseWriter, methods ...string) {
	w.Header().Set("Allow", joinMethods(methods))
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}
func joinMethods(methods []string) string {
	out := ""
	for i, m := range methods {
		if i > 0 {
			out += ", "
		}
		out += m
	}
	return out
}
