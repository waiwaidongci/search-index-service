// Package implementation for tenant-isolated indexing and full-text search.
package application

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"time"
)

type ChangeQueue interface {
	Consume(context.Context) (<-chan domain.Document, error)
	DeadLetter(context.Context, domain.Document, error) error
}
type IndexWorker struct {
	service *SearchService
	retry   int
	backoff time.Duration
}

func NewIndexWorker(s *SearchService) *IndexWorker {
	return &IndexWorker{service: s, retry: 3, backoff: 100 * time.Millisecond}
}
func (w *IndexWorker) Run(ctx context.Context, q ChangeQueue) error {
	jobs, err := q.Consume(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-jobs:
			if !ok {
				return nil
			}
			var last error
			for i := 0; i < w.retry; i++ {
				_, last = w.service.Index(ctx, d, "")
				if last == nil {
					break
				}
				if i == w.retry-1 {
					break
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(w.backoff):
				}
			}
			if last != nil {
				_ = q.DeadLetter(ctx, d, last)
				continue
			}
		}
	}
}
