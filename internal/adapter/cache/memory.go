// Package implementation for tenant-isolated indexing and full-text search.
package cache

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]domain.TemplateRevision
}

func NewMemory() *Memory { return &Memory{data: map[string]domain.TemplateRevision{}} }
func (m *Memory) Get(_ context.Context, k string) (*domain.TemplateRevision, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[k]
	if !ok {
		return nil, false
	}
	v.Rules = domain.CloneRules(v.Rules)
	return &v, true
}
func (m *Memory) Set(_ context.Context, k string, v domain.TemplateRevision) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[k] = v
}
func (m *Memory) Delete(_ context.Context, k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, k)
}
