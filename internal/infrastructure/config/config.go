// Package implementation for tenant-isolated indexing and full-text search.
package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr        string
	ShutdownSeconds int
	Environment     string
}

func Load() Config {
	c := Config{HTTPAddr: ":8083", ShutdownSeconds: 10, Environment: "local"}
	path := os.Getenv("SEARCH_INDEX_CONFIG")
	if path == "" {
		path = "configs/config.yaml"
	}
	if parsed, err := ParseYAML(path); err == nil {
		c = Merge(c, parsed)
	}
	c = ApplyEnvironment(c)
	if v := os.Getenv("SEARCH_INDEX_SHUTDOWN_SECONDS"); v != "" {
		if n, e := strconv.Atoi(v); e == nil {
			c.ShutdownSeconds = n
		}
	}
	return c
}
func ParseTags(s string) map[string]string {
	out := map[string]string{}
	for _, p := range strings.Split(s, ",") {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = kv[1]
		}
	}
	return out
}
