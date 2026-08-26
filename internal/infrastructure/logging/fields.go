// Package implementation for tenant-isolated indexing and full-text search.
package logging

import (
	"fmt"
	"runtime"
	"time"
)

func Fields(pairs ...any) map[string]any {
	out := map[string]any{}
	for i := 0; i+1 < len(pairs); i += 2 {
		out[fmt.Sprint(pairs[i])] = pairs[i+1]
	}
	return out
}
func WithCaller(fields map[string]any) map[string]any {
	copy := map[string]any{}
	for k, v := range fields {
		copy[k] = v
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		copy["caller"] = fmt.Sprintf("%s:%d", file, line)
	}
	return copy
}
func WithDuration(fields map[string]any, start time.Time) map[string]any {
	copy := map[string]any{}
	for k, v := range fields {
		copy[k] = v
	}
	copy["duration_ms"] = time.Since(start).Milliseconds()
	return copy
}
