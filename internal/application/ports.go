// Package implementation for tenant-isolated indexing and full-text search.
package application

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
)

type Store interface {
	CreateSearchTenant(context.Context, domain.SearchTenant) error
	GetSearchTenant(context.Context, string) (domain.SearchTenant, error)
	CreateIndexNamespace(context.Context, domain.IndexNamespace) error
	ListIndexNamespaces(context.Context, string) ([]domain.IndexNamespace, error)
	CreateQueryTemplate(context.Context, domain.QueryTemplate) error
	GetQueryTemplate(context.Context, string) (domain.QueryTemplate, error)
	ListQueryTemplates(context.Context, string, string) ([]domain.QueryTemplate, error)
	SaveTemplateRevision(context.Context, domain.TemplateRevision) error
	GetTemplateRevision(context.Context, string, int) (domain.TemplateRevision, error)
	ListTemplateRevisions(context.Context, string) ([]domain.TemplateRevision, error)
	SaveIndexPublication(context.Context, domain.IndexPublication) error
	ListIndexPublications(context.Context, string) ([]domain.IndexPublication, error)
}

type Cache interface {
	Get(context.Context, string) (*domain.TemplateRevision, bool)
	Set(context.Context, string, domain.TemplateRevision)
	Delete(context.Context, string)
}
