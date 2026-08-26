package repository

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"sync"
	"testing"
)

func TestConcurrentTemplateRevisionAccessNoRace(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	if err := m.SaveTemplateRevision(ctx, domain.TemplateRevision{QueryTemplateID: "q", Number: 1, Value: "one"}); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _ = m.GetTemplateRevision(ctx, "q", 1)
		}()
	}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _ = m.ListTemplateRevisions(ctx, "q")
		}()
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_ = m.SaveTemplateRevision(ctx, domain.TemplateRevision{QueryTemplateID: "q", Number: i + 2, Value: "next"})
		}(i)
	}
	close(start)
	wg.Wait()
}
