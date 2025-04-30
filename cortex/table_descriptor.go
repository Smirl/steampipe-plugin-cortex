package cortex

import (
	"context"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type CortexDescriptorsResponse struct {
	Descriptors []Cortex `yaml:"descriptors"`
	Page        int      `yaml:"page"`
	TotalPages  int      `yaml:"totalPages"`
	Total       int      `yaml:"total"`
}

func tableCortexDescriptor() *plugin.Table {
	return &plugin.Table{
		Name:        "cortex_descriptor",
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
			{Name: "slos", Type: proto.ColumnType_JSON, Description: "SLOs from each integration if any", Transform: transform.FromField("SLOs")},
			{Name: "static_analysis", Type: proto.ColumnType_JSON, Description: "Static analysis", Transform: transform.FromField("StaticAnalysis")},
		},
	}
}

func listDescriptors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	config := GetConfig(d.Connection)
	client := CortexHTTPClient(ctx, config)

	var response CortexDescriptorsResponse
	var page int = 0
	for {
		logger.Debug("listDescriptors", "page", page)
		err := client.
			Get("/api/v1/catalog/descriptors").
			// Options
			SetQueryParam("yaml", "false").
			// Pagination
			SetQueryParam("pageSize", "1000").
			SetQueryParam("page", strconv.Itoa(page)).
			Do(ctx).
			Into(&response)
		if err != nil {
			logger.Error("listDescriptors", "Error", err)
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
		page++
		if page >= response.TotalPages {
			break
		}
	}
	return nil, nil
}
