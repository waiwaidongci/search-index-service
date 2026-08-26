// Package implementation for tenant-isolated indexing and full-text search.
package domain

import (
	"encoding/json"
	"fmt"
)

func EncodeValue(v any) ([]byte, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, fmt.Errorf("encode value: %w", e)
	}
	return b, nil
}
func DecodeValue(data []byte, t ValueType) (any, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("decode value: %w", err)
	}
	if err := ValidateValue(t, v); err != nil {
		return nil, err
	}
	return v, nil
}
func CloneRules(in []Rule) []Rule {
	out := make([]Rule, len(in))
	for i, r := range in {
		out[i] = r
		if r.Tags != nil {
			out[i].Tags = map[string]string{}
			for k, v := range r.Tags {
				out[i].Tags[k] = v
			}
		}
	}
	return out
}
func CopyQueryTemplate(in QueryTemplate) QueryTemplate { in.Rules = CloneRules(in.Rules); return in }
