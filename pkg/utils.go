package pkg

import (
	"context"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Get field from the data and for each item of type T, get the nested field "child"
// always returns a string array
func FromStructSlice[T any](field string, child string) *transform.ColumnTransforms {
	return &transform.ColumnTransforms{Transforms: []*transform.TransformCall{
		{Transform: transform.FieldValue, Param: field},
		{Transform: func(ctx context.Context, td *transform.TransformData) (interface{}, error) {
			var output []string
			vals, ok := td.Value.([]T)
			if !ok {
				return nil, nil
			}
			for _, val := range vals {
				newVal, _ := helpers.GetNestedFieldValueFromInterface(val, child)
				output = append(output, newVal.(string))
			}
			return output, nil
		}},
		{Transform: transform.EnsureStringArray},
	}}
}

func TagArrayToMap(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	result := map[string]interface{}{}
	for _, value := range d.Value.([]CortexEntityElementMetadata) {
		result[value.Key] = value.Value.Value()
	}
	return result, nil
}
