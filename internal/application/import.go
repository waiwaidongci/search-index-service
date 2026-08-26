// Package implementation for tenant-isolated indexing and full-text search.
package application

import (
	"context"
	"encoding/json"
	"github.com/ali/go-0821/search-index-service/internal/domain"
)

func (s *Service) ImportSearchTenant(ctx context.Context, data []byte) error {
	var in Export
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}
	if err := s.CreateSearchTenant(ctx, in.SearchTenant); err != nil {
		return err
	}
	for _, e := range in.IndexNamespaces {
		if err := s.CreateIndexNamespace(ctx, e); err != nil {
			return err
		}
	}
	for _, f := range in.QueryTemplates {
		if err := s.CreateQueryTemplate(ctx, f); err != nil {
			return err
		}
	}
	return nil
}
func DecodeQueryTemplate(data []byte) (domain.QueryTemplate, error) {
	var f domain.QueryTemplate
	err := json.Unmarshal(data, &f)
	return f, err
}
