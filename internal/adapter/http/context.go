// Package implementation for tenant-isolated indexing and full-text search.
package httpadapter

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func withRequestID(r *http.Request, id string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), requestIDKey, id))
}
func requestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}
func requestTimeout(next http.Handler, d time.Duration) http.Handler {
	return http.TimeoutHandler(next, d, "request timeout")
}
