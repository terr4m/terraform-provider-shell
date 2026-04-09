package script

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/terr4m/terraform-provider-shell/internal/shell"
)

// RunOptions contains the options for running a command.
type RunOptions struct {
	Interpreter      []string
	Environment      map[string]string
	WorkingDirectory string
	Command          string
	Lifecycle        Lifecycle
	Inputs           any
	StateOutput      any
	ReadJSON         bool
}

// RunResult represents the result of running a command.
type RunResult struct {
	Meta   ResultMetadata
	Output any
}

// ResultMetadata represents metadata from running a command.
type ResultMetadata struct {
	OutputDriftDetected bool `json:"output_drift_detected"`
}

// CommandRunner runs shell scripts and returns parsed results.
type CommandRunner interface {
	Run(ctx context.Context, opts RunOptions) (RunResult, diag.Diagnostics)
}

type shellCommandRunner struct {
	logProvider *shell.LogProvider
}

// NewCommandRunner creates a new CommandRunner.
func NewCommandRunner(logProvider *shell.LogProvider) CommandRunner {
	return &shellCommandRunner{logProvider: logProvider}
}

// Run runs a shell script with the given options and returns the result.
func (r *shellCommandRunner) Run(ctx context.Context, opts RunOptions) (RunResult, diag.Diagnostics) {
	var diags diag.Diagnostics
	var res RunResult

	outFilePath, err := shell.GetOutFilePath()
	if err != nil {
		diags.AddError("Failed to get output file path.", err.Error())
		return res, diags
	}
	defer os.Remove(outFilePath)

	errorFilePath, err := shell.GetErrorFilePath()
	if err != nil {
		diags.AddError("Failed to get error file path.", err.Error())
		return res, diags
	}
	defer os.Remove(errorFilePath)

	environment := make(map[string]string, len(opts.Environment)+5)
	maps.Copy(environment, opts.Environment)

	environment[LifecycleEnv] = string(opts.Lifecycle)
	environment[ScriptOutputFilePathEnv] = outFilePath
	environment[ScriptErrorFilePathEnv] = errorFilePath

	if opts.Inputs != nil {
		by, err := json.Marshal(opts.Inputs)
		if err != nil {
			diags.AddError("Failed to marshal inputs.", err.Error())
			return res, diags
		}

		environment[InputsEnv] = string(by)
	}

	if opts.StateOutput != nil {
		by, err := json.Marshal(opts.StateOutput)
		if err != nil {
			diags.AddError("Failed to marshal state output.", err.Error())
			return res, diags
		}

		environment[StateOutputEnv] = string(by)
	}

	err = shell.RunCommand(ctx, opts.Interpreter, environment, opts.WorkingDirectory, opts.Command, r.logProvider)
	if err != nil {
		exitError := &exec.ExitError{}
		if errors.As(err, &exitError) {
			detail := ""
			by, err := os.ReadFile(errorFilePath)
			if err == nil {
				detail = string(by)
			}
			diags.AddError(fmt.Sprintf("Command failed with exit code: %d", exitError.ExitCode()), detail)
			return res, diags
		}

		diags.AddError("Failed to run command.", err.Error())
		return res, diags
	}

	if !opts.ReadJSON {
		return res, diags
	}

	out, err := shell.ReadJSON(outFilePath)
	if err != nil {
		diags.AddError("Failed to read output file.", err.Error())
		return res, diags
	}

	res = GetRunCommandResult(out)

	return res, diags
}

// GetRunCommandResult extracts the RunResult from the output.
func GetRunCommandResult(o any) RunResult {
	if om, ok := o.(map[string]any); ok {
		if m, ok := om["__meta"].(map[string]any); ok {
			meta := ResultMetadata{}
			if outputDriftDetected, ok := m["output_drift_detected"].(bool); ok {
				meta.OutputDriftDetected = outputDriftDetected
			}
			delete(om, "__meta")
			return RunResult{
				Meta:   meta,
				Output: om,
			}
		}
	}

	return RunResult{
		Output: o,
	}
}
