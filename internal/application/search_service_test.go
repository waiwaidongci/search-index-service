package application

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"testing"
)

func TestIdempotencyIsTenantScoped(t *testing.T) {
	s := NewSearchService()
	ctx := context.Background()
	for _, tenant := range []string{"t1", "t2"} {
		if _, err := s.CreateCollection(ctx, domain.Collection{ID: "c", TenantID: tenant, Name: tenant}); err != nil {
			t.Fatal(err)
		}
	}
	first, err := s.Index(ctx, domain.Document{ID: "d1", TenantID: "t1", CollectionID: "c", Version: 1, Fields: map[string]any{"body": "one"}}, "same")
	if err != nil {
		t.Fatal(err)
	}
	second, err := s.Index(ctx, domain.Document{ID: "d2", TenantID: "t2", CollectionID: "c", Version: 1, Fields: map[string]any{"body": "two"}}, "same")
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == second.ID {
		t.Fatalf("idempotency key crossed tenants: %#v", second)
	}
}
func TestBatchDeleteAndTask(t *testing.T) {
	s := NewSearchService()
	ctx := context.Background()
	s.CreateCollection(ctx, domain.Collection{ID: "c", TenantID: "t", Name: "C"})
	docs, err := s.BatchIndex(ctx, []domain.Document{{ID: "1", TenantID: "t", CollectionID: "c", Version: 1}, {ID: "2", TenantID: "t", CollectionID: "c", Version: 1}}, "batch")
	if err != nil || len(docs) != 2 {
		t.Fatalf("batch failed: %#v %v", docs, err)
	}
	task, err := s.Rebuild(ctx, "t", "c")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Task(task["task_id"].(string)); err != nil {
		t.Fatal(err)
	}
	if err = s.Delete(ctx, "t", "c", "1"); err != nil {
		t.Fatal(err)
	}
}
