package script_test

import (
	"context"
	"encoding/json"
	"os"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/terr4m/terraform-provider-shell/internal/script"
	"github.com/terr4m/terraform-provider-shell/internal/shell"
)

func TestNewCommandRunner(t *testing.T) {
	t.Parallel()

	t.Run("nil_log_provider", func(t *testing.T) {
		t.Parallel()

		runner := script.NewCommandRunner(nil)
		if runner == nil {
			t.Fatal("expected non-nil runner")
		}
	})

	t.Run("with_log_provider", func(t *testing.T) {
		t.Parallel()

		runner := script.NewCommandRunner(&shell.LogProvider{Logger: &script.TFLogLogger{}})
		if runner == nil {
			t.Fatal("expected non-nil runner")
		}
	})
}

func TestShellCommandRunner_Run(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	for _, d := range []struct {
		testName       string
		opts           script.RunOptions
		wantResult     script.RunResult
		wantError      bool
		wantErrorCount int
	}{
		{
			testName: "invalid_interpreter",
			opts: script.RunOptions{
				Interpreter: []string{"xxxxxx"},
				Command:     `echo "hello"`,
				Lifecycle:   script.LifecycleRead,
			},
			wantResult:     script.RunResult{},
			wantError:      true,
			wantErrorCount: 1,
		},
		{
			testName: "command_exit_code_1",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     testExitCommand(1),
				Lifecycle:   script.LifecycleCreate,
			},
			wantResult:     script.RunResult{},
			wantError:      true,
			wantErrorCount: 1,
		},
		{
			testName: "command_exit_code_1_with_error_file",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     testWriteErrorCommand("something went wrong"),
				Lifecycle:   script.LifecycleCreate,
			},
			wantResult:     script.RunResult{},
			wantError:      true,
			wantErrorCount: 1,
		},
		{
			testName: "read_json_false_no_output",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     `echo "ignored"`,
				Lifecycle:   script.LifecycleDelete,
				ReadJSON:    false,
			},
			wantResult:     script.RunResult{},
			wantError:      false,
			wantErrorCount: 0,
		},
		{
			testName: "read_json_true_missing_output_file",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     `echo "no output written"`,
				Lifecycle:   script.LifecycleRead,
				ReadJSON:    true,
			},
			wantResult:     script.RunResult{},
			wantError:      true,
			wantErrorCount: 1,
		},
		{
			testName: "read_json_true_with_output",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     testWriteOutputCommand(`{"key":"value"}`),
				Lifecycle:   script.LifecycleRead,
				ReadJSON:    true,
			},
			wantResult: script.RunResult{
				Meta:   script.ResultMetadata{},
				Output: map[string]any{"key": "value"},
			},
			wantError:      false,
			wantErrorCount: 0,
		},
		{
			testName: "read_json_true_with_metadata",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     testWriteOutputCommand(`{"key":"value","__meta":{"output_drift_detected":true}}`),
				Lifecycle:   script.LifecycleRead,
				ReadJSON:    true,
			},
			wantResult: script.RunResult{
				Meta:   script.ResultMetadata{OutputDriftDetected: true},
				Output: map[string]any{"key": "value"},
			},
			wantError:      false,
			wantErrorCount: 0,
		},
		{
			testName: "with_inputs",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     testWriteOutputCommand(`{"received":true}`),
				Lifecycle:   script.LifecycleCreate,
				Inputs:      map[string]any{"foo": "bar"},
				ReadJSON:    true,
			},
			wantResult: script.RunResult{
				Meta:   script.ResultMetadata{},
				Output: map[string]any{"received": true},
			},
			wantError:      false,
			wantErrorCount: 0,
		},
		{
			testName: "with_state_output",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     testWriteOutputCommand(`{"updated":true}`),
				Lifecycle:   script.LifecycleUpdate,
				StateOutput: map[string]any{"old": "state"},
				ReadJSON:    true,
			},
			wantResult: script.RunResult{
				Meta:   script.ResultMetadata{},
				Output: map[string]any{"updated": true},
			},
			wantError:      false,
			wantErrorCount: 0,
		},
		{
			testName: "with_environment",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Environment: map[string]string{"CUSTOM_VAR": "custom_value"},
				Command:     testWriteOutputCommand(`{"env":true}`),
				Lifecycle:   script.LifecycleRead,
				ReadJSON:    true,
			},
			wantResult: script.RunResult{
				Meta:   script.ResultMetadata{},
				Output: map[string]any{"env": true},
			},
			wantError:      false,
			wantErrorCount: 0,
		},
		{
			testName: "unmarshalable_inputs",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     `echo "should not run"`,
				Lifecycle:   script.LifecycleCreate,
				Inputs:      make(chan int),
				ReadJSON:    true,
			},
			wantResult:     script.RunResult{},
			wantError:      true,
			wantErrorCount: 1,
		},
		{
			testName: "unmarshalable_state_output",
			opts: script.RunOptions{
				Interpreter: interpreter,
				Command:     `echo "should not run"`,
				Lifecycle:   script.LifecycleUpdate,
				StateOutput: make(chan int),
				ReadJSON:    true,
			},
			wantResult:     script.RunResult{},
			wantError:      true,
			wantErrorCount: 1,
		},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			runner := script.NewCommandRunner(nil)

			got, diags := runner.Run(ctx, d.opts)

			if diags.HasError() != d.wantError {
				t.Errorf("expected error=%v, got diags: %v", d.wantError, diags.Errors())
			}

			if d.wantErrorCount > 0 && diags.ErrorsCount() != d.wantErrorCount {
				t.Errorf("expected %d errors, got %d", d.wantErrorCount, diags.ErrorsCount())
			}

			if !d.wantError {
				if diff := cmp.Diff(d.wantResult, got); diff != "" {
					t.Errorf("Run() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestShellCommandRunner_Run_LifecycleEnv(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	for _, d := range []struct {
		testName  string
		lifecycle script.Lifecycle
		want      string
	}{
		{testName: "plan", lifecycle: script.LifecyclePlan, want: "plan"},
		{testName: "create", lifecycle: script.LifecycleCreate, want: "create"},
		{testName: "read", lifecycle: script.LifecycleRead, want: "read"},
		{testName: "update", lifecycle: script.LifecycleUpdate, want: "update"},
		{testName: "delete", lifecycle: script.LifecycleDelete, want: "delete"},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			runner := script.NewCommandRunner(nil)

			var cmd string
			if runtime.GOOS == "windows" {
				cmd = `[IO.File]::WriteAllText($env:TF_SCRIPT_OUTPUT, ('"' + $env:TF_SCRIPT_LIFECYCLE + '"'))`
			} else {
				cmd = `printf '"%s"' "${TF_SCRIPT_LIFECYCLE}" > "${TF_SCRIPT_OUTPUT}"`
			}

			res, diags := runner.Run(ctx, script.RunOptions{
				Interpreter: interpreter,
				Command:     cmd,
				Lifecycle:   d.lifecycle,
				ReadJSON:    true,
			})
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags.Errors())
			}

			got, ok := res.Output.(string)
			if !ok {
				t.Fatalf("expected string output, got %T: %v", res.Output, res.Output)
			}

			if got != d.want {
				t.Errorf("expected lifecycle %q, got %q", d.want, got)
			}
		})
	}
}

func TestShellCommandRunner_Run_InputsEnv(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	var cmd string
	if runtime.GOOS == "windows" {
		cmd = `[IO.File]::WriteAllText($env:TF_SCRIPT_OUTPUT, $env:TF_SCRIPT_INPUTS)`
	} else {
		cmd = `printf '%s' "${TF_SCRIPT_INPUTS}" > "${TF_SCRIPT_OUTPUT}"`
	}

	ctx := t.Context()
	runner := script.NewCommandRunner(nil)

	inputs := map[string]any{"name": "test", "count": float64(42)}

	res, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Command:     cmd,
		Lifecycle:   script.LifecycleRead,
		Inputs:      inputs,
		ReadJSON:    true,
	})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags.Errors())
	}

	got, ok := res.Output.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", res.Output)
	}

	if diff := cmp.Diff(inputs, got); diff != "" {
		t.Errorf("inputs mismatch (-want +got):\n%s", diff)
	}
}

func TestShellCommandRunner_Run_StateOutputEnv(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	var cmd string
	if runtime.GOOS == "windows" {
		cmd = `[IO.File]::WriteAllText($env:TF_SCRIPT_OUTPUT, $env:TF_SCRIPT_STATE_OUTPUT)`
	} else {
		cmd = `printf '%s' "${TF_SCRIPT_STATE_OUTPUT}" > "${TF_SCRIPT_OUTPUT}"`
	}

	ctx := t.Context()
	runner := script.NewCommandRunner(nil)

	stateOutput := map[string]any{"existing": "state"}

	res, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Command:     cmd,
		Lifecycle:   script.LifecycleUpdate,
		StateOutput: stateOutput,
		ReadJSON:    true,
	})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags.Errors())
	}

	got, ok := res.Output.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", res.Output)
	}

	if diff := cmp.Diff(stateOutput, got); diff != "" {
		t.Errorf("state output mismatch (-want +got):\n%s", diff)
	}
}

func TestShellCommandRunner_Run_EnvironmentMerge(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	var cmd string
	if runtime.GOOS == "windows" {
		cmd = `$j = (@{a=$env:CUSTOM_A; b=$env:CUSTOM_B} | ConvertTo-Json -Compress); [IO.File]::WriteAllText($env:TF_SCRIPT_OUTPUT, $j)`
	} else {
		cmd = `printf '{"a":"%s","b":"%s"}' "${CUSTOM_A}" "${CUSTOM_B}" > "${TF_SCRIPT_OUTPUT}"`
	}

	ctx := t.Context()
	runner := script.NewCommandRunner(nil)

	res, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Environment: map[string]string{"CUSTOM_A": "alpha", "CUSTOM_B": "beta"},
		Command:     cmd,
		Lifecycle:   script.LifecycleRead,
		ReadJSON:    true,
	})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags.Errors())
	}

	got, ok := res.Output.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", res.Output)
	}

	want := map[string]any{"a": "alpha", "b": "beta"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("environment mismatch (-want +got):\n%s", diff)
	}
}

func TestShellCommandRunner_Run_OutputFileCleaned(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	// Run a command that writes to the output file, then check the temp files are cleaned up.
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = `$p = $env:TF_SCRIPT_OUTPUT; [IO.File]::WriteAllText($p, ('{"path":"' + $p.Replace('\', '\\') + '"}'))`
	} else {
		cmd = `printf '{"path":"%s"}' "${TF_SCRIPT_OUTPUT}" > "${TF_SCRIPT_OUTPUT}"`
	}

	ctx := t.Context()
	runner := script.NewCommandRunner(nil)

	res, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Command:     cmd,
		Lifecycle:   script.LifecycleRead,
		ReadJSON:    true,
	})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags.Errors())
	}

	outMap, ok := res.Output.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", res.Output)
	}

	outPath, ok := outMap["path"].(string)
	if !ok {
		t.Fatal("expected path in output")
	}

	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Errorf("expected output file %q to be cleaned up", outPath)
	}
}

func TestShellCommandRunner_Run_CancelledContext(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	runner := script.NewCommandRunner(nil)

	_, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Command:     `echo "should not run"`,
		Lifecycle:   script.LifecycleRead,
	})

	if !diags.HasError() {
		t.Error("expected error from cancelled context")
	}
}

func TestShellCommandRunner_Run_InvalidJSON(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	ctx := t.Context()
	runner := script.NewCommandRunner(nil)

	_, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Command:     testWriteOutputCommand("not-json"),
		Lifecycle:   script.LifecycleRead,
		ReadJSON:    true,
	})

	if !diags.HasError() {
		t.Error("expected error for invalid JSON output")
	}
}

func TestShellCommandRunner_Run_NilInputsAndStateOutput(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()

	// Verify that TF_SCRIPT_INPUTS and TF_SCRIPT_STATE_OUTPUT are NOT set when nil.
	var cmd string
	if runtime.GOOS == "windows" {
		cmd = `$r = @{}; if ($env:TF_SCRIPT_INPUTS) { $r.inputs = 'set' } else { $r.inputs = 'unset' }; if ($env:TF_SCRIPT_STATE_OUTPUT) { $r.state = 'set' } else { $r.state = 'unset' }; [IO.File]::WriteAllText($env:TF_SCRIPT_OUTPUT, ($r | ConvertTo-Json -Compress))`
	} else {
		cmd = `printf '{"inputs":"%s","state":"%s"}' "${TF_SCRIPT_INPUTS:-unset}" "${TF_SCRIPT_STATE_OUTPUT:-unset}" > "${TF_SCRIPT_OUTPUT}"`
	}

	ctx := t.Context()
	runner := script.NewCommandRunner(nil)

	res, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Command:     cmd,
		Lifecycle:   script.LifecycleRead,
		ReadJSON:    true,
	})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags.Errors())
	}

	got, ok := res.Output.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", res.Output)
	}

	if got["inputs"] != "unset" {
		t.Errorf("expected inputs=unset, got %q", got["inputs"])
	}
	if got["state"] != "unset" {
		t.Errorf("expected state=unset, got %q", got["state"])
	}
}

func TestShellCommandRunner_Run_InputsJSON(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()
	ctx := t.Context()
	runner := script.NewCommandRunner(nil)

	inputs := map[string]any{"nested": map[string]any{"key": "val"}, "list": []any{"a", "b"}}
	wantBytes, _ := json.Marshal(inputs)

	var cmd string
	if runtime.GOOS == "windows" {
		cmd = `[IO.File]::WriteAllText($env:TF_SCRIPT_OUTPUT, $env:TF_SCRIPT_INPUTS)`
	} else {
		cmd = `printf '%s' "${TF_SCRIPT_INPUTS}" > "${TF_SCRIPT_OUTPUT}"`
	}

	res, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Command:     cmd,
		Lifecycle:   script.LifecycleRead,
		Inputs:      inputs,
		ReadJSON:    true,
	})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags.Errors())
	}

	gotBytes, _ := json.Marshal(res.Output)

	var want, got any
	_ = json.Unmarshal(wantBytes, &want)
	_ = json.Unmarshal(gotBytes, &got)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("complex inputs mismatch (-want +got):\n%s", diff)
	}
}

func TestShellCommandRunner_Run_WithLogProvider(t *testing.T) {
	t.Parallel()

	interpreter := testInterpreter()
	logger := &mockLogger{}
	runner := script.NewCommandRunner(&shell.LogProvider{Logger: logger})

	var cmd string
	if runtime.GOOS == "windows" {
		cmd = `Write-Output "[INFO] hello from script"; [IO.File]::WriteAllText($env:TF_SCRIPT_OUTPUT, '{"logged":true}')`
	} else {
		cmd = `echo "[INFO] hello from script"; printf '{"logged":true}' > "${TF_SCRIPT_OUTPUT}"`
	}

	ctx := t.Context()
	res, diags := runner.Run(ctx, script.RunOptions{
		Interpreter: interpreter,
		Command:     cmd,
		Lifecycle:   script.LifecycleRead,
		ReadJSON:    true,
	})
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags.Errors())
	}

	got, ok := res.Output.(map[string]any)
	if !ok {
		t.Fatalf("expected map output, got %T", res.Output)
	}

	if got["logged"] != true {
		t.Errorf("expected logged=true, got %v", got["logged"])
	}

	entries := logger.getEntries()
	if len(entries) == 0 {
		t.Fatal("expected at least one log entry")
	}

	found := false
	for _, e := range entries {
		if e.level == "info" && e.msg == "hello from script" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected INFO log entry with 'hello from script', got: %v", entries)
	}
}
