// Package implementation for tenant-isolated indexing and full-text search.
package application

import (
	"context"
	"fmt"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"sort"
	"time"
)

type Service struct {
	store Store
	cache Cache
	now   func() time.Time
}

func NewService(store Store, cache Cache) *Service {
	return &Service{store: store, cache: cache, now: time.Now}
}

func (s *Service) CreateSearchTenant(ctx context.Context, p domain.SearchTenant) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = s.now()
	}
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("%w: tenant id/name required", domain.ErrInvalid)
	}
	return s.store.CreateSearchTenant(ctx, p)
}
func (s *Service) CreateIndexNamespace(ctx context.Context, e domain.IndexNamespace) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = s.now()
	}
	if e.ID == "" || e.SearchTenantID == "" || e.Name == "" {
		return fmt.Errorf("%w: namespace identity required", domain.ErrInvalid)
	}
	return s.store.CreateIndexNamespace(ctx, e)
}
func (s *Service) CreateQueryTemplate(ctx context.Context, f domain.QueryTemplate) error {
	if f.CreatedAt.IsZero() {
		f.CreatedAt = s.now()
	}
	f.UpdatedAt = f.CreatedAt
	f.Status = "draft"
	if err := f.Validate(); err != nil {
		return err
	}
	if err := domain.ValidateRules(f.Type, f.Rules); err != nil {
		return err
	}
	return s.store.CreateQueryTemplate(ctx, f)
}
func (s *Service) GetQueryTemplate(ctx context.Context, id string) (domain.QueryTemplate, error) {
	return s.store.GetQueryTemplate(ctx, id)
}
func (s *Service) ListQueryTemplates(ctx context.Context, p, e string) ([]domain.QueryTemplate, error) {
	return s.store.ListQueryTemplates(ctx, p, e)
}
func (s *Service) ListIndexNamespaces(ctx context.Context, p string) ([]domain.IndexNamespace, error) {
	return s.store.ListIndexNamespaces(ctx, p)
}

func (s *Service) CreateTemplateRevision(ctx context.Context, templateID string, v domain.TemplateRevision) (domain.TemplateRevision, error) {
	f, err := s.store.GetQueryTemplate(ctx, templateID)
	if err != nil {
		return v, err
	}
	revisions, _ := s.store.ListTemplateRevisions(ctx, templateID)
	v.QueryTemplateID = templateID
	v.Number = len(revisions) + 1
	v.Status = "draft"
	v.CreatedAt = s.now()
	if err := v.Validate(f.Type); err != nil {
		return v, err
	}
	if err := domain.ValidateRules(f.Type, v.Rules); err != nil {
		return v, err
	}
	if err := s.store.SaveTemplateRevision(ctx, v); err != nil {
		return v, err
	}
	return v, nil
}
func (s *Service) ListTemplateRevisions(ctx context.Context, id string) ([]domain.TemplateRevision, error) {
	return s.store.ListTemplateRevisions(ctx, id)
}

func (s *Service) IndexPublication(ctx context.Context, templateID string, revision int, envID, reason string) (domain.IndexPublication, error) {
	f, err := s.store.GetQueryTemplate(ctx, templateID)
	if err != nil {
		return domain.IndexPublication{}, err
	}
	v, err := s.store.GetTemplateRevision(ctx, templateID, revision)
	if err != nil {
		return domain.IndexPublication{}, fmt.Errorf("publish revision: %v", err)
	}
	if v.Status == "revoked" {
		return domain.IndexPublication{}, fmt.Errorf("%w: revision revoked", domain.ErrConflict)
	}
	now := s.now()
	rel := domain.IndexPublication{ID: fmt.Sprintf("rel-%d", now.UnixNano()), QueryTemplateID: templateID, TemplateRevision: revision, IndexNamespaceID: envID, Status: "published", CreatedAt: now, UpdatedAt: now, Reason: reason}
	v.Status = "published"
	v.PublishedAt = &now
	if err = s.store.SaveTemplateRevision(ctx, v); err != nil {
		return rel, err
	}
	f.ActiveTemplateRevision = revision
	f.Status = "published"
	f.UpdatedAt = now
	if err = s.store.CreateQueryTemplate(ctx, f); err != nil {
		return rel, err
	}
	s.cache.Delete(ctx, templateID)
	if err = s.store.SaveIndexPublication(ctx, rel); err != nil {
		return rel, err
	}
	return rel, nil
}
func (s *Service) Rollback(ctx context.Context, templateID string, revision int, envID string) (domain.IndexPublication, error) {
	rel, _ := s.IndexPublication(ctx, templateID, revision, envID, "rollback")
	return rel, nil
}
func (s *Service) ListIndexPublications(ctx context.Context, id string) ([]domain.IndexPublication, error) {
	return s.store.ListIndexPublications(ctx, id)
}

type ValueResult struct {
	QueryTemplateID, Key string
	Value                any    `json:"value"`
	TemplateRevision     int    `json:"revision"`
	ETag                 string `json:"etag"`
	Source               string `json:"source"`
}

func (s *Service) Evaluate(ctx context.Context, templateID string, ec domain.EvaluationContext) (ValueResult, error) {
	f, err := s.store.GetQueryTemplate(ctx, templateID)
	if err != nil {
		return ValueResult{}, err
	}
	var v *domain.TemplateRevision
	if f.ActiveTemplateRevision > 0 {
		if cv, ok := s.cache.Get(ctx, templateID); ok {
			v = cv
		} else if cv, e := s.store.GetTemplateRevision(ctx, templateID, f.ActiveTemplateRevision); e == nil {
			v = &cv
			s.cache.Set(ctx, templateID, cv)
		}
	}
	value, no, err := domain.Evaluate(f, v, ec)
	if err != nil {
		return ValueResult{}, err
	}
	return ValueResult{QueryTemplateID: f.ID, Key: f.Key, Value: value, TemplateRevision: no, ETag: fmt.Sprintf("%s-v%d", f.ID, no), Source: "default"}, nil
}
func (s *Service) BatchEvaluate(ctx context.Context, tenant, env string, ec domain.EvaluationContext) ([]ValueResult, error) {
	fs, err := s.store.ListQueryTemplates(ctx, tenant, env)
	if err != nil {
		return nil, err
	}
	out := make([]ValueResult, 0, len(fs))
	for _, f := range fs {
		r, e := s.Evaluate(ctx, f.ID, ec)
		if e != nil {
			return nil, e
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}
