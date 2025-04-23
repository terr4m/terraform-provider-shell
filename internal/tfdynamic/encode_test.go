package tfdynamic

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestEncodeDynamicObject(t *testing.T) {
	t.Parallel()

	simpleObject, _ := types.ObjectValue(map[string]attr.Type{"foo": types.StringType}, map[string]attr.Value{"foo": types.StringValue("bar")})

	for _, d := range []struct {
		testName string
		dyn      types.Dynamic
		expected any
		errMsg   string
	}{
		{
			testName: "unknown",
			dyn:      types.DynamicUnknown(),
			expected: nil,
			errMsg:   "underlying value is unknown",
		},
		{
			testName: "null",
			dyn:      types.DynamicValue(types.ObjectNull(nil)),
			expected: nil,
			errMsg:   "",
		},
		{
			testName: "object",
			dyn:      types.DynamicValue(simpleObject),
			expected: map[string]any{"foo": "bar"},
			errMsg:   "",
		},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			actual, err := EncodeDynamic(ctx, d.dyn)

			if !reflect.DeepEqual(actual, d.expected) {
				t.Errorf("expected %v, got %v", d.expected, actual)
			}

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != d.errMsg {
				t.Errorf("expected error message %s, got %s", d.errMsg, errMsg)
			}
		})
	}
}
