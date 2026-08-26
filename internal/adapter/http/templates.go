package httpadapter

import (
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"net/http"
	"strings"
)

func (s *Server) searchTenants(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v domain.SearchTenant
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.templates.CreateSearchTenant(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
func (s *Server) queryTemplateRoutes(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 5 {
		w.WriteHeader(404)
		return
	}
	id, action := parts[3], parts[4]
	switch action {
	case "revisions":
		if r.Method == "GET" {
			o, e := s.templates.ListTemplateRevisions(r.Context(), id)
			if e != nil {
				writeErr(w, e)
				return
			}
			writeJSON(w, 200, o)
			return
		}
		var in domain.TemplateRevision
		if e := decode(r, &in); e != nil {
			writeErr(w, e)
			return
		}
		o, e := s.templates.CreateTemplateRevision(r.Context(), id, in)
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 201, o)
	case "switch":
		var in struct {
			Revision    int    `json:"revision"`
			NamespaceID string `json:"namespace_id"`
		}
		if e := decode(r, &in); e != nil {
			writeErr(w, e)
			return
		}
		o, e := s.templates.IndexPublication(r.Context(), id, in.Revision, in.NamespaceID, "version switch")
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, o)
	default:
		w.WriteHeader(404)
	}
}
func (s *Server) indexNamespaces(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		o, e := s.templates.ListIndexNamespaces(r.Context(), r.URL.Query().Get("tenant_id"))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, o)
		return
	}
	var v domain.IndexNamespace
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.templates.CreateIndexNamespace(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
func (s *Server) queryTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		o, e := s.templates.ListQueryTemplates(r.Context(), r.URL.Query().Get("tenant_id"), r.URL.Query().Get("namespace_id"))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 200, o)
		return
	}
	var v domain.QueryTemplate
	if err := decode(r, &v); err != nil {
		writeErr(w, err)
		return
	}
	if err := s.templates.CreateQueryTemplate(r.Context(), v); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, 201, v)
}
