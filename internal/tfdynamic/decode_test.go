package tfdynamic

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	emptyObject, _ := basetypes.NewObjectValue(map[string]attr.Type{}, map[string]attr.Value{})
	simpleObject, _ := basetypes.NewObjectValue(map[string]attr.Type{"foo": basetypes.StringType{}}, map[string]attr.Value{"foo": basetypes.NewStringValue("bar")})

	emptyTuple, _ := basetypes.NewTupleValue([]attr.Type{}, []attr.Value{})
	stringTuple, _ := basetypes.NewTupleValue([]attr.Type{basetypes.StringType{}, basetypes.StringType{}}, []attr.Value{basetypes.NewStringValue("foo"), basetypes.NewStringValue("bar")})

	for _, d := range []struct {
		testName string
		obj      any
		expected basetypes.DynamicValue
		errMsg   string
	}{
		{
			testName: "unexpected_type",
			obj:      1,
			expected: basetypes.DynamicValue{},
			errMsg:   "Unexpected type.",
		},
		{
			testName: "null",
			obj:      nil,
			expected: basetypes.NewDynamicNull(),
			errMsg:   "",
		},
		{
			testName: "int64",
			obj:      int64(1),
			expected: basetypes.NewDynamicValue(basetypes.NewNumberValue(big.NewFloat(float64(1)))),
			errMsg:   "",
		},
		{
			testName: "float64",
			obj:      float64(1.1),
			expected: basetypes.NewDynamicValue(basetypes.NewNumberValue(big.NewFloat(1.1))),
			errMsg:   "",
		},
		{
			testName: "bool",
			obj:      true,
			expected: basetypes.NewDynamicValue(basetypes.NewBoolValue(true)),
			errMsg:   "",
		},
		{
			testName: "string",
			obj:      "foo",
			expected: basetypes.NewDynamicValue(basetypes.NewStringValue("foo")),
			errMsg:   "",
		},
		{
			testName: "object_empty",
			obj:      map[string]any{},
			expected: basetypes.NewDynamicValue(emptyObject),
		},
		{
			testName: "object_simple",
			obj:      map[string]any{"foo": "bar"},
			expected: basetypes.NewDynamicValue(simpleObject),
			errMsg:   "",
		},
		{
			testName: "array_empty",
			obj:      []any{},
			expected: basetypes.NewDynamicValue(emptyTuple),
			errMsg:   "",
		},
		{
			testName: "array_strings",
			obj:      []any{"foo", "bar"},
			expected: basetypes.NewDynamicValue(stringTuple),
			errMsg:   "",
		},
		//[]any{map[string]any{"foo": "bar"}, map[string]any{"foo": "baz"}}
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			dyn, diags := Decode(ctx, d.obj)

			if !dyn.Equal(d.expected) {
				// if !reflect.DeepEqual(dyn.UnderlyingValue(), d.expected.UnderlyingValue()) {
				t.Errorf("expected %v, got %v", d.expected, dyn)
			}

			var errMsg string
			if diags.HasError() {
				for i, diag := range diags.Errors() {
					if i == 0 {
						errMsg = diag.Summary()
						continue
					}
					errMsg = fmt.Sprintf("%s: %s", errMsg, diag.Summary())
				}
			}

			if errMsg != d.errMsg {
				t.Errorf("expected error message %s, got %s", d.errMsg, errMsg)
			}
		})
	}
}
