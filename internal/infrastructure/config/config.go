// Package implementation for tenant-isolated indexing and full-text search.
package config

import (
	"fmt"
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
func ParseTags(s string) (map[string]string, error) {
	out := map[string]string{}
	if s == "" {
		return out, nil
	}
	for _, p := range strings.Split(s, ",") {
		entry := strings.TrimSpace(p)
		if entry == "" {
			return out, fmt.Errorf("empty tag entry in %q", s)
		}
		kv := strings.SplitN(entry, "=", 2)
		if len(kv) != 2 || kv[0] == "" {
			return out, fmt.Errorf("malformed tag entry %q", entry)
		}
		out[kv[0]] = kv[1]
	}
	return out, nil
}
