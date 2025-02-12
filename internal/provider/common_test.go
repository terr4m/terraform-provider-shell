package provider

import (
	"testing"

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
