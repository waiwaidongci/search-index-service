// Package implementation for tenant-isolated indexing and full-text search.
package repository

import (
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"strings"
)

type QueryTemplateIndex struct {
	bySearchTenant map[string][]string
	byKey          map[string]string
}

func NewQueryTemplateIndex() *QueryTemplateIndex {
	return &QueryTemplateIndex{bySearchTenant: map[string][]string{}, byKey: map[string]string{}}
}
func (i *QueryTemplateIndex) Add(f domain.QueryTemplate) {
	i.bySearchTenant[f.SearchTenantID] = append(i.bySearchTenant[f.SearchTenantID], f.ID)
	i.byKey[strings.ToLower(f.Key)] = f.ID
}
func (i *QueryTemplateIndex) Remove(f domain.QueryTemplate) {
	ids := i.bySearchTenant[f.SearchTenantID]
	out := ids[:0]
	for _, id := range ids {
		if id != f.ID {
			out = append(out, id)
		}
	}
	i.bySearchTenant[f.SearchTenantID] = out
	delete(i.byKey, strings.ToLower(f.Key))
}
func (i *QueryTemplateIndex) FindByKey(key string) string { return i.byKey[strings.ToLower(key)] }
