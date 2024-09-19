package pkg

import (
	"context"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
	Owners      CortexEntityOwners            `yaml:"owners"`
}

type CortexEntityElementHierarchy struct {
	Parents []CortexTag `yaml:"parents"`
}

type CortexEntityElementMetadata struct {
	Key   string      `yaml:"key"`
	Value ScalarOrMap `yaml:"value"`
}

type CortexEntityOwners struct {
	Teams       []CortexEntityOwnersTeam       `yaml:"teams"`
	Individuals []CortexEntityOwnersIndividual `yaml:"individuals"`
}

type CortexEntityOwnersTeam struct {
	Tag string `yaml:"tag"`
}

type CortexEntityOwnersIndividual struct {
	Email string `yaml:"email"`
}

func tableCortexEntity() *plugin.Table {
	return &plugin.Table{
		Name:        "cortex_entity",
		Description: "Cortex list entities api.",
		List: &plugin.ListConfig{
			Hydrate: listEntities,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "archived", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
			},
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
			{Name: "owner_teams", Type: proto.ColumnType_JSON, Description: "List of owning team tags", Transform: FromStructSlice[CortexEntityOwnersTeam]("Owners.Teams", "Tag")},
			{Name: "owner_individuals", Type: proto.ColumnType_JSON, Description: "List of owning individuals emails", Transform: FromStructSlice[CortexEntityOwnersIndividual]("Owners.Individuals", "Email")},
		},
	}
}

func listEntities(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	config := GetConfig(d.Connection)
	client := CortexHTTPClient(ctx, config)

	// Archive filter
	var archived string = "false"
	if d.EqualsQuals["archived"] != nil && d.EqualsQuals["archived"].GetBoolValue() {
		archived = "true"
	}
	logger.Info("listEntities", "archived", archived)
	// Type filter
	// When doing a "where in ()" steampipe does multiple separate calls to listEntities
	var types string
	if d.EqualsQuals["type"] != nil {
		types = d.EqualsQuals["type"].GetStringValue()
	}
	logger.Info("listEntities", "types", types)

	var response CortexEntityResponse
	var page int = 0
	for {
		logger.Debug("listEntities", "page", page)
		err := client.
			Get("/api/v1/catalog").
			// Filters
			SetQueryParam("includeArchived", archived).
			SetQueryParam("types", types).
			// Options
			SetQueryParam("yaml", "false").
			SetQueryParam("includeMetadata", "true").
			SetQueryParam("includeLinks", "true").
			SetQueryParam("includeSlackChannels", "true").
			SetQueryParam("includeOwners", "true").
			SetQueryParam("includeHierarchyFields", "true").
			// Pagination
			SetQueryParam("pageSize", "1000").
			SetQueryParam("page", strconv.Itoa(page)).
			Do(ctx).
			Into(&response)
		if err != nil {
			logger.Error("listEntities", "Error", err)
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
		page++
		if page == response.TotalPages {
			break
		}
	}
	return nil, nil
}
