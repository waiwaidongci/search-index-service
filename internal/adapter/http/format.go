// Package implementation for tenant-isolated indexing and full-text search.
package httpadapter

import (
	"encoding/json"
	"net/http"
	"time"
)

type Envelope struct {
	Data any            `json:"data"`
	Meta map[string]any `json:"meta,omitempty"`
	At   time.Time      `json:"at"`
}

func writeEnvelope(w http.ResponseWriter, status int, data any, meta map[string]any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{Data: data, Meta: meta, At: time.Now().UTC()})
}
func writeNoContent(w http.ResponseWriter) { w.WriteHeader(http.StatusNoContent) }
