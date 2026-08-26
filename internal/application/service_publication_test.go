package application

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ali/go-0821/search-index-service/internal/adapter/cache"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/repository"
	"testing"
)

type failingStore struct {
	*repository.Memory
	failNamespace bool
	failTemplate  bool
}

func (f *failingStore) CreateIndexNamespace(ctx context.Context, e domain.IndexNamespace) error {
	if f.failNamespace {
		return domain.ErrInvalid
	}
	return f.Memory.CreateIndexNamespace(ctx, e)
}

func (f *failingStore) CreateQueryTemplate(ctx context.Context, template domain.QueryTemplate) error {
	if f.failTemplate {
		return domain.ErrInvalid
	}
	return f.Memory.CreateQueryTemplate(ctx, template)
}

func seedPublicationService(t *testing.T) *Service {
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
	if err := s.CreateQueryTemplate(ctx, domain.QueryTemplate{ID: "q", SearchTenantID: "t", IndexNamespaceID: "e", Key: "flag", Type: domain.TypeString, DefaultValue: "off"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateTemplateRevision(ctx, "q", domain.TemplateRevision{Value: "on"}); err != nil {
		t.Fatal(err)
	}
	return s
}

func TestPublicationErrorPreservesNotFoundChain(t *testing.T) {
	s := seedPublicationService(t)
	_, err := s.IndexPublication(context.Background(), "q", 999, "e", "version switch")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound in error chain, got %v", err)
	}
}

func TestRollbackRejectsInvalidRevision(t *testing.T) {
	s := seedPublicationService(t)
	_, err := s.Rollback(context.Background(), "q", 999, "e")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound from rollback, got %v", err)
	}
}

func TestImportPartialFailureIsNotSilentlyIgnored(t *testing.T) {
	store := &failingStore{Memory: repository.NewMemory(), failNamespace: true}
	s := NewService(store, cache.NewMemory())
	payload, err := json.Marshal(Export{
		SearchTenant:    domain.SearchTenant{ID: "t", Name: "T"},
		IndexNamespaces: []domain.IndexNamespace{{ID: "e", SearchTenantID: "t", Name: "E"}},
		QueryTemplates:  []domain.QueryTemplate{{ID: "q", SearchTenantID: "t", IndexNamespaceID: "e", Key: "flag", Type: domain.TypeString, DefaultValue: "off"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.ImportSearchTenant(context.Background(), payload); err == nil {
		t.Fatal("expected import partial failure to return an error")
	}
}
