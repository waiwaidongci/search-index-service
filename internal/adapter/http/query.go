// Package implementation for tenant-isolated indexing and full-text search.
package httpadapter

import (
	"fmt"
	"github.com/ali/go-0821/search-index-service/internal/domain"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/config"
	"net/http"
	"strconv"
	"strings"
)

type Query struct {
	ProjectID     string
	EnvironmentID string
	Page          int
	Size          int
	Tags          map[string]string
}

func ParseQuery(r *http.Request) (Query, error) {
	q := Query{Tags: map[string]string{}}
	v := r.URL.Query()
	page, err := strconv.Atoi(v.Get("page"))
	if err != nil && v.Get("page") != "" {
		return q, fmt.Errorf("%w: invalid page parameter", domain.ErrInvalid)
	}
	size, err := strconv.Atoi(v.Get("size"))
	if err != nil && v.Get("size") != "" {
		return q, fmt.Errorf("%w: invalid size parameter", domain.ErrInvalid)
	}
	tags, err := config.ParseTags(v.Get("tags"))
	if err != nil {
		return q, fmt.Errorf("%w: invalid tags parameter", domain.ErrInvalid)
	}
	return Query{ProjectID: strings.TrimSpace(v.Get("project_id")), EnvironmentID: strings.TrimSpace(v.Get("environment_id")), Page: page, Size: size, Tags: tags}, nil
}
func boolQuery(r *http.Request, key string, def bool) (bool, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def, nil
	}
	b, e := strconv.ParseBool(v)
	if e != nil {
		return false, fmt.Errorf("%w: invalid %s parameter", domain.ErrInvalid, key)
	}
	return b, nil
}
