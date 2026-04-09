package script_test

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

func testInterpreter() []string {
	if runtime.GOOS == "windows" {
		return []string{"pwsh", "-c"}
	}

	return []string{"/bin/bash", "-c"}
}

func testWriteOutputCommand(output string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("[IO.File]::WriteAllText($env:TF_SCRIPT_OUTPUT, '%s')", output)
	}

	return fmt.Sprintf(`printf '%s' > "${TF_SCRIPT_OUTPUT}"`, output)
}

func testWriteErrorCommand(errMsg string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("[IO.File]::WriteAllText($env:TF_SCRIPT_ERROR, '%s'); exit 1", errMsg)
	}

	return fmt.Sprintf(`printf '%s' > "${TF_SCRIPT_ERROR}"; exit 1`, errMsg)
}

func testExitCommand(code int) string {
	return fmt.Sprintf("exit %d", code)
}

// mockLogger records log calls for testing.
type mockLogger struct {
	mu      sync.Mutex
	entries []logEntry
}

type logEntry struct {
	level string
	msg   string
}

func (m *mockLogger) record(level, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, logEntry{level: level, msg: msg})
}

func (m *mockLogger) getEntries() []logEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	copied := make([]logEntry, len(m.entries))
	copy(copied, m.entries)
	return copied
}

func (m *mockLogger) Error(_ context.Context, msg string, _ ...map[string]any) {
	m.record("error", msg)
}

func (m *mockLogger) Warn(_ context.Context, msg string, _ ...map[string]any) {
	m.record("warn", msg)
}

func (m *mockLogger) Info(_ context.Context, msg string, _ ...map[string]any) {
	m.record("info", msg)
}

func (m *mockLogger) Debug(_ context.Context, msg string, _ ...map[string]any) {
	m.record("debug", msg)
}

func (m *mockLogger) Trace(_ context.Context, msg string, _ ...map[string]any) {
	m.record("trace", msg)
}
