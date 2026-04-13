package tfdynamic

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	emptyObject, _ := types.ObjectValue(map[string]attr.Type{}, map[string]attr.Value{})
	simpleObject, _ := types.ObjectValue(map[string]attr.Type{"foo": types.StringType}, map[string]attr.Value{"foo": types.StringValue("bar")})
	unknownObject, _ := types.ObjectValue(map[string]attr.Type{"foo": types.StringType, "baz": types.DynamicType}, map[string]attr.Value{"foo": types.StringValue("bar"), "baz": types.DynamicUnknown()})

	nestedObject, _ := types.ObjectValue(
		map[string]attr.Type{"child": types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType}}},
		map[string]attr.Value{"child": func() attr.Value {
			v, _ := types.ObjectValue(map[string]attr.Type{"key": types.StringType}, map[string]attr.Value{"key": types.StringValue("val")})
			return v
		}()},
	)

	obj1, _ := types.ObjectValue(map[string]attr.Type{"foo": types.StringType}, map[string]attr.Value{"foo": types.StringValue("bar")})
	obj2, _ := types.ObjectValue(map[string]attr.Type{"foo": types.StringType}, map[string]attr.Value{"foo": types.StringValue("baz")})
	arrayOfObjects, _ := types.TupleValue(
		[]attr.Type{obj1.Type(ctx), obj2.Type(ctx)},
		[]attr.Value{obj1, obj2},
	)

	innerTuple, _ := types.TupleValue(
		[]attr.Type{types.NumberType, types.NumberType},
		[]attr.Value{types.NumberValue(big.NewFloat(1)), types.NumberValue(big.NewFloat(2))},
	)
	objectWithArray, _ := types.ObjectValue(
		map[string]attr.Type{"nums": innerTuple.Type(ctx)},
		map[string]attr.Value{"nums": innerTuple},
	)

	emptyTuple, _ := types.TupleValue([]attr.Type{}, []attr.Value{})
	stringTuple, _ := types.TupleValue([]attr.Type{types.StringType, types.StringType}, []attr.Value{types.StringValue("foo"), types.StringValue("bar")})
	unknownTuple, _ := types.TupleValue([]attr.Type{types.StringType, types.StringType, types.DynamicType}, []attr.Value{types.StringValue("foo"), types.StringValue("bar"), types.DynamicUnknown()})

	for _, d := range []struct {
		testName string
		obj      any
		expected types.Dynamic
		errMsg   string
	}{
		{
			testName: "unexpected_type",
			obj:      1,
			expected: types.Dynamic{},
			errMsg:   "Unexpected type.",
		},
		{
			testName: "null",
			obj:      nil,
			expected: types.DynamicNull(),
			errMsg:   "",
		},
		{
			testName: "int64",
			obj:      int64(1),
			expected: types.DynamicValue(types.NumberValue(big.NewFloat(float64(1)))),
			errMsg:   "",
		},
		{
			testName: "float64",
			obj:      float64(1.1),
			expected: types.DynamicValue(types.NumberValue(big.NewFloat(1.1))),
			errMsg:   "",
		},
		{
			testName: "bool",
			obj:      true,
			expected: types.DynamicValue(types.BoolValue(true)),
			errMsg:   "",
		},
		{
			testName: "string",
			obj:      "foo",
			expected: types.DynamicValue(types.StringValue("foo")),
			errMsg:   "",
		},
		{
			testName: "string_unknown",
			obj:      "???",
			expected: types.DynamicUnknown(),
			errMsg:   "",
		},
		{
			testName: "object_empty",
			obj:      map[string]any{},
			expected: types.DynamicValue(emptyObject),
		},
		{
			testName: "object_simple",
			obj:      map[string]any{"foo": "bar"},
			expected: types.DynamicValue(simpleObject),
			errMsg:   "",
		},
		{
			testName: "object_with_unknown",
			obj:      map[string]any{"foo": "bar", "baz": "???"},
			expected: types.DynamicValue(unknownObject),
			errMsg:   "",
		},
		{
			testName: "array_empty",
			obj:      []any{},
			expected: types.DynamicValue(emptyTuple),
			errMsg:   "",
		},
		{
			testName: "array_strings",
			obj:      []any{"foo", "bar"},
			expected: types.DynamicValue(stringTuple),
			errMsg:   "",
		},
		{
			testName: "array_with_unknown",
			obj:      []any{"foo", "bar", "???"},
			expected: types.DynamicValue(unknownTuple),
			errMsg:   "",
		},
		{
			testName: "nested_object",
			obj:      map[string]any{"child": map[string]any{"key": "val"}},
			expected: types.DynamicValue(nestedObject),
			errMsg:   "",
		},
		{
			testName: "array_of_objects",
			obj:      []any{map[string]any{"foo": "bar"}, map[string]any{"foo": "baz"}},
			expected: types.DynamicValue(arrayOfObjects),
			errMsg:   "",
		},
		{
			testName: "object_with_array",
			obj:      map[string]any{"nums": []any{float64(1), float64(2)}},
			expected: types.DynamicValue(objectWithArray),
			errMsg:   "",
		},
		{
			testName: "array_unexpected_type",
			obj:      []any{1},
			expected: types.Dynamic{},
			errMsg:   "Unexpected type.",
		},
		{
			testName: "object_unexpected_type",
			obj:      map[string]any{"bad": 1},
			expected: types.Dynamic{},
			errMsg:   "Unexpected type.",
		},
	} {
		t.Run(d.testName, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()

			dyn, diags := Decode(ctx, d.obj)

			if len(d.errMsg) > 0 {
				if !diags.HasError() {
					t.Errorf("expected error message %s, got none", d.errMsg)
					return
				}

				var errMsg string
				for i, diag := range diags.Errors() {
					if i == 0 {
						errMsg = diag.Summary()
						continue
					}
					errMsg = fmt.Sprintf("%s: %s", errMsg, diag.Summary())
				}

				if errMsg != d.errMsg {
					t.Errorf("expected error message %s, got %s", d.errMsg, errMsg)
					return
				}
			}

			if !dyn.Equal(d.expected) {
				t.Errorf("expected %v, got %v", d.expected, dyn)
			}
		})
	}
}
