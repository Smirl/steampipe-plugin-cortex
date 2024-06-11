package pkg

import (
	"context"

	"github.com/imroc/req/v3"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"gopkg.in/yaml.v2"
)

type CortexDescriptorsResponse struct {
	Descriptors []Cortex `yaml:"descriptors"`
	Page        int      `yaml:"page"`
	TotalPages  int      `yaml:"totalPages"`
	Total       int      `yaml:"total"`
}

func tableCortexDescriptors() *plugin.Table {
	return &plugin.Table{
		Name:        "cortex_descriptors",
		Description: "Cortex openapi descriptors.",
		List: &plugin.ListConfig{
			Hydrate: listDescriptors,
		},
		Columns: []*plugin.Column{
			{Name: "tag", Type: proto.ColumnType_STRING, Description: "The x-cortex-tag of the entity."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Entity Type."},
			{Name: "parents", Type: proto.ColumnType_JSON, Description: "Parent tags.", Transform: FromStructSlice[CortexTag]("Parents", "Tag")},
			{Name: "groups", Type: proto.ColumnType_JSON, Description: "Groups, kind of like tags."},
			{Name: "team", Type: proto.ColumnType_JSON, Description: "Raw team"},
			{Name: "owners", Type: proto.ColumnType_JSON, Description: "Raw owner"},
			{Name: "slack", Type: proto.ColumnType_JSON, Description: "Raw slack"},
			{Name: "links", Type: proto.ColumnType_JSON, Description: "List of links", Transform: FromStructSlice[CortexLink]("Links", "Url")},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Raw custom metadata", Transform: transform.FromField("CustomMetadata")},
			{Name: "repository", Type: proto.ColumnType_STRING, Description: "Git repo full name", Transform: transform.FromField("Git.Github.Repository")},
			{Name: "victorops", Type: proto.ColumnType_STRING, Description: "Victorops team slug", Transform: transform.FromField("Oncall.VictorOps.ID")},
			{Name: "jira", Type: proto.ColumnType_JSON, Description: "List of jira projects", Transform: transform.FromField("Issues.Jira.Projects").Transform(transform.EnsureStringArray)},
		},
	}
}

func listDescriptors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	config := GetConfig(d.Connection)
	var response CortexDescriptorsResponse
	err := req.C().
		SetJsonUnmarshal(yaml.Unmarshal).
		Get("https://api.getcortexapp.com/api/v1/catalog/descriptors").
		SetBearerAuthToken(*config.ApiKey).
		SetQueryParam("yaml", "false").
		Do(ctx).
		Into(&response)
	if err != nil {
		return nil, err
	}
	for _, result := range response.Descriptors {
		// send the item to steampipe
		d.StreamListItem(ctx, result.Info)
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}