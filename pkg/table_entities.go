package pkg

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"gopkg.in/yaml.v2"
)

type ScalarOrMap struct {
	Scalar interface{}
	Map    map[string]interface{}
}

func (s *ScalarOrMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&s.Map); err == nil {
		return nil
	}
	if err := unmarshal(&s.Scalar); err != nil {
		return err
	}
	return nil
}

func (s *ScalarOrMap) Value() interface{} {
	if s.Scalar != nil {
		return s.Scalar
	}
	return s.Map
}

type CortexEntityResponse struct {
	Entities   []CortexEntityElement `yaml:"entities"`
	Page       int                   `yaml:"page"`
	TotalPages int                   `yaml:"totalPages"`
	Total      int                   `yaml:"total"`
}

type CortexEntityElement struct {
	Name        string                        `yaml:"name"`
	Tag         string                        `yaml:"tag"`
	Description string                        `yaml:"description"`
	Type        string                        `yaml:"type"`
	Hierarchy   CortexEntityElementHierarchy  `yaml:"hierarchy"`
	Groups      []string                      `yaml:"groups"`
	Metadata    []CortexEntityElementMetadata `yaml:"metadata"`
	LastUpdated string                        `yaml:"lastUpdated"`
	Links       []CortexLink                  `yaml:"links"`
	Archived    bool                          `yaml:"isArchived"`
	Git         CortexGithub                  `yaml:"git"`
	Slack       []CortexSlackChannel          `yaml:"slackChannels"`
	// Members TODO
}

type CortexEntityElementHierarchy struct {
	Parents []CortexTag `yaml:"parents"`
}

type CortexEntityElementMetadata struct {
	Key   string      `yaml:"key"`
	Value ScalarOrMap `yaml:"value"`
}

func tableCortexEntities() *plugin.Table {
	return &plugin.Table{
		Name:        "cortex_entities",
		Description: "Cortex list entities api.",
		List: &plugin.ListConfig{
			Hydrate: listEntities,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Pretty name of the entity."},
			{Name: "tag", Type: proto.ColumnType_STRING, Description: "The x-cortex-tag of the entity."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Entity Type."},
			{Name: "parents", Type: proto.ColumnType_JSON, Description: "Parents of the entity.", Transform: FromStructSlice[CortexTag]("Hierarchy.Parents", "Tag")},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "Groups, kind of like tags."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Raw custom metadata", Transform: transform.FromField("Metadata").Transform(TagArrayToMap)},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Last updated time."},
			{Name: "links", Type: proto.ColumnType_JSON, Description: "List of links", Transform: FromStructSlice[CortexLink]("Links", "Url")},
			{Name: "archived", Type: proto.ColumnType_BOOL, Description: "Is archived."},
			{Name: "repository", Type: proto.ColumnType_STRING, Description: "Git repo full name", Transform: transform.FromField("Git.Repository")},
			{Name: "slack_channels", Type: proto.ColumnType_JSON, Description: "List of string slack channels"},
		},
	}
}

func listEntities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	var response CortexEntityResponse

	err := req.C().
		SetJsonUnmarshal(yaml.Unmarshal).
		Get("https://api.getcortexapp.com/api/v1/catalog").
		SetBearerAuthToken(*config.ApiKey).
		SetQueryParam("yaml", "false").
		SetQueryParam("includeArchived", "false").
		SetQueryParam("includeMetadata", "true").
		SetQueryParam("includeLinks", "true").
		SetQueryParam("includeSlackChannels", "true").
		SetQueryParam("includeOwners", "true").
		Do(ctx).
		Into(&response)
	if err != nil {
		return nil, err
	}
	for _, result := range response.Entities {
		// send the item to steampipe
		d.StreamListItem(ctx, result)
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
