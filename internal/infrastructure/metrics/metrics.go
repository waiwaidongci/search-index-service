// Package implementation for tenant-isolated indexing and full-text search.
package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Metrics struct{ requests uint64 }

func New() *Metrics     { return &Metrics{} }
func (m *Metrics) Inc() { atomic.AddUint64(&m.requests, 1) }
func (m *Metrics) Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "search_index_requests_total %d\n", atomic.LoadUint64(&m.requests))
}
