package application

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"sync"
	"testing"
)

func TestConcurrentIndexAndQueryDoNotRace(t *testing.T) {
	s := NewSearchService()
	ctx := context.Background()
	if _, err := s.CreateCollection(ctx, domain.Collection{ID: "c", TenantID: "t", Name: "C"}); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, _ = s.Index(ctx, domain.Document{ID: "d" + string(rune('a'+i%26)), TenantID: "t", CollectionID: "c", Version: int64(i + 1), Fields: map[string]any{"body": "text"}}, "")
		}(i)
	}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _ = s.Query(ctx, domain.Query{TenantID: "t", CollectionID: "c", Text: "text", Limit: 20})
		}()
	}
	close(start)
	wg.Wait()
}

func TestQueryReturnsIndependentSnapshot(t *testing.T) {
	s := NewSearchService()
	ctx := context.Background()
	if _, err := s.CreateCollection(ctx, domain.Collection{ID: "c", TenantID: "t", Name: "C"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Index(ctx, domain.Document{ID: "d", TenantID: "t", CollectionID: "c", Version: 1, Fields: map[string]any{"body": "original"}}, ""); err != nil {
		t.Fatal(err)
	}
	first, err := s.Query(ctx, domain.Query{TenantID: "t", CollectionID: "c", Text: "original", Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Hits) != 1 {
		t.Fatalf("expected one hit, got %d", len(first.Hits))
	}
	first.Hits[0].Fields["body"] = "changed"
	second, err := s.Query(ctx, domain.Query{TenantID: "t", CollectionID: "c", Text: "original", Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Hits) != 1 || second.Hits[0].Fields["body"] != "original" {
		t.Fatalf("query returned a shared field map: %#v", second.Hits)
	}
}
