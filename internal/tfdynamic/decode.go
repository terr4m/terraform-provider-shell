package tfdynamic

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Decode decodes an object into a Terraform attribute value.
func Decode(ctx context.Context, obj any) (types.Dynamic, diag.Diagnostics) {
	if obj == nil {
		return types.DynamicNull(), nil
	}

	val, diags := decodeScalar(ctx, obj, path.Empty())
	if diags.HasError() {
		return types.Dynamic{}, diags
	}

	return types.DynamicValue(val), diags
}

// decodeScalar decodes a scalar value into a Terraform attribute value.
func decodeScalar(ctx context.Context, a any, thisPath path.Path) (attr.Value, diag.Diagnostics) {
	switch v := a.(type) {
	case nil:
		return types.DynamicNull(), nil
	case int64:
		return types.NumberValue(big.NewFloat(float64(v))), nil
	case float64:
		return types.NumberValue(big.NewFloat(v)), nil
	case bool:
		return types.BoolValue(v), nil
	case string:
		return types.StringValue(v), nil
	case []any:
		return decodeSequence(ctx, v, thisPath)
	case map[string]any:
		return decodeMapping(ctx, v, thisPath)
	default:
		diagnostics := diag.Diagnostics{}
		diagnostics.AddError("Unexpected type.", fmt.Sprintf("unexpected type: %T for value %#v", v, v))
		return nil, diagnostics
	}
}

// decodeMapping decodes a mapping value into a Terraform attribute value.
func decodeMapping(ctx context.Context, m map[string]any, thisPath path.Path) (attr.Value, diag.Diagnostics) {
	l := len(m)
	vm := make(map[string]attr.Value, l)
	tm := make(map[string]attr.Type, l)

	for k, v := range m {
		p := thisPath.AtName(k)
		vv, diags := decodeScalar(ctx, v, p)
		if diags.HasError() {
			return nil, diags
		}

		vm[k] = vv
		tm[k] = vv.Type(ctx)
	}

	return types.ObjectValue(tm, vm)
}

// decodeSequence decodes a sequence value into a Terraform attribute value.
func decodeSequence(ctx context.Context, s []any, thisPath path.Path) (attr.Value, diag.Diagnostics) {
	l := len(s)
	vl := make([]attr.Value, l)
	tl := make([]attr.Type, l)

	for i, v := range s {
		p := thisPath.AtListIndex(i)
		vv, err := decodeScalar(ctx, v, p)
		if err != nil {
			return nil, err
		}
		vl[i] = vv
		tl[i] = vv.Type(ctx)
	}

	return types.TupleValue(tl, vl)
}
