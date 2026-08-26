// Package implementation for tenant-isolated indexing and full-text search.
package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource conflict")
	ErrInvalid  = errors.New("invalid resource")
)

type ValueType string

const (
	TypeBool   ValueType = "boolean"
	TypeString ValueType = "string"
	TypeInt    ValueType = "integer"
	TypeJSON   ValueType = "json"
)

type SearchTenant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
type IndexNamespace struct {
	ID             string    `json:"id"`
	SearchTenantID string    `json:"tenant_id"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
}

type Rule struct {
	ID         string            `json:"id"`
	Priority   int               `json:"priority"`
	Tags       map[string]string `json:"tags,omitempty"`
	Percentage *int              `json:"percentage,omitempty"`
	StartAt    *time.Time        `json:"start_at,omitempty"`
	EndAt      *time.Time        `json:"end_at,omitempty"`
	Value      any               `json:"value"`
}

type QueryTemplate struct {
	ID                     string    `json:"id"`
	SearchTenantID         string    `json:"tenant_id"`
	IndexNamespaceID       string    `json:"namespace_id"`
	Key                    string    `json:"key"`
	Description            string    `json:"description,omitempty"`
	Type                   ValueType `json:"type"`
	DefaultValue           any       `json:"default_value"`
	Rules                  []Rule    `json:"rules,omitempty"`
	ActiveTemplateRevision int       `json:"active_revision"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type TemplateRevision struct {
	Number          int    `json:"revision"`
	QueryTemplateID string `json:"template_id"`
	Value           any    `json:"value"`
	Rules           []Rule `json:"rules,omitempty"`
	Status          string `json:"status"`
	CreatedAt       time.Time
	PublishedAt     *time.Time `json:"published_at,omitempty"`
}
type IndexPublication struct {
	ID, QueryTemplateID  string
	TemplateRevision     int
	IndexNamespaceID     string
	Status               string
	CreatedAt, UpdatedAt time.Time
	Reason               string
}

func ValidateValue(t ValueType, value any) error {
	switch t {
	case TypeBool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("%w: expected boolean", ErrInvalid)
		}
	case TypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("%w: expected string", ErrInvalid)
		}
	case TypeInt:
		switch value.(type) {
		case int, int64, float64, float32:
		default:
			return fmt.Errorf("%w: expected integer", ErrInvalid)
		}
	case TypeJSON:
		if value == nil {
			return fmt.Errorf("%w: JSON value cannot be nil", ErrInvalid)
		}
	default:
		return fmt.Errorf("%w: unsupported value type %q", ErrInvalid, t)
	}
	return nil
}

func (f *QueryTemplate) Validate() error {
	if f.ID == "" || f.SearchTenantID == "" || f.IndexNamespaceID == "" || f.Key == "" {
		return fmt.Errorf("%w: template identity is required", ErrInvalid)
	}
	if err := ValidateValue(f.Type, f.DefaultValue); err != nil {
		return err
	}
	return nil
}

func (v *TemplateRevision) Validate(t ValueType) error {
	if v.Number < 1 {
		return fmt.Errorf("%w: revision must be positive", ErrInvalid)
	}
	return ValidateValue(t, v.Value)
}
