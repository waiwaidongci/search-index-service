// Package implementation for tenant-isolated indexing and full-text search.
package repository

import "github.com/ali/go-0821/search-index-service/internal/domain"

func cloneSearchTenant(p domain.SearchTenant) domain.SearchTenant       { return p }
func cloneIndexNamespace(e domain.IndexNamespace) domain.IndexNamespace { return e }
func cloneTemplateRevision(v domain.TemplateRevision) domain.TemplateRevision {
	v.Rules = domain.CloneRules(v.Rules)
	return v
}
func cloneIndexPublication(r domain.IndexPublication) domain.IndexPublication { return r }
func cloneQueryTemplate(f domain.QueryTemplate) domain.QueryTemplate {
	return domain.CopyQueryTemplate(f)
}
