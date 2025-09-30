package cortex

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/quals"
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
			Hydrate: listEntitiesHydrator,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "archived", Require: plugin.Optional},
				{Name: "type", Require: plugin.Optional},
				{Name: "groups", Require: plugin.Optional, Operators: []string{"=", "?", "?|"}},
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

func listEntitiesHydrator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	config := GetConfig(d.Connection)
	client := CortexHTTPClient(ctx, config)
	hydratorWriter := QueryDataWriter{d}

	// Extract parameters from QueryData
	archived := "false"
	if d.EqualsQuals["archived"] != nil && d.EqualsQuals["archived"].GetBoolValue() {
		logger.Debug("listEntitiesHydrator", "archived", d.EqualsQuals["archived"])
		archived = "true"
	}
	types := ""
	if d.EqualsQuals["type"] != nil {
		// When doing a "where in ()" steampipe does multiple separate calls to listEntities
		types = d.EqualsQuals["type"].GetStringValue()
	}
	groups := ""
	logger.Debug("listEntitiesHydrator", "quals", d.Quals)
	if d.Quals["groups"] != nil {
		groupFilters := buildGroupFilters(d.Quals["groups"].Quals)
		logger.Debug("listEntitiesHydrator", "groupFilters", groupFilters)
		if len(groupFilters) > 0 {
			groups = strings.Join(groupFilters, ",")
		}
	}
	logger.Info("listEntitiesHydrator", "archived", archived, "types", types, "groups", groups)
	return nil, listEntities(ctx, client, &hydratorWriter, archived, types, groups)
}

func listEntities(ctx context.Context, client *req.Client, writer HydratorWriter, archived string, types string, groups string) error {
	logger := plugin.Logger(ctx)

	var response CortexEntityResponse
	var page int = 0
	for {
		logger.Debug("listEntities", "page", page)
		resp := client.
			Get("/api/v1/catalog").
			// Filters
			SetQueryParam("includeArchived", archived).
			SetQueryParam("types", types).
			SetQueryParam("groups", groups).
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
			Do(ctx)

		// Check for HTTP errors
		if resp.IsErrorState() {
			logger.Error("listEntities", "Status", resp.Status, "Body", resp.String())
			return fmt.Errorf("error from cortex API %s: %s", resp.Status, resp.String())
		}

		// Unmarshal the response and check for unmarshal errors
		err := resp.Into(&response)
		if err != nil {
			logger.Error("listEntities", "page", page, "Error", err)
			return err
		}

		logger.Debug("listEntities", "totalPages", response.TotalPages, "total", response.Total)

		for _, result := range response.Entities {
			// send the item to steampipe
			writer.StreamListItem(ctx, result)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if writer.RowsRemaining(ctx) == 0 {
				logger.Debug("listEntities", "RowsRemaining", writer.RowsRemaining(ctx))
				return nil
			}
		}
		page++
		if page >= response.TotalPages {
			logger.Debug("listEntities", "page", page, "totalPages", response.TotalPages)
			break
		}
	}
	return nil
}

func buildGroupFilters(groupQuals []*quals.Qual) []string {
	var groupFilters []string
	for _, q := range groupQuals {
		switch q.Operator {
		case quals.QualOperatorJsonbExistsOne, quals.QualOperatorEqual:
			if value := q.Value.GetStringValue(); value != "" {
				groupFilters = append(groupFilters, value)
			}
		case quals.QualOperatorJsonbExistsAny:
			if listValue := q.Value.GetListValue(); listValue != nil {
				for _, v := range listValue.Values {
					if value := v.GetStringValue(); value != "" {
						groupFilters = append(groupFilters, value)
					}
				}
			}
		}
	}
	return groupFilters
}
