// Package implementation for tenant-isolated indexing and full-text search.
package metrics

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Registry struct {
	mu       sync.RWMutex
	counters map[string]*uint64
}

func NewRegistry() *Registry { return &Registry{counters: map[string]*uint64{}} }
func (r *Registry) Counter(name string) *uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	v := uint64(0)
	r.counters[name] = &v
	return &v
}
func (r *Registry) Inc(name string) { atomic.AddUint64(r.Counter(name), 1) }
func (r *Registry) Text() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := ""
	for n, c := range r.counters {
		out += fmt.Sprintf("%s %d\n", n, atomic.LoadUint64(c))
	}
	return out
}
