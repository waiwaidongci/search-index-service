// Package implementation for tenant-isolated indexing and full-text search.
package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type YAML struct{ Values map[string]string }

func ParseYAML(path string) (YAML, error) {
	f, e := os.Open(path)
	if e != nil {
		return YAML{}, e
	}
	defer f.Close()
	y := YAML{Values: map[string]string{}}
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}
		parts := strings.SplitN(raw, ":", 2)
		if len(parts) != 2 {
			return y, fmt.Errorf("line %d: invalid YAML", line)
		}
		y.Values[strings.TrimSpace(parts[0])] = strings.Trim(strings.TrimSpace(parts[1]), "\"'")
	}
	return y, scanner.Err()
}
func (y YAML) String(key, def string) string {
	if v, ok := y.Values[key]; ok && v != "" {
		return v
	}
	return def
}
func (y YAML) Int(key string, def int) int {
	v := y.String(key, "")
	if n, e := strconv.Atoi(v); e == nil {
		return n
	}
	return def
}
func Merge(base Config, y YAML) Config {
	base.HTTPAddr = y.String("http_addr", base.HTTPAddr)
	base.Environment = y.String("environment", base.Environment)
	base.ShutdownSeconds = y.Int("shutdown_seconds", base.ShutdownSeconds)
	return base
}
