package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_resolveInterpreter(t *testing.T) {
	t.Parallel()

	defaultInterpreter := []string{"/bin/bash", "-c"}

	for _, d := range []struct {
		testName           string
		tfInterpreter      types.List
		defaultInterpreter []string
		want               []string
		wantError          bool
	}{
		{
			testName:           "null_uses_default",
			tfInterpreter:      types.ListNull(types.StringType),
			defaultInterpreter: defaultInterpreter,
			want:               defaultInterpreter,
			wantError:          false,
		},
		{
			testName:           "set_overrides_default",
			tfInterpreter:      mustStringList(t, []string{"/bin/sh", "-c"}),
			defaultInterpreter: defaultInterpreter,
			want:               []string{"/bin/sh", "-c"},
			wantError:          false,
		},
		{
			testName:           "single_element",
			tfInterpreter:      mustStringList(t, []string{"pwsh"}),
			defaultInterpreter: defaultInterpreter,
			want:               []string{"pwsh"},
			wantError:          false,
		},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			got, diags := resolveInterpreter(ctx, d.tfInterpreter, d.defaultInterpreter)

			if diags.HasError() != d.wantError {
				t.Errorf("expected error=%v, got diags: %v", d.wantError, diags.Errors())
			}

			if !d.wantError {
				if diff := cmp.Diff(d.want, got); diff != "" {
					t.Errorf("resolveInterpreter() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func Test_resolveEnvironment(t *testing.T) {
	t.Parallel()

	for _, d := range []struct {
		testName           string
		tfEnvironment      types.Map
		defaultEnvironment map[string]string
		want               map[string]string
		wantError          bool
	}{
		{
			testName:           "null_uses_default",
			tfEnvironment:      types.MapNull(types.StringType),
			defaultEnvironment: map[string]string{"A": "1"},
			want:               map[string]string{"A": "1"},
			wantError:          false,
		},
		{
			testName:           "null_empty_default",
			tfEnvironment:      types.MapNull(types.StringType),
			defaultEnvironment: map[string]string{},
			want:               map[string]string{},
			wantError:          false,
		},
		{
			testName:           "null_nil_default",
			tfEnvironment:      types.MapNull(types.StringType),
			defaultEnvironment: nil,
			want:               map[string]string{},
			wantError:          false,
		},
		{
			testName:           "set_merges_with_default",
			tfEnvironment:      mustStringMap(t, map[string]string{"B": "2"}),
			defaultEnvironment: map[string]string{"A": "1"},
			want:               map[string]string{"A": "1", "B": "2"},
			wantError:          false,
		},
		{
			testName:           "set_overrides_default",
			tfEnvironment:      mustStringMap(t, map[string]string{"A": "override"}),
			defaultEnvironment: map[string]string{"A": "1", "B": "2"},
			want:               map[string]string{"A": "override", "B": "2"},
			wantError:          false,
		},
		{
			testName:           "set_empty_default",
			tfEnvironment:      mustStringMap(t, map[string]string{"X": "y"}),
			defaultEnvironment: map[string]string{},
			want:               map[string]string{"X": "y"},
			wantError:          false,
		},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			got, diags := resolveEnvironment(ctx, d.tfEnvironment, d.defaultEnvironment)

			if diags.HasError() != d.wantError {
				t.Errorf("expected error=%v, got diags: %v", d.wantError, diags.Errors())
			}

			if !d.wantError {
				if diff := cmp.Diff(d.want, got); diff != "" {
					t.Errorf("resolveEnvironment() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestScriptResource_Configure_NilProviderData(t *testing.T) {
	t.Parallel()

	r := &ScriptResource{}
	resp := &resource.ConfigureResponse{}
	r.Configure(t.Context(), resource.ConfigureRequest{ProviderData: nil}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error for nil provider data, got: %v", resp.Diagnostics.Errors())
	}

	if r.providerData != nil {
		t.Error("expected providerData to remain nil")
	}
}

func TestScriptResource_Configure_WrongType(t *testing.T) {
	t.Parallel()

	r := &ScriptResource{}
	resp := &resource.ConfigureResponse{}
	r.Configure(t.Context(), resource.ConfigureRequest{ProviderData: "wrong-type"}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for wrong provider data type")
	}

	if resp.Diagnostics.Errors()[0].Summary() != "Unexpected resource provider data." {
		t.Errorf("unexpected error summary: %s", resp.Diagnostics.Errors()[0].Summary())
	}
}

func TestScriptDataSource_Configure_NilProviderData(t *testing.T) {
	t.Parallel()

	d := &ScriptDataSource{}
	resp := &datasource.ConfigureResponse{}
	d.Configure(t.Context(), datasource.ConfigureRequest{ProviderData: nil}, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no error for nil provider data, got: %v", resp.Diagnostics.Errors())
	}

	if d.providerData != nil {
		t.Error("expected providerData to remain nil")
	}
}

func TestScriptDataSource_Configure_WrongType(t *testing.T) {
	t.Parallel()

	d := &ScriptDataSource{}
	resp := &datasource.ConfigureResponse{}
	d.Configure(t.Context(), datasource.ConfigureRequest{ProviderData: 42}, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for wrong provider data type")
	}

	if resp.Diagnostics.Errors()[0].Summary() != "Unexpected data source provider data." {
		t.Errorf("unexpected error summary: %s", resp.Diagnostics.Errors()[0].Summary())
	}
}
