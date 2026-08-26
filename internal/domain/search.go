// Package implementation for tenant-isolated indexing and full-text search.
package domain

import (
	"fmt"
	"strings"
	"time"
)

type FieldMapping struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Searchable bool   `json:"searchable"`
	Filterable bool   `json:"filterable"`
}
type Collection struct {
	ID        string         `json:"id"`
	TenantID  string         `json:"tenant_id"`
	Name      string         `json:"name"`
	Version   int            `json:"version"`
	Status    string         `json:"status"`
	Mappings  []FieldMapping `json:"mappings"`
	CreatedAt time.Time      `json:"created_at"`
}
type Document struct {
	ID           string         `json:"id"`
	TenantID     string         `json:"tenant_id"`
	CollectionID string         `json:"collection_id"`
	Version      int64          `json:"version"`
	Fields       map[string]any `json:"fields"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
type Query struct {
	TenantID     string            `json:"tenant_id"`
	CollectionID string            `json:"collection_id"`
	Text         string            `json:"text"`
	Filters      map[string]string `json:"filters"`
	Limit        int               `json:"limit"`
	Cursor       string            `json:"cursor"`
}
type Hit struct {
	ID         string            `json:"id"`
	Score      float64           `json:"score"`
	Fields     map[string]any    `json:"fields"`
	Highlights map[string]string `json:"highlights,omitempty"`
}
type SearchResult struct {
	Hits       []Hit  `json:"hits"`
	Total      int    `json:"total"`
	NextCursor string `json:"next_cursor,omitempty"`
	TookMS     int64  `json:"took_ms"`
}

func (c Collection) Validate() error {
	if c.ID == "" || c.TenantID == "" || c.Name == "" {
		return fmt.Errorf("%w: collection identity required", ErrInvalid)
	}
	for _, m := range c.Mappings {
		switch m.Type {
		case "text", "keyword", "integer", "number", "boolean", "datetime":
		default:
			return fmt.Errorf("%w: field %s has unknown type", ErrInvalid, m.Name)
		}
	}
	return nil
}
func Tokenize(text string) []string {
	parts := strings.Fields(strings.ToLower(text))
	seen := map[string]bool{}
	out := []string{}
	for _, p := range parts {
		p = strings.Trim(p, ".,!?;:\"'()[]{}")
		if p != "" && !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}
