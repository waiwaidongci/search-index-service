package httpadapter

import (
	"github.com/ali/go-0821/search-index-service/internal/adapter/cache"
	"github.com/ali/go-0821/search-index-service/internal/application"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/logging"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/metrics"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/repository"
	"net/http/httptest"
	"testing"
)

func TestTemplateRouteRejectsInvalidMethodWithErrorChain(t *testing.T) {
	s := New(application.NewService(repository.NewMemory(), cache.NewMemory()), logging.New(), metrics.New())
	req := httptest.NewRequest("GET", "/v1/search/templates/q/switch", nil)
	rec := httptest.NewRecorder()
	s.queryTemplateRoutes(rec, req)
	if rec.Code != 405 {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestBoolQueryPreservesParseError(t *testing.T) {
	req := httptest.NewRequest("GET", "/?enabled=not-bool", nil)
	if _, err := boolQuery(req, "enabled", true); err == nil {
		t.Fatal("expected bool parse error")
	}
}

func TestParseQueryRejectsMalformedTags(t *testing.T) {
	req := httptest.NewRequest("GET", "/?tags=env-prod", nil)
	if _, err := ParseQuery(req); err == nil {
		t.Fatal("expected malformed tags error")
	}
}
