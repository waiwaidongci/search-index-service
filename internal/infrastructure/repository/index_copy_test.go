package repository

import (
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"testing"
)

func TestIndexRemoveDoesNotMutateExistingSlice(t *testing.T) {
	idx := NewQueryTemplateIndex()
	idx.Add(domain.QueryTemplate{ID: "a", SearchTenantID: "t", Key: "a"})
	idx.Add(domain.QueryTemplate{ID: "b", SearchTenantID: "t", Key: "b"})
	before := idx.bySearchTenant["t"]
	idx.Remove(domain.QueryTemplate{ID: "a", SearchTenantID: "t", Key: "a"})
	if len(before) != 2 || before[0] != "a" || before[1] != "b" {
		t.Fatalf("old holder was mutated: %#v", before)
	}
}

func TestOldIndexHolderStaysStableAfterRemove(t *testing.T) {
	idx := NewQueryTemplateIndex()
	idx.Add(domain.QueryTemplate{ID: "a", SearchTenantID: "t", Key: "a"})
	idx.Add(domain.QueryTemplate{ID: "b", SearchTenantID: "t", Key: "b"})
	idx.Add(domain.QueryTemplate{ID: "c", SearchTenantID: "t", Key: "c"})
	before := idx.bySearchTenant["t"]
	idx.Remove(domain.QueryTemplate{ID: "b", SearchTenantID: "t", Key: "b"})
	idx.Remove(domain.QueryTemplate{ID: "a", SearchTenantID: "t", Key: "a"})
	if len(before) != 3 || before[0] != "a" || before[1] != "b" || before[2] != "c" {
		t.Fatalf("old holder changed after removals: %#v", before)
	}
}
