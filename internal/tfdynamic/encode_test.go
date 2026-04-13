package tfdynamic

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestEncodeDynamic(t *testing.T) {
	t.Parallel()

	simpleObject, _ := types.ObjectValue(
		map[string]attr.Type{"foo": types.StringType},
		map[string]attr.Value{"foo": types.StringValue("bar")},
	)

	stringList, _ := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("a"),
		types.StringValue("b"),
	})

	stringSet, _ := types.SetValue(types.StringType, []attr.Value{
		types.StringValue("x"),
	})

	stringTuple, _ := types.TupleValue(
		[]attr.Type{types.StringType, types.StringType},
		[]attr.Value{types.StringValue("t1"), types.StringValue("t2")},
	)

	stringMap, _ := types.MapValue(types.StringType, map[string]attr.Value{
		"k1": types.StringValue("v1"),
	})

	nestedObject, _ := types.ObjectValue(
		map[string]attr.Type{
			"name": types.StringType,
			"list": types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"name": types.StringValue("nested"),
			"list": stringList,
		},
	)

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
			testName: "null_dynamic",
			dyn:      types.DynamicNull(),
			expected: nil,
			errMsg:   "",
		},
		{
			testName: "bool_true",
			dyn:      types.DynamicValue(types.BoolValue(true)),
			expected: true,
			errMsg:   "",
		},
		{
			testName: "bool_false",
			dyn:      types.DynamicValue(types.BoolValue(false)),
			expected: false,
			errMsg:   "",
		},
		{
			testName: "string",
			dyn:      types.DynamicValue(types.StringValue("hello")),
			expected: "hello",
			errMsg:   "",
		},
		{
			testName: "int64",
			dyn:      types.DynamicValue(types.Int64Value(42)),
			expected: int64(42),
			errMsg:   "",
		},
		{
			testName: "float64",
			dyn:      types.DynamicValue(types.Float64Value(3.14)),
			expected: float64(3.14),
			errMsg:   "",
		},
		{
			testName: "number",
			dyn:      types.DynamicValue(types.NumberValue(big.NewFloat(99.5))),
			expected: float64(99.5),
			errMsg:   "",
		},
		{
			testName: "object",
			dyn:      types.DynamicValue(simpleObject),
			expected: map[string]any{"foo": "bar"},
			errMsg:   "",
		},
		{
			testName: "list",
			dyn:      types.DynamicValue(stringList),
			expected: []any{"a", "b"},
			errMsg:   "",
		},
		{
			testName: "set",
			dyn:      types.DynamicValue(stringSet),
			expected: []any{"x"},
			errMsg:   "",
		},
		{
			testName: "tuple",
			dyn:      types.DynamicValue(stringTuple),
			expected: []any{"t1", "t2"},
			errMsg:   "",
		},
		{
			testName: "map",
			dyn:      types.DynamicValue(stringMap),
			expected: map[string]any{"k1": "v1"},
			errMsg:   "",
		},
		{
			testName: "nested_object_with_list",
			dyn:      types.DynamicValue(nestedObject),
			expected: map[string]any{"name": "nested", "list": []any{"a", "b"}},
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
