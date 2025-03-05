package provider

import (
	"context"
	"encoding/json"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terr4m/terraform-provider-shell/internal/shell"
)

const (
	ScriptOutPutFilePathEnv string = "TF_SCRIPT_OUTPUT"
	StateOutputEnv          string = "TF_STATE_OUTPUT"
)

// runCommand runs a shell script and returns the output.
func runCommand(ctx context.Context, providerData *ShellProviderData, tfInterpreter types.List, tfEnvironment types.Map, tfWorkingDirectory, tfCommand types.String, stateOutput any, readJSON bool) (any, diag.Diagnostics) {
	var diags diag.Diagnostics

	var interpreter []string
	if !tfInterpreter.IsNull() {
		if diags.Append(tfInterpreter.ElementsAs(ctx, &interpreter, false)...); diags.HasError() {
			return nil, diags
		}
	} else {
		interpreter = providerData.Interpreter
	}

	environment := map[string]string{}
	for k, v := range providerData.Environment {
		environment[k] = v
	}

	if !tfEnvironment.IsNull() {
		if diags.Append(tfEnvironment.ElementsAs(ctx, &environment, false)...); diags.HasError() {
			return nil, diags
		}
	}

	outFilePath, err := shell.GetOutFilePath()
	if err != nil {
		diags.AddError("Failed to get output file path.", err.Error())
		return nil, diags
	}
	defer os.Remove(outFilePath)

	environment[ScriptOutPutFilePathEnv] = outFilePath

	if stateOutput != nil {
		by, err := json.Marshal(stateOutput)
		if err != nil {
			diags.AddError("Failed to marshal state output.", err.Error())
			return nil, diags
		}

		environment[StateOutputEnv] = string(by)
	}

	var log *tflogLogger
	if providerData.LogOutput {
		log = &tflogLogger{}
	}

	err = shell.RunCommand(ctx, interpreter, environment, tfWorkingDirectory.ValueString(), tfCommand.ValueString(), log)
	if err != nil {
		diags.AddError("Failed to run command.", err.Error())
		return nil, diags
	}

	if !readJSON {
		return nil, diags
	}

	a, err := shell.ReadOutJSON(outFilePath)
	if err != nil {
		diags.AddError("Failed to read output file.", err.Error())
		return nil, diags
	}

	return a, diags
}
