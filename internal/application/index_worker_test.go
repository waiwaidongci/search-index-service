package application

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"testing"
)

type fakeChangeQueue struct {
	docs  []domain.Document
	dead  []domain.Document
	errs  []error
}

func (q *fakeChangeQueue) Consume(ctx context.Context) (<-chan domain.Document, error) {
	ch := make(chan domain.Document)
	go func() {
		defer close(ch)
		for _, d := range q.docs {
			select {
			case <-ctx.Done():
				return
			case ch <- d:
			}
		}
	}()
	return ch, nil
}

func (q *fakeChangeQueue) DeadLetter(_ context.Context, d domain.Document, err error) error {
	q.dead = append(q.dead, d)
	q.errs = append(q.errs, err)
	return nil
}

func TestWorkerResetAfterFailedDocument(t *testing.T) {
	s := NewSearchService()
	ctx := context.Background()
	if _, err := s.CreateCollection(ctx, domain.Collection{ID: "c", TenantID: "t", Name: "C"}); err != nil {
		t.Fatal(err)
	}
	q := &fakeChangeQueue{docs: []domain.Document{
		{ID: "bad", TenantID: "t", CollectionID: "missing", Version: 1},
		{ID: "good", TenantID: "t", CollectionID: "c", Version: 1, Fields: map[string]any{"body": "text"}},
	}}
	worker := NewIndexWorker(s)
	if err := worker.Run(ctx, q); err != nil {
		t.Fatal(err)
	}
	result, err := s.Query(ctx, domain.Query{TenantID: "t", CollectionID: "c", Text: "text", Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 1 {
		t.Fatalf("worker did not continue to next document after failed one: %#v", result)
	}
}

func TestWorkerDeadLettersOnExhaustedRetries(t *testing.T) {
	s := NewSearchService()
	ctx := context.Background()
	if _, err := s.CreateCollection(ctx, domain.Collection{ID: "c", TenantID: "t", Name: "C"}); err != nil {
		t.Fatal(err)
	}
	q := &fakeChangeQueue{docs: []domain.Document{{ID: "bad", TenantID: "t", CollectionID: "missing", Version: 1}}}
	worker := NewIndexWorker(s)
	if err := worker.Run(ctx, q); err != nil {
		t.Fatal(err)
	}
	if len(q.dead) != 1 {
		t.Fatalf("expected one dead-lettered document, got %d", len(q.dead))
	}
}
