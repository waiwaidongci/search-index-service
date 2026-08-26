package application

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/adapter/cache"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/repository"
	"testing"
)

func seedExportService(t *testing.T) (*Service, *repository.Memory) {
	t.Helper()
	store := repository.NewMemory()
	s := NewService(store, cache.NewMemory())
	ctx := context.Background()
	if err := s.CreateSearchTenant(ctx, domain.SearchTenant{ID: "t", Name: "T"}); err != nil {
		t.Fatal(err)
	}
	if err := s.CreateIndexNamespace(ctx, domain.IndexNamespace{ID: "e", SearchTenantID: "t", Name: "E"}); err != nil {
		t.Fatal(err)
	}
	template := domain.QueryTemplate{
		ID:               "q",
		SearchTenantID:   "t",
		IndexNamespaceID: "e",
		Key:              "flag",
		Type:             domain.TypeString,
		DefaultValue:     "off",
		Rules: []domain.Rule{{
			ID:    "r",
			Value: "original",
			Tags:  map[string]string{"env": "prod"},
		}},
	}
	if err := s.CreateQueryTemplate(ctx, template); err != nil {
		t.Fatal(err)
	}
	return s, store
}

func TestExportSearchTenantDoesNotAliasStore(t *testing.T) {
	s, store := seedExportService(t)
	ctx := context.Background()
	first, err := s.ExportSearchTenant(ctx, "t")
	if err != nil {
		t.Fatal(err)
	}
	stored, err := store.GetQueryTemplate(ctx, "q")
	if err != nil {
		t.Fatal(err)
	}
	stored.Rules[0].Value = "changed"
	stored.Rules[0].Tags["env"] = "staging"
	if err := store.CreateQueryTemplate(ctx, stored); err != nil {
		t.Fatal(err)
	}
	if first.QueryTemplates[0].Rules[0].Value != "original" || first.QueryTemplates[0].Rules[0].Tags["env"] != "prod" {
		t.Fatalf("export aliases store rules: %#v", first.QueryTemplates[0].Rules[0])
	}
}

func TestMutatingTemplateDoesNotChangeExport(t *testing.T) {
	s, _ := seedExportService(t)
	ctx := context.Background()
	first, err := s.ExportSearchTenant(ctx, "t")
	if err != nil {
		t.Fatal(err)
	}
	first.QueryTemplates[0].Rules[0].Value = "changed"
	second, err := s.ExportSearchTenant(ctx, "t")
	if err != nil {
		t.Fatal(err)
	}
	if second.QueryTemplates[0].Rules[0].Value != "original" {
		t.Fatalf("mutating exported result changed the next export: %#v", second.QueryTemplates[0].Rules[0])
	}
}
