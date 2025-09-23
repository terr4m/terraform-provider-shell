package tfdynamic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EncodeDynamic encodes a Terraform dynamic value.
func EncodeDynamic(ctx context.Context, d types.Dynamic) (any, error) {
	if d.IsUnknown() {
		return nil, fmt.Errorf("underlying value is unknown")
	}

	if d.IsNull() || d.IsUnderlyingValueNull() {
		return nil, nil
	}

	return encodeScalar(d.UnderlyingValue())
}

// encodeScalar encodes a scalar attribute value into an any value.
func encodeScalar(v attr.Value) (any, error) {
	switch val := v.(type) {
	case types.Bool:
		return val.ValueBool(), nil
	case types.String:
		return val.ValueString(), nil
	case types.Int64:
		return val.ValueInt64(), nil
	case types.Float64:
		return val.ValueFloat64(), nil
	case types.Number:
		f, _ := val.ValueBigFloat().Float64()
		return f, nil
	case types.List:
		return encodeSequence(val.Elements())
	case types.Set:
		return encodeSequence(val.Elements())
	case types.Tuple:
		return encodeSequence(val.Elements())
	case types.Map:
		return encodeMapping(val.Elements())
	case types.Object:
		return encodeMapping(val.Attributes())
	default:
		return nil, fmt.Errorf("unexpected type: %T", val)
	}
}

// encodeMapping encodes a map of attributes to a map of any.
func encodeMapping(m map[string]attr.Value) (map[string]any, error) {
	result := make(map[string]any, len(m))
	for k, v := range m {
		a, err := encodeScalar(v)
		if err != nil {
			return nil, err
		}

		result[k] = a
	}
	return result, nil
}

// encodeSequence encodes a list of attributes to a list of any.
func encodeSequence(s []attr.Value) ([]any, error) {
	result := make([]any, len(s))
	for i, v := range s {
		a, err := encodeScalar(v)
		if err != nil {
			return nil, err
		}

		result[i] = a
	}
	return result, nil
}
