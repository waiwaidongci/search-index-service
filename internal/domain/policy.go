// Package implementation for tenant-isolated indexing and full-text search.
package domain

import (
	"fmt"
	"strings"
	"time"
)

type Policy struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	IndexNamespaceID string    `json:"namespace_id"`
	Enabled          bool      `json:"enabled"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (p Policy) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("%w: policy id is required", ErrInvalid)
	}
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("%w: policy name is required", ErrInvalid)
	}
	if strings.TrimSpace(p.IndexNamespaceID) == "" {
		return fmt.Errorf("%w: policy namespace is required", ErrInvalid)
	}
	return nil
}

type Constraint struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

func (c Constraint) Matches(tags map[string]string) bool {
	actual, ok := tags[c.Key]
	if !ok {
		return false
	}
	switch c.Operator {
	case "equals", "":
		return actual == c.Value
	case "not_equals":
		return actual != c.Value
	case "contains":
		return strings.Contains(actual, c.Value)
	case "prefix":
		return strings.HasPrefix(actual, c.Value)
	default:
		return false
	}
}

type Segment struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Constraints []Constraint `json:"constraints"`
	Percentage  int          `json:"percentage"`
}

func (s Segment) Matches(tags map[string]string) bool {
	for _, c := range s.Constraints {
		if !c.Matches(tags) {
			return false
		}
	}
	return true
}
func (s Segment) Validate() error {
	if s.ID == "" || s.Name == "" {
		return fmt.Errorf("%w: segment identity required", ErrInvalid)
	}
	if s.Percentage < 0 || s.Percentage > 100 {
		return fmt.Errorf("%w: percentage must be 0..100", ErrInvalid)
	}
	return nil
}
