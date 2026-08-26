// Package implementation for tenant-isolated indexing and full-text search.
package repository

import (
	"fmt"
	"github.com/ali/go-0821/search-index-service/internal/domain"
)

func validateSearchTenant(p domain.SearchTenant) error {
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("%w: invalid tenant", domain.ErrInvalid)
	}
	return nil
}
func validateIndexNamespace(e domain.IndexNamespace) error {
	if e.ID == "" || e.SearchTenantID == "" || e.Name == "" {
		return fmt.Errorf("%w: invalid namespace", domain.ErrInvalid)
	}
	return nil
}
func validateQueryTemplate(f domain.QueryTemplate) error { return f.Validate() }
func validateTemplateRevision(v domain.TemplateRevision, t domain.ValueType) error {
	return v.Validate(t)
}
