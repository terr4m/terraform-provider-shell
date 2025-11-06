package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terr4m/terraform-provider-shell/internal/shell"
)

const (
	LifecycleEnv            string = "TF_SCRIPT_LIFECYCLE"
	InputsEnv               string = "TF_SCRIPT_INPUTS"
	StateOutputEnv          string = "TF_SCRIPT_STATE_OUTPUT"
	ScriptOutputFilePathEnv string = "TF_SCRIPT_OUTPUT"
	ScriptErrorFilePathEnv  string = "TF_SCRIPT_ERROR"
)

type TFLifecycle string

const (
	TFLifecyclePlan   TFLifecycle = "plan"
	TFLifecycleCreate TFLifecycle = "create"
	TFLifecycleRead   TFLifecycle = "read"
	TFLifecycleUpdate TFLifecycle = "update"
	TFLifecycleDelete TFLifecycle = "delete"
)

// runCommandResult represents the result of running a command.
type runCommandResult struct {
	Meta   runCommandResultMetadata
	Output any
}

// runCommandResultMetadata represents metadata from running a command.
type runCommandResultMetadata struct {
	OutputDriftDetected bool `json:"output_drift_detected"`
}

// runCommand runs a shell script and returns the output.
func runCommand(ctx context.Context, providerData *ShellProviderData, tfInterpreter types.List, tfEnvironment types.Map, tfWorkingDirectory, tfCommand types.String, lifecycle TFLifecycle, inputs, stateOutput any, readJSON bool) (runCommandResult, diag.Diagnostics) {
	var diags diag.Diagnostics
	var res runCommandResult

	var interpreter []string
	if !tfInterpreter.IsNull() {
		if diags.Append(tfInterpreter.ElementsAs(ctx, &interpreter, false)...); diags.HasError() {
			return res, diags
		}
	} else {
		interpreter = providerData.DefaultInterpreter
	}

	environment := map[string]string{}
	for k, v := range providerData.Environment {
		environment[k] = v
	}

	if !tfEnvironment.IsNull() {
		if diags.Append(tfEnvironment.ElementsAs(ctx, &environment, false)...); diags.HasError() {
			return res, diags
		}
	}

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

	environment[LifecycleEnv] = string(lifecycle)
	environment[ScriptOutputFilePathEnv] = outFilePath
	environment[ScriptErrorFilePathEnv] = errorFilePath

	if inputs != nil {
		by, err := json.Marshal(inputs)
		if err != nil {
			diags.AddError("Failed to marshal inputs.", err.Error())
			return res, diags
		}

		environment[InputsEnv] = string(by)
	}

	if stateOutput != nil {
		by, err := json.Marshal(stateOutput)
		if err != nil {
			diags.AddError("Failed to marshal state output.", err.Error())
			return res, diags
		}

		environment[StateOutputEnv] = string(by)
	}

	var logProvider *shell.LogProvider
	if providerData.LogOutput {
		logProvider = &shell.LogProvider{
			Logger: &tflogLogger{},
		}
	}

	err = shell.RunCommand(ctx, interpreter, environment, tfWorkingDirectory.ValueString(), tfCommand.ValueString(), logProvider)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
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

	if !readJSON {
		return res, diags
	}

	out, err := shell.ReadJSON(outFilePath)
	if err != nil {
		diags.AddError("Failed to read output file.", err.Error())
		return res, diags
	}

	res = getRunCommandResult(out)

	return res, diags
}

// getRunCommandResult extracts the runCommandResult from the output.
func getRunCommandResult(o any) runCommandResult {
	if om, ok := o.(map[string]any); ok {
		if m, ok := om["__meta"].(map[string]any); ok {
			meta := runCommandResultMetadata{}
			if outputDriftDetected, ok := m["output_drift_detected"].(bool); ok {
				meta.OutputDriftDetected = outputDriftDetected
			}
			delete(om, "__meta")
			return runCommandResult{
				Meta:   meta,
				Output: om,
			}
		}
	}

	return runCommandResult{
		Output: o,
	}
}
