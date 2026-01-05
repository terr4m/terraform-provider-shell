package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_runScript(t *testing.T) {
	t.Parallel()

	for _, d := range []struct {
		testName     string
		providerData *ShellProviderData
		interpreter  types.List
		env          types.Map
		dir          types.String
		command      types.String
		readJSON     bool
	}{} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()
		})
	}
}

func Test_getRunCommandResult(t *testing.T) {
	t.Parallel()

	for _, d := range []struct {
		testName string
		output   any
		want     runCommandResult
	}{
		{
			testName: "empty",
			output:   map[string]any{},
			want: runCommandResult{
				Meta:   runCommandResultMetadata{},
				Output: map[string]any{},
			},
		},
		{
			testName: "no_metadata",
			output:   map[string]any{"foo": "bar"},
			want: runCommandResult{
				Meta:   runCommandResultMetadata{},
				Output: map[string]any{"foo": "bar"},
			},
		},
		{
			testName: "metadata_empty",
			output:   map[string]any{"foo": "bar", "__meta": map[string]any{}},
			want: runCommandResult{
				Meta:   runCommandResultMetadata{},
				Output: map[string]any{"foo": "bar"},
			},
		},
		{
			testName: "metadata_unknown",
			output:   map[string]any{"foo": "bar", "__meta": map[string]any{"unknown_key": 123}},
			want: runCommandResult{
				Meta:   runCommandResultMetadata{},
				Output: map[string]any{"foo": "bar"},
			},
		},
		{
			testName: "metadata_drift_detected",
			output:   map[string]any{"foo": "bar", "__meta": map[string]any{"output_drift_detected": true}},
			want: runCommandResult{
				Meta:   runCommandResultMetadata{OutputDriftDetected: true},
				Output: map[string]any{"foo": "bar"},
			},
		},
		{
			testName: "metadata_no_drift_detected",
			output:   map[string]any{"foo": "bar", "__meta": map[string]any{"output_drift_detected": false}},
			want: runCommandResult{
				Meta:   runCommandResultMetadata{OutputDriftDetected: false},
				Output: map[string]any{"foo": "bar"},
			},
		},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			got := getRunCommandResult(d.output)

			if diff := cmp.Diff(d.want, got); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
