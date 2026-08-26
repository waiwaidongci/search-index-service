// Package implementation for tenant-isolated indexing and full-text search.
package logging

import (
	"encoding/json"
	"log"
	"time"
)

type Logger struct{ base *log.Logger }

func New() *Logger { return &Logger{base: log.Default()} }
func (l *Logger) Event(level, msg string, fields map[string]any) {
	if fields == nil {
		fields = map[string]any{}
	}
	fields["level"] = level
	fields["message"] = msg
	fields["ts"] = time.Now().UTC().Format(time.RFC3339Nano)
	b, _ := json.Marshal(fields)
	l.base.Print(string(b))
}
func (l *Logger) Info(msg string, fields map[string]any)  { l.Event("info", msg, fields) }
func (l *Logger) Error(msg string, fields map[string]any) { l.Event("error", msg, fields) }
