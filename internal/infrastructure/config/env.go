// Package implementation for tenant-isolated indexing and full-text search.
package config

import (
	"os"
	"strings"
)

func ApplyEnvironment(c Config) Config {
	for _, item := range []struct {
		key   string
		apply func(string)
	}{{"SEARCH_INDEX_HTTP_ADDR", func(v string) { c.HTTPAddr = v }}, {"SEARCH_INDEX_ENVIRONMENT", func(v string) { c.Environment = v }}} {
		if v := strings.TrimSpace(os.Getenv(item.key)); v != "" {
			item.apply(v)
		}
	}
	return c
}
func (c Config) Public() map[string]any {
	return map[string]any{"http_addr": c.HTTPAddr, "environment": c.Environment, "shutdown_seconds": c.ShutdownSeconds}
}
