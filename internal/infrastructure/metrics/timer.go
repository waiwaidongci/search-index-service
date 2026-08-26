// Package implementation for tenant-isolated indexing and full-text search.
package metrics

import (
	"sync"
	"time"
)

type Timer struct {
	mu    sync.Mutex
	count uint64
	total time.Duration
	max   time.Duration
}

func (t *Timer) Observe(d time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.count++
	t.total += d
	if d > t.max {
		t.max = d
	}
}
func (t *Timer) Snapshot() (uint64, time.Duration, time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.count, t.total, t.max
}
func (t *Timer) Average() time.Duration {
	c, total, _ := t.Snapshot()
	if c == 0 {
		return 0
	}
	return total / time.Duration(c)
}
