// Package implementation for tenant-isolated indexing and full-text search.
package httpadapter

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"net/http"
	"strings"
)

func (s *Server) collections(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var c domain.Collection
	if err := decode(r, &c); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.search.CreateCollection(r.Context(), c)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, o)
}
func (s *Server) documents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var d domain.Document
	if err := decode(r, &d); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.search.Index(r.Context(), d, r.Header.Get("Idempotency-Key"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, o)
}
func (s *Server) batchDocuments(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Documents []domain.Document `json:"documents"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.search.BatchIndex(context.Background(), in.Documents, r.Header.Get("Idempotency-Key"))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, map[string]any{"documents": o, "count": len(o)})
}
func (s *Server) deleteDocument(w http.ResponseWriter, r *http.Request) {
	if err := s.search.Delete(context.Background(), r.URL.Query().Get("tenant_id"), r.URL.Query().Get("collection_id"), r.URL.Query().Get("id")); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, map[string]string{"status": "deleted"})
}
func (s *Server) searchTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/search/tasks/")
	o, err := s.search.Task(id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, o)
}
func (s *Server) searchQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var q domain.Query
	if err := decode(r, &q); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.search.Query(context.Background(), q)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 200, o)
}
func (s *Server) rebuildIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct {
		TenantID     string `json:"tenant_id"`
		CollectionID string `json:"collection_id"`
	}
	if err := decode(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	o, err := s.search.Rebuild(context.Background(), in.TenantID, in.CollectionID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 202, o)
}
