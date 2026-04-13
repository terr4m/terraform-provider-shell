package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"shell": providerserver.NewProtocol6WithError(New("test", "test")()),
}

func testAccPreCheck(_ *testing.T) {
}

// mustStringList creates a types.List of strings for testing.
func mustStringList(t *testing.T, values []string) types.List {
	t.Helper()

	elems := make([]types.String, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}

	list, diags := types.ListValueFrom(t.Context(), types.StringType, elems)
	if diags.HasError() {
		t.Fatalf("failed to create list: %v", diags.Errors())
	}

	return list
}

// mustStringMap creates a types.Map of strings for testing.
func mustStringMap(t *testing.T, values map[string]string) types.Map {
	t.Helper()

	elems := make(map[string]types.String, len(values))
	for k, v := range values {
		elems[k] = types.StringValue(v)
	}

	m, diags := types.MapValueFrom(t.Context(), types.StringType, elems)
	if diags.HasError() {
		t.Fatalf("failed to create map: %v", diags.Errors())
	}

	return m
}
