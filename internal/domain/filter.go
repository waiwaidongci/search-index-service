// Package implementation for tenant-isolated indexing and full-text search.
package domain

import (
	"sort"
	"strings"
)

type QueryTemplateFilter struct {
	KeyContains      string
	Status           string
	Type             ValueType
	IndexNamespaceID string
	SortBy           string
	Descending       bool
}

func FilterQueryTemplates(templates []QueryTemplate, f QueryTemplateFilter) []QueryTemplate {
	out := make([]QueryTemplate, 0, len(templates))
	for _, item := range templates {
		if f.KeyContains != "" && !strings.Contains(strings.ToLower(item.Key), strings.ToLower(f.KeyContains)) {
			continue
		}
		if f.Status != "" && item.Status != f.Status {
			continue
		}
		if f.Type != "" && item.Type != f.Type {
			continue
		}
		if f.IndexNamespaceID != "" && item.IndexNamespaceID != f.IndexNamespaceID {
			continue
		}
		out = append(out, CopyQueryTemplate(item))
	}
	sort.SliceStable(out, func(i, j int) bool {
		var less bool
		switch f.SortBy {
		case "updated_at":
			less = out[i].UpdatedAt.Before(out[j].UpdatedAt)
		case "created_at":
			less = out[i].CreatedAt.Before(out[j].CreatedAt)
		default:
			less = out[i].Key < out[j].Key
		}
		if f.Descending {
			return !less
		}
		return less
	})
	return out
}
func UniqueKeys(templates []QueryTemplate) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, f := range templates {
		if !seen[f.Key] {
			seen[f.Key] = true
			out = append(out, f.Key)
		}
	}
	sort.Strings(out)
	return out
}
