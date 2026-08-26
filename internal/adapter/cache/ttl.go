// Package implementation for tenant-isolated indexing and full-text search.
package cache

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"sync"
	"time"
)

type entry struct {
	revision domain.TemplateRevision
	expires  time.Time
}
type TTL struct {
	mu   sync.RWMutex
	data map[string]entry
	ttl  time.Duration
}

func NewTTL(ttl time.Duration) *TTL {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &TTL{data: map[string]entry{}, ttl: ttl}
}
func (c *TTL) Get(_ context.Context, key string) (*domain.TemplateRevision, bool) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expires) {
		if ok {
			c.Delete(context.Background(), key)
		}
		return nil, false
	}
	v := e.revision
	return &v, true
}
func (c *TTL) Set(_ context.Context, key string, v domain.TemplateRevision) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = entry{revision: v, expires: time.Now().Add(c.ttl)}
}
func (c *TTL) Delete(_ context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
func (c *TTL) Size() int { c.mu.RLock(); defer c.mu.RUnlock(); return len(c.data) }
