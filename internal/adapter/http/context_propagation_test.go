package httpadapter

import (
	"bytes"
	"context"
	"github.com/ali/go-0821/search-index-service/internal/adapter/cache"
	"github.com/ali/go-0821/search-index-service/internal/application"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/logging"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/metrics"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer() *Server {
	templates := application.NewService(repository.NewMemory(), cache.NewMemory())
	return New(templates, logging.New(), metrics.New())
}

func seedSearchForContextTest(s *Server, t *testing.T) {
	t.Helper()
	ctx := context.Background()
	if _, err := s.search.CreateCollection(ctx, domain.Collection{ID: "c", TenantID: "t", Name: "C"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.search.Index(ctx, domain.Document{ID: "d", TenantID: "t", CollectionID: "c", Version: 1, Fields: map[string]any{"body": "text"}}, ""); err != nil {
		t.Fatal(err)
	}
}

func TestBatchDeleteHonorsCanceledContext(t *testing.T) {
	s := newTestServer()
	seedSearchForContextTest(s, t)
	req := httptest.NewRequest("DELETE", "/v1/search/documents/delete?tenant_id=t&collection_id=c&id=d", nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	s.deleteDocument(rec, req)
	got, err := s.search.Query(context.Background(), domain.Query{TenantID: "t", CollectionID: "c", Text: "text", Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if got.Total != 1 {
		t.Fatalf("canceled delete still removed the document: %#v", got)
	}
}

func TestBatchDocumentsUsesRequestContext(t *testing.T) {
	s := newTestServer()
	ctx := context.Background()
	if _, err := s.search.CreateCollection(ctx, domain.Collection{ID: "c", TenantID: "t", Name: "C"}); err != nil {
		t.Fatal(err)
	}
	body := []byte(`{"documents":[{"id":"d","tenant_id":"t","collection_id":"c","version":1,"fields":{"body":"text"}}]}`)
	req := httptest.NewRequest("POST", "/v1/search/documents/batch", bytes.NewReader(body))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	s.batchDocuments(rec, req)
	got, err := s.search.Query(context.Background(), domain.Query{TenantID: "t", CollectionID: "c", Text: "text", Limit: 20})
	if err != nil {
		t.Fatal(err)
	}
	if got.Total != 0 {
		t.Fatalf("canceled batch write still inserted documents: %#v", got)
	}
}

func TestTimeoutContextReachesSearchHandler(t *testing.T) {
	s := newTestServer()
	var observed context.Context
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		observed = r.Context()
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()
	s.middleware(next).ServeHTTP(rec, req)
	if _, hasDeadline := observed.Deadline(); !hasDeadline {
		t.Fatalf("middleware did not propagate timeout context")
	}
}

func TestRebuildHonorsCanceledContext(t *testing.T) {
	s := newTestServer()
	seedSearchForContextTest(s, t)
	req := httptest.NewRequest("POST", "/v1/search/rebuild", bytes.NewBufferString(`{"tenant_id":"t","collection_id":"c"}`))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	s.rebuildIndex(rec, req)
	if rec.Code == 202 {
		t.Fatalf("canceled rebuild still created a task: %#v", rec.Body.String())
	}
}

func TestSearchQueryHonorsCanceledContext(t *testing.T) {
	s := newTestServer()
	seedSearchForContextTest(s, t)
	req := httptest.NewRequest("POST", "/v1/search/query", bytes.NewBufferString(`{"tenant_id":"t","collection_id":"c","text":"text"}`))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	s.searchQuery(rec, req)
	if rec.Code == 200 {
		t.Fatalf("canceled query still succeeded: %#v", rec.Body.String())
	}
}
