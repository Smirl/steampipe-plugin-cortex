package pkg

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type CortexTeamResponse struct {
	Teams []CortexTeamElement `yaml:"teams"`
}

type CortexTeamElement struct {
	Tag      string                 `yaml:"teamTag"`
	Metadata map[string]interface{} `yaml:"metadata"`
	Links    []CortexLink           `yaml:"links"`
	Archived bool                   `yaml:"isArchived"`
	Slack    []CortexSlackChannel   `yaml:"slackChannels"`
	IDPGroup CortexTeamIDPGroup     `yaml:"idpGroup"`
}

type CortexTeamIDPGroup struct {
	Group    string             `yaml:"group"`
	Provider string             `yaml:"provider"`
	Members  []CortexTeamMember `yaml:"members"`
}

func tableCortexTeam() *plugin.Table {
	return &plugin.Table{
		Name:        "cortex_team",
		Description: "Cortex list teams api.",
		List: &plugin.ListConfig{
			Hydrate: listTeams,
		},
		Columns: []*plugin.Column{
			{Name: "Name", Type: proto.ColumnType_STRING, Description: "The pretty name of the team.", Transform: transform.FromField("Metadata.name")},
			{Name: "tag", Type: proto.ColumnType_STRING, Description: "The teamTag of the team."},
			// {Name: "parents", Type: proto.ColumnType_JSON, Description: "Parents of the entity.", Transform: FromStructSlice[CortexTag]("Hierarchy.Parents", "Tag")},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Raw custom metadata"},
			{Name: "links", Type: proto.ColumnType_JSON, Description: "List of links", Transform: FromStructSlice[CortexLink]("Links", "Url")},
			{Name: "archived", Type: proto.ColumnType_BOOL, Description: "Is archived."},
			{Name: "slack_channels", Type: proto.ColumnType_JSON, Description: "List of string slack channels"},
			{Name: "members", Type: proto.ColumnType_JSON, Description: "List of members", Transform: transform.FromField("IDPGroup.Members")},
		},
	}
}

func listTeams(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	config := GetConfig(d.Connection)
	client := CortexHTTPClient(ctx, config)

	var response CortexTeamResponse
	err := client.
		EnableDumpAllToFile("/tmp/cortex").
		Get("/api/v1/teams").
		SetQueryParam("includeTeamsWithoutMembers", "true").
		Do(ctx).
		Into(&response)
	if err != nil {
		logger.Error("listTeams", "Error", err)
		return nil, err
	}
	logger.Info("listTeams", "results", len(response.Teams))
	for _, result := range response.Teams {
		// send the item to steampipe
		d.StreamListItem(ctx, result)
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
