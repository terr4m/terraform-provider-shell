package shell

import (
	"context"
	"reflect"
	"runtime"
	"testing"
)

type testLogger struct {
	errors   []string
	warnings []string
	infos    []string
	debugs   []string
	traces   []string
}

func (l *testLogger) Error(ctx context.Context, msg string, additionalFields ...map[string]any) {
	l.errors = append(l.errors, msg)
}

func (l *testLogger) Warn(ctx context.Context, msg string, additionalFields ...map[string]any) {
	l.warnings = append(l.warnings, msg)
}

func (l *testLogger) Info(ctx context.Context, msg string, additionalFields ...map[string]any) {
	l.infos = append(l.infos, msg)
}

func (l *testLogger) Debug(ctx context.Context, msg string, additionalFields ...map[string]any) {
	l.debugs = append(l.debugs, msg)
}

func (l *testLogger) Trace(ctx context.Context, msg string, additionalFields ...map[string]any) {
	l.traces = append(l.traces, msg)
}

func TestRunCommand(t *testing.T) {
	t.Parallel()

	interpreter := []string{"/bin/bash", "-c"}
	if runtime.GOOS == "windows" {
		interpreter = []string{"pwsh", "-c"}
	}

	for _, d := range []struct {
		testName     string
		winOK        bool
		interpreter  []string
		env          map[string]string
		dir          string
		command      string
		logger       *testLogger
		hasErr       bool
		loggerResult *testLogger
	}{
		{
			testName:     "invalid_interpreter",
			winOK:        true,
			interpreter:  []string{"xxxxxx"},
			env:          nil,
			dir:          "",
			command:      `echo "hello"`,
			logger:       nil,
			hasErr:       true,
			loggerResult: nil,
		},
		{
			testName:     "invalid_command",
			winOK:        true,
			interpreter:  interpreter,
			env:          nil,
			dir:          "",
			command:      `xxxxxx`,
			logger:       nil,
			hasErr:       true,
			loggerResult: nil,
		},
		{
			testName:     "echo",
			winOK:        true,
			interpreter:  interpreter,
			env:          nil,
			dir:          "",
			command:      `echo "hello"`,
			logger:       nil,
			hasErr:       false,
			loggerResult: nil,
		},
		{
			testName:     "echo_env",
			winOK:        true,
			interpreter:  interpreter,
			env:          map[string]string{"TEST": "hello"},
			dir:          "",
			command:      `echo "${TEST}"`,
			logger:       nil,
			hasErr:       false,
			loggerResult: nil,
		},
		{
			testName:     "echo_change_dir",
			winOK:        false,
			interpreter:  interpreter,
			env:          nil,
			dir:          "/tmp",
			command:      `echo "hello"`,
			logger:       nil,
			hasErr:       false,
			loggerResult: nil,
		},
		{
			testName:     "check_log_empty",
			winOK:        true,
			interpreter:  interpreter,
			env:          nil,
			dir:          "/tmp",
			command:      `echo "Test..."`,
			logger:       &testLogger{},
			hasErr:       false,
			loggerResult: &testLogger{},
		},
		{
			testName:     "check_log_error",
			winOK:        true,
			interpreter:  interpreter,
			env:          nil,
			dir:          "/tmp",
			command:      `echo "[ERROR] Test error"`,
			logger:       &testLogger{},
			hasErr:       false,
			loggerResult: &testLogger{errors: []string{"Test error"}},
		},
		{
			testName:     "check_log_warning",
			winOK:        true,
			interpreter:  interpreter,
			env:          nil,
			dir:          "/tmp",
			command:      `echo "[WARN] Test warn"`,
			logger:       &testLogger{},
			hasErr:       false,
			loggerResult: &testLogger{warnings: []string{"Test warn"}},
		},
		{
			testName:     "check_log_info",
			winOK:        true,
			interpreter:  interpreter,
			env:          nil,
			dir:          "/tmp",
			command:      `echo "[INFO] Test info"`,
			logger:       &testLogger{},
			hasErr:       false,
			loggerResult: &testLogger{infos: []string{"Test info"}},
		},
		{
			testName:     "check_log_debug",
			winOK:        true,
			interpreter:  interpreter,
			env:          nil,
			dir:          "/tmp",
			command:      `echo "[DEBUG] Test debug"`,
			logger:       &testLogger{},
			hasErr:       false,
			loggerResult: &testLogger{debugs: []string{"Test debug"}},
		},
		{
			testName:     "check_log_trace",
			winOK:        true,
			interpreter:  interpreter,
			env:          nil,
			dir:          "/tmp",
			command:      `echo "[TRACE] Test trace"`,
			logger:       &testLogger{},
			hasErr:       false,
			loggerResult: &testLogger{traces: []string{"Test trace"}},
		},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			if runtime.GOOS == "windows" && !d.winOK {
				t.Skip("Test is not valid on Windows")
			}

			ctx := t.Context()
			err := RunCommand(ctx, d.interpreter, d.env, d.dir, d.command, d.logger)

			hasErr := err != nil
			if hasErr != d.hasErr {
				t.Errorf("unexpected error state")
			}

			if d.logger != nil {
				if !reflect.DeepEqual(d.logger.errors, d.loggerResult.errors) {
					t.Errorf("expected errors %v, got %v", d.loggerResult.errors, d.logger.errors)
				}

				if !reflect.DeepEqual(d.logger.warnings, d.loggerResult.warnings) {
					t.Errorf("expected warnings %v, got %v", d.loggerResult.warnings, d.logger.warnings)
				}

				if !reflect.DeepEqual(d.logger.infos, d.loggerResult.infos) {
					t.Errorf("expected infos %v, got %v", d.loggerResult.infos, d.logger.infos)
				}

				if !reflect.DeepEqual(d.logger.debugs, d.loggerResult.debugs) {
					t.Errorf("expected debugs %v, got %v", d.loggerResult.debugs, d.logger.debugs)
				}

				if !reflect.DeepEqual(d.logger.traces, d.loggerResult.traces) {
					t.Errorf("expected traces %v, got %v", d.loggerResult.traces, d.logger.traces)
				}
			}
		})
	}
}
