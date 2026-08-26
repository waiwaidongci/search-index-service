package cache

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"testing"
	"time"
)

func TestListTemplateRevisionsReturnsOwnedRules(t *testing.T) {
	c := NewTTL(time.Minute)
	ctx := context.Background()
	c.Set(ctx, "q", domain.TemplateRevision{Number: 1, Value: "one", Rules: []domain.Rule{{ID: "r", Value: "on", Tags: map[string]string{"env": "prod"}}}})
	got, ok := c.Get(ctx, "q")
	if !ok {
		t.Fatal("expected cached revision")
	}
	got.Rules[0].Value = "changed"
	got.Rules[0].Tags["env"] = "staging"
	again, ok := c.Get(ctx, "q")
	if !ok || again.Rules[0].Value != "on" || again.Rules[0].Tags["env"] != "prod" {
		t.Fatalf("cache returned shared revision: %#v", again)
	}
}

func TestExpiredCacheEntryDoesNotAliasCaller(t *testing.T) {
	c := NewTTL(10 * time.Millisecond)
	ctx := context.Background()
	input := domain.TemplateRevision{Number: 1, Value: "one", Rules: []domain.Rule{{ID: "r", Value: "on"}}}
	c.Set(ctx, "q", input)
	input.Rules[0].Value = "changed"
	got, ok := c.Get(ctx, "q")
	if !ok || got.Rules[0].Value != "on" {
		t.Fatalf("cache did not own input revision: %#v", got)
	}
}
