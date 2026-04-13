package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resolveInterpreter resolves the interpreter from the TF type or falls back to the default.
func resolveInterpreter(ctx context.Context, tfInterpreter types.List, defaultInterpreter []string) ([]string, diag.Diagnostics) {
	if !tfInterpreter.IsNull() {
		var interpreter []string
		diags := tfInterpreter.ElementsAs(ctx, &interpreter, false)
		return interpreter, diags
	}

	return defaultInterpreter, nil
}

// resolveEnvironment resolves the environment by merging the default and TF map.
func resolveEnvironment(ctx context.Context, tfEnvironment types.Map, defaultEnvironment map[string]string) (map[string]string, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	environment := make(map[string]string, len(defaultEnvironment))
	maps.Copy(environment, defaultEnvironment)

	if !tfEnvironment.IsNull() {
		resourceEnv := map[string]string{}
		if diags.Append(tfEnvironment.ElementsAs(ctx, &resourceEnv, false)...); diags.HasError() {
			return nil, diags
		}

		maps.Copy(environment, resourceEnv)
	}

	return environment, diags
}
