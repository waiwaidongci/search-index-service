// Package implementation for tenant-isolated indexing and full-text search.
package httpadapter

import (
	"net/http"
	"strings"
)

func setJSONHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
}
func checkIfMatch(r *http.Request, current string) bool {
	v := r.Header.Get("If-Match")
	return v == "" || strings.Trim(v, "\"") == current
}
func etagValue(value string) string { return `"` + value + `"` }
