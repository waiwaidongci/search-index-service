// Package implementation for tenant-isolated indexing and full-text search.
// Package httpadapter exposes tenant-isolated indexing and search workflows.
package httpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ali/go-0821/search-index-service/internal/application"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/logging"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/metrics"
	"io"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	templates *application.Service
	search    *application.SearchService
	log       *logging.Logger
	metrics   *metrics.Metrics
	started   time.Time
}

func New(templates *application.Service, l *logging.Logger, m *metrics.Metrics) *Server {
	return &Server{templates: templates, search: application.NewSearchService(), log: l, metrics: m, started: time.Now()}
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/readyz", s.ready)
	mux.HandleFunc("/metrics", s.metrics.Handler)
	mux.HandleFunc("/v1/search/collections", s.collections)
	mux.HandleFunc("/v1/search/documents", s.documents)
	mux.HandleFunc("/v1/search/documents/batch", s.batchDocuments)
	mux.HandleFunc("/v1/search/documents/delete", s.deleteDocument)
	mux.HandleFunc("/v1/search/query", s.searchQuery)
	mux.HandleFunc("/v1/search/rebuild", s.rebuildIndex)
	mux.HandleFunc("/v1/search/tasks/", s.searchTask)
	mux.HandleFunc("/v1/search/tenants", s.searchTenants)
	mux.HandleFunc("/v1/search/namespaces", s.indexNamespaces)
	mux.HandleFunc("/v1/search/templates", s.queryTemplates)
	mux.HandleFunc("/v1/search/templates/", s.queryTemplateRoutes)
	return s.middleware(mux)
}
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.metrics.Inc()
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		w.Header().Set("X-Request-ID", id)
		start := time.Now()
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		_ = ctx
		next.ServeHTTP(w, r)
		s.log.Info("http_request", map[string]any{"method": r.Method, "path": r.URL.Path, "duration_ms": time.Since(start).Milliseconds(), "request_id": id})
	})
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok", "uptime": time.Since(s.started).String()})
}
func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func decode(r *http.Request, v any) error {
	b, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	if len(b) == 0 {
		return fmt.Errorf("%w: empty body", domain.ErrInvalid)
	}
	if err = json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("%w: invalid JSON", domain.ErrInvalid)
	}
	return nil
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, err error) {
	status := 500
	if strings.Contains(err.Error(), domain.ErrNotFound.Error()) {
		status = 404
	} else if strings.Contains(err.Error(), domain.ErrConflict.Error()) {
		status = 409
	} else if strings.Contains(err.Error(), domain.ErrInvalid.Error()) {
		status = 400
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
func Serve(ctx context.Context, addr string, h http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 5 * time.Second}
	errs := make(chan error, 1)
	go func() { errs <- srv.ListenAndServe() }()
	select {
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdown)
	case err := <-errs:
		if err == http.ErrServerClosed {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	}
}
