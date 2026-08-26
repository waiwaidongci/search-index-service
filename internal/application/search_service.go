// Package implementation for tenant-isolated indexing and full-text search.
package application

import (
	"context"
	"fmt"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SearchService struct {
	mu          sync.RWMutex
	collections map[string]domain.Collection
	documents   map[string]domain.Document
	idempotency map[string]string
	tasks       map[string]map[string]any
}

func NewSearchService() *SearchService {
	return &SearchService{collections: map[string]domain.Collection{}, documents: map[string]domain.Document{}, idempotency: map[string]string{}, tasks: map[string]map[string]any{}}
}
func key(t, c, id string) string { return t + "/" + c + "/" + id }
func (s *SearchService) CreateCollection(ctx context.Context, c domain.Collection) (domain.Collection, error) {
	if err := ctx.Err(); err != nil {
		return c, err
	}
	if err := c.Validate(); err != nil {
		return c, err
	}
	c.Version = 1
	c.Status = "active"
	c.CreatedAt = time.Now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(c.TenantID, "collection", c.ID)
	if _, ok := s.collections[k]; ok {
		return c, fmt.Errorf("%w: collection exists", domain.ErrConflict)
	}
	s.collections[k] = c
	return c, nil
}
func (s *SearchService) Index(ctx context.Context, d domain.Document, idem string) (domain.Document, error) {
	if err := ctx.Err(); err != nil {
		return d, err
	}
	if d.ID == "" || d.TenantID == "" || d.CollectionID == "" {
		return d, fmt.Errorf("%w: document identity required", domain.ErrInvalid)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.collections[key(d.TenantID, "collection", d.CollectionID)]; !ok {
		return d, domain.ErrNotFound
	}
	if old, ok := s.documents[key(d.TenantID, d.CollectionID, d.ID)]; ok && d.Version <= old.Version {
		return old, fmt.Errorf("%w: stale document version", domain.ErrConflict)
	}
	if idem != "" {
		scoped := d.TenantID + "|" + idem
		if id, ok := s.idempotency[scoped]; ok {
			return s.documents[id], nil
		}
	}
	d.UpdatedAt = time.Now().UTC()
	s.documents[key(d.TenantID, d.CollectionID, d.ID)] = d
	if idem != "" {
		s.idempotency[d.TenantID+"|"+idem] = key(d.TenantID, d.CollectionID, d.ID)
	}
	return d, nil
}
func (s *SearchService) Query(ctx context.Context, q domain.Query) (domain.SearchResult, error) {
	start := time.Now()
	if err := ctx.Err(); err != nil {
		return domain.SearchResult{}, err
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	offset, _ := strconv.Atoi(q.Cursor)
	tokens := domain.Tokenize(q.Text)
	hits := []domain.Hit{}
	for _, d := range s.documents {
		if d.TenantID != q.TenantID || d.CollectionID != q.CollectionID {
			continue
		}
		matched := len(tokens) == 0
		score := 0.0
		high := map[string]string{}
		filterOK := true
		for k, v := range q.Filters {
			if fmt.Sprint(d.Fields[k]) != v {
				filterOK = false
			}
		}
		if !filterOK {
			continue
		}
		for name, value := range d.Fields {
			text := strings.ToLower(fmt.Sprint(value))
			for _, token := range tokens {
				if strings.Contains(text, token) {
					matched = true
					score++
					high[name] = strings.ReplaceAll(text, token, "<em>"+token+"</em>")
				}
			}
		}
		if matched {
			hits = append(hits, domain.Hit{ID: d.ID, Score: score, Fields: d.Fields, Highlights: high})
		}
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	total := len(hits)
	if offset > total {
		offset = total
	}
	end := offset + q.Limit
	if end > total {
		end = total
	}
	next := ""
	if end < total {
		next = strconv.Itoa(end)
	}
	return domain.SearchResult{Hits: hits[offset:end], Total: total, NextCursor: next, TookMS: time.Since(start).Milliseconds()}, nil
}
func (s *SearchService) Rebuild(ctx context.Context, tenant, collection string) (map[string]any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	count := 0
	for _, d := range s.documents {
		if d.TenantID == tenant && d.CollectionID == collection {
			count++
		}
	}
	s.mu.RUnlock()
	task := fmt.Sprintf("rebuild-%d", time.Now().UnixNano())
	s.mu.Lock()
	s.tasks[task] = map[string]any{"task_id": task, "status": "completed", "documents": count}
	s.mu.Unlock()
	return map[string]any{"task_id": task, "status": "completed", "documents": count}, nil
}

func (s *SearchService) BatchIndex(ctx context.Context, docs []domain.Document, idemPrefix string) ([]domain.Document, error) {
	out := make([]domain.Document, 0, len(docs))
	for i, d := range docs {
		o, e := s.Index(ctx, d, fmt.Sprintf("%s-%d", idemPrefix, i))
		if e != nil {
			return nil, e
		}
		out = append(out, o)
	}
	return out, nil
}
func (s *SearchService) Delete(ctx context.Context, tenant, collection, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k := key(tenant, collection, id)
	if _, ok := s.documents[k]; !ok {
		return domain.ErrNotFound
	}
	delete(s.documents, k)
	return nil
}
func (s *SearchService) Task(id string) (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return v, nil
}
