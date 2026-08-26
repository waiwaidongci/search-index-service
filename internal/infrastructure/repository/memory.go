// Package implementation for tenant-isolated indexing and full-text search.
package repository

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"sync"
)

type Memory struct {
	mu           sync.RWMutex
	tenants      map[string]domain.SearchTenant
	envs         map[string]domain.IndexNamespace
	templates    map[string]domain.QueryTemplate
	revisions    map[string][]domain.TemplateRevision
	publications map[string][]domain.IndexPublication
}

func NewMemory() *Memory {
	return &Memory{tenants: map[string]domain.SearchTenant{}, envs: map[string]domain.IndexNamespace{}, templates: map[string]domain.QueryTemplate{}, revisions: map[string][]domain.TemplateRevision{}, publications: map[string][]domain.IndexPublication{}}
}
func (m *Memory) CreateSearchTenant(_ context.Context, p domain.SearchTenant) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tenants[p.ID]; ok {
		return domain.ErrConflict
	}
	m.tenants[p.ID] = p
	return nil
}
func (m *Memory) GetSearchTenant(_ context.Context, id string) (domain.SearchTenant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.tenants[id]
	if !ok {
		return domain.SearchTenant{}, domain.ErrNotFound
	}
	return p, nil
}
func (m *Memory) CreateIndexNamespace(_ context.Context, e domain.IndexNamespace) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.envs[e.ID]; ok {
		return domain.ErrConflict
	}
	e2 := e
	m.envs[e.ID] = e2
	return nil
}
func (m *Memory) ListIndexNamespaces(_ context.Context, p string) ([]domain.IndexNamespace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.IndexNamespace{}
	for _, e := range m.envs {
		if e.SearchTenantID == p {
			out = append(out, e)
		}
	}
	return out, nil
}
func (m *Memory) CreateQueryTemplate(_ context.Context, f domain.QueryTemplate) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if old, ok := m.templates[f.ID]; ok {
		f.CreatedAt = old.CreatedAt
	}
	m.templates[f.ID] = f
	return nil
}
func (m *Memory) GetQueryTemplate(_ context.Context, id string) (domain.QueryTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.templates[id]
	if !ok {
		return domain.QueryTemplate{}, domain.ErrNotFound
	}
	return f, nil
}
func (m *Memory) ListQueryTemplates(_ context.Context, p, e string) ([]domain.QueryTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.QueryTemplate{}
	for _, f := range m.templates {
		if f.SearchTenantID == p && (e == "" || f.IndexNamespaceID == e) {
			out = append(out, f)
		}
	}
	return out, nil
}
func (m *Memory) SaveTemplateRevision(_ context.Context, v domain.TemplateRevision) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	arr := m.revisions[v.QueryTemplateID]
	for i, x := range arr {
		if x.Number == v.Number {
			arr[i] = v
			m.revisions[v.QueryTemplateID] = arr
			return nil
		}
	}
	m.revisions[v.QueryTemplateID] = append(arr, v)
	return nil
}
func (m *Memory) GetTemplateRevision(_ context.Context, id string, n int) (domain.TemplateRevision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.revisions[id] {
		if v.Number == n {
			return v, nil
		}
	}
	return domain.TemplateRevision{}, domain.ErrNotFound
}
func (m *Memory) ListTemplateRevisions(_ context.Context, id string) ([]domain.TemplateRevision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.TemplateRevision(nil), m.revisions[id]...), nil
}
func (m *Memory) SaveIndexPublication(_ context.Context, r domain.IndexPublication) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publications[r.QueryTemplateID] = append(m.publications[r.QueryTemplateID], r)
	return nil
}
func (m *Memory) ListIndexPublications(_ context.Context, id string) ([]domain.IndexPublication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.IndexPublication(nil), m.publications[id]...), nil
}
