package cortex

import (
	"context"

	"github.com/imroc/req/v3"
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

	// Enriched data
	Children []string `yaml:"-"`
	Parents  []string `yaml:"-"`
}

type CortexTeamIDPGroup struct {
	Group    string             `yaml:"group"`
	Provider string             `yaml:"provider"`
	Members  []CortexTeamMember `yaml:"members"`
}

type CortexRelationshipsResponse struct {
	Edges []CortexRelationshipsEdge `yaml:"edges"`
}

type CortexRelationshipsEdge struct {
	Child  string `yaml:"childTeamTag"`
	Parent string `yaml:"parentTeamTag"`
}

type Relationships struct {
	Children []string
	Parents  []string
}

func tableCortexTeam() *plugin.Table {
	return &plugin.Table{
		Name:        "cortex_team",
		Description: "Cortex list teams api.",
		List: &plugin.ListConfig{
			Hydrate: listTeams,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The pretty name of the team.", Transform: transform.FromField("Metadata.name")},
			{Name: "tag", Type: proto.ColumnType_STRING, Description: "The teamTag of the team."},
			{Name: "parents", Type: proto.ColumnType_JSON, Description: "Parents of the entity."},
			{Name: "children", Type: proto.ColumnType_JSON, Description: "Parents of the entity."},
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

	relationships, err := getTeamRelationships(ctx, client)
	if err != nil {
		logger.Warn("listTeams", "Error", err)
	}

	var response CortexTeamResponse
	err = client.
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
		// enrich the data
		relationships, ok := relationships[result.Tag]
		logger.Debug("listTeams", "relationships", relationships, "ok", ok)
		if ok {
			result.Children = relationships.Children
			result.Parents = relationships.Parents
		}
		// send the item to steampipe
		d.StreamListItem(ctx, result)
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func getTeamRelationships(ctx context.Context, client *req.Client) (map[string]Relationships, error) {
	logger := plugin.Logger(ctx)
	relationships := make(map[string]Relationships)

	var response CortexRelationshipsResponse
	err := client.
		Get("/api/v1/teams/relationships").
		Do(ctx).
		Into(&response)
	if err != nil {
		logger.Error("getTeamRelationships", "Error", err)
		return nil, err
	}
	logger.Info("getTeamRelationships", "results", len(response.Edges))
	for _, edges := range response.Edges {
		child := relationships[edges.Child]
		parent := relationships[edges.Parent]
		child.Parents = append(child.Parents, edges.Parent)
		parent.Children = append(parent.Children, edges.Parent)
		relationships[edges.Child] = child
		relationships[edges.Parent] = parent
	}
	return relationships, nil
}
