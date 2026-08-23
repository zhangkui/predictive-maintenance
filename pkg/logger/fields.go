package logger

import (
	"fmt"
	"runtime"
	"time"
)

func Duration(start time.Time) map[string]any {
	return map[string]any{"durationMs": time.Since(start).Milliseconds()}
}
func Error(err error) map[string]any {
	if err == nil {
		return nil
	}
	return map[string]any{"error": err.Error()}
}
func Caller(skip int) map[string]any {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return nil
	}
	return map[string]any{"file": file, "line": line}
}
func Request(method, path string) map[string]any {
	return map[string]any{"method": method, "path": path}
}
func String(key string, value any) map[string]any { return map[string]any{key: fmt.Sprint(value)} }
