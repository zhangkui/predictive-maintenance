package logger

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	mu     sync.Mutex
	std    *log.Logger
	fields map[string]any
}

func New() *Logger { return &Logger{std: log.New(os.Stdout, "", 0), fields: map[string]any{}} }
func (l *Logger) With(key string, value any) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	next := New()
	for k, v := range l.fields {
		next.fields[k] = v
	}
	next.fields[key] = value
	return next
}
func (l *Logger) write(level, message string, fields map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	payload := map[string]any{"time": time.Now().UTC().Format(time.RFC3339Nano), "level": level, "message": message}
	for k, v := range l.fields {
		payload[k] = v
	}
	for k, v := range fields {
		payload[k] = v
	}
	b, _ := json.Marshal(payload)
	l.std.Println(string(b))
}
func (l *Logger) Info(message string, fields ...map[string]any) {
	l.write("info", message, merge(fields))
}
func (l *Logger) Warn(message string, fields ...map[string]any) {
	l.write("warn", message, merge(fields))
}
func (l *Logger) Error(message string, fields ...map[string]any) {
	l.write("error", message, merge(fields))
}
func (l *Logger) Debug(message string, fields ...map[string]any) {
	l.write("debug", message, merge(fields))
}
func merge(values []map[string]any) map[string]any {
	out := map[string]any{}
	for _, v := range values {
		for k, x := range v {
			out[k] = x
		}
	}
	return out
}
