// Package implementation for tenant-isolated indexing and full-text search.
package domain

import (
	"encoding/json"
	"reflect"
	"sort"
)

type Change struct {
	Path   string `json:"path"`
	Before any    `json:"before"`
	After  any    `json:"after"`
	Kind   string `json:"kind"`
}

func CompareTemplateRevisions(a, b TemplateRevision) []Change {
	out := []Change{}
	if !reflect.DeepEqual(a.Value, b.Value) {
		out = append(out, Change{Path: "value", Before: a.Value, After: b.Value, Kind: "changed"})
	}
	la, lb := len(a.Rules), len(b.Rules)
	if la != lb {
		out = append(out, Change{Path: "rules", Before: la, After: lb, Kind: "changed"})
	}
	n := la
	if lb > n {
		n = lb
	}
	for i := 0; i < n; i++ {
		switch {
		case i >= la:
			out = append(out, Change{Path: "rules[" + itoa(i) + "]", Before: nil, After: b.Rules[i], Kind: "added"})
		case i >= lb:
			out = append(out, Change{Path: "rules[" + itoa(i) + "]", Before: a.Rules[i], After: nil, Kind: "removed"})
		default:
			if !reflect.DeepEqual(a.Rules[i], b.Rules[i]) {
				out = append(out, Change{Path: "rules[" + itoa(i) + "]", Before: a.Rules[i], After: b.Rules[i], Kind: "changed"})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
func itoa(n int) string { b, _ := json.Marshal(n); return string(b) }
