// Package implementation for tenant-isolated indexing and full-text search.
package application

import (
	"context"
	"encoding/json"
	"github.com/ali/go-0821/search-index-service/internal/domain"
)

type Export struct {
	SearchTenant    domain.SearchTenant     `json:"tenant"`
	IndexNamespaces []domain.IndexNamespace `json:"namespaces"`
	QueryTemplates  []domain.QueryTemplate  `json:"templates"`
}

func (s *Service) ExportSearchTenant(ctx context.Context, id string) (Export, error) {
	p, e := s.store.GetSearchTenant(ctx, id)
	if e != nil {
		return Export{}, e
	}
	envs, _ := s.store.ListIndexNamespaces(ctx, id)
	templates := []domain.QueryTemplate{}
	for _, env := range envs {
		items, _ := s.store.ListQueryTemplates(ctx, id, env.ID)
		templates = append(templates, items...)
	}
	return Export{SearchTenant: p, IndexNamespaces: envs, QueryTemplates: templates}, nil
}
func (s *Service) ExportJSON(ctx context.Context, id string) ([]byte, error) {
	v, e := s.ExportSearchTenant(ctx, id)
	if e != nil {
		return nil, e
	}
	return json.MarshalIndent(v, "", "  ")
}
