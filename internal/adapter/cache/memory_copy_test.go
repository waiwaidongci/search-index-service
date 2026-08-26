package cache

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"testing"
)

func TestCacheMemoryReturnsIndependentRevision(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	m.Set(ctx, "q", domain.TemplateRevision{Number: 1, Value: "one", Rules: []domain.Rule{{ID: "r", Value: "on"}}})
	got, ok := m.Get(ctx, "q")
	if !ok {
		t.Fatal("expected cached revision")
	}
	got.Rules[0].Value = "changed"
	again, ok := m.Get(ctx, "q")
	if !ok || again.Rules[0].Value != "on" {
		t.Fatalf("cache returned shared revision: %#v", again)
	}
}
