package cortex

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/imroc/req/v3"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"gopkg.in/yaml.v3"
)

// The Cortex API is always requested with `yaml=false`, i.e. it returns JSON, but
// responses are parsed with gopkg.in/yaml.v3 so the structs can share a single set
// of `yaml:` tags for both this and any genuinely-YAML endpoints. That works for
// almost all JSON, since JSON is a subset of YAML - except for one case: JSON
// encodes any character outside the Basic Multilingual Plane (e.g. an emoji) as a
// *pair* of `\uXXXX` UTF-16 surrogate escapes, whereas yaml.v3's `\uXXXX` escape
// represents a raw Unicode scalar and always rejects the surrogate range
// (U+D800-U+DFFF) outright, even when the pair is well-formed. That single invalid
// escape aborts parsing of the entire response it's part of, so a single Cortex
// entity with an emoji in e.g. its description can break listing every entity.
//
// This resolves surrogate pairs into their real UTF-8 encoding before handing the
// bytes to yaml.Unmarshal, so it never sees a raw surrogate escape.
var surrogatePairPattern = regexp.MustCompile(`\\u([dD][89abAB][0-9a-fA-F]{2})\\u([dD][c-fC-F][0-9a-fA-F]{2})`)

func unescapeSurrogatePairs(data []byte) []byte {
	return surrogatePairPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		groups := surrogatePairPattern.FindSubmatch(match)

		high, err := strconv.ParseUint(string(groups[1]), 16, 32)
		if err != nil {
			return match
		}
		low, err := strconv.ParseUint(string(groups[2]), 16, 32)
		if err != nil {
			return match
		}

		r := (rune(high)-0xD800)<<10 + (rune(low) - 0xDC00) + 0x10000
		return []byte(string(r))
	})
}

// yamlUnmarshalJSON parses a JSON (or YAML) document using yaml.v3, after first
// resolving any UTF-16 surrogate pair escapes JSON may have produced. See
// unescapeSurrogatePairs for why that's necessary.
func yamlUnmarshalJSON(data []byte, v interface{}) error {
	return yaml.Unmarshal(unescapeSurrogatePairs(data), v)
}

// Create a req http client for the Cortex API.
// This will set the BaseURL and Auth from config, as well as common retry settings.
func CortexHTTPClient(ctx context.Context, config *SteampipeConfig) *req.Client {
	return req.C().
		SetBaseURL(*config.BaseURL).
		SetJsonUnmarshal(yamlUnmarshalJSON).
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
