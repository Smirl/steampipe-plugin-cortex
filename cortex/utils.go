package cortex

import (
	"context"
	"time"

	"github.com/imroc/req/v3"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"gopkg.in/yaml.v3"
)

// Create a req http client for the Cortex API.
// This will set the BaseURL and Auth from config, as well as common retry settings.
func CortexHTTPClient(ctx context.Context, config *SteampipeConfig) *req.Client {
	return req.C().
		SetBaseURL(*config.BaseURL).
		SetJsonUnmarshal(yaml.Unmarshal).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(time.Second, 5*time.Second).
		SetCommonBearerAuthToken(*config.ApiKey)
}

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

// Writer is a generic interface to stream items of any type.
type HydratorWriter interface {
	StreamListItem(ctx context.Context, items ...interface{})
	RowsRemaining(ctx context.Context) int64
}

// Production implementation that wraps a *plugin.QueryData.
type QueryDataWriter struct {
	QueryData *plugin.QueryData
}

func (h *QueryDataWriter) StreamListItem(ctx context.Context, items ...interface{}) {
	h.QueryData.StreamListItem(ctx, items...)
}

func (h *QueryDataWriter) RowsRemaining(ctx context.Context) int64 {
	return h.QueryData.RowsRemaining(ctx)
}

// Testing implementation that writes to a slice up to a fixed limit.
type SliceWriter[T any] struct {
	Limit int64
	Items []T
}

// NewSliceWriter creates a new SliceWriter with the given limit.
func NewSliceWriter[T any](limit int64) *SliceWriter[T] {
	return &SliceWriter[T]{
		Limit: limit,
		Items: make([]T, 0, limit),
	}
}

func (s *SliceWriter[T]) StreamListItem(ctx context.Context, items ...interface{}) {
	for _, item := range items {
		if typedItem, ok := item.(T); ok {
			s.Items = append(s.Items, typedItem)
		}
	}
}

func (s *SliceWriter[T]) RowsRemaining(ctx context.Context) int64 {
	return s.Limit - int64(len(s.Items))
}
