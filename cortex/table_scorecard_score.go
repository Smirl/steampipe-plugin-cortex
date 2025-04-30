package cortex

import (
	"context"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Response elements for the /scorecard/{tag} endpoint
type CortexScorecardResponse struct {
	Scorecard CortexScorecard `yaml:"scorecard"`
}

type CortexScorecard struct {
	Levels []*CortexScorecardLevel `yaml:"levels"`
	Rules  []*CortexRuleInfo       `yaml:"rules"`
}

type CortexScorecardLevel struct {
	Level CortexLevel `yaml:"level"`
}

type CortexLevel struct {
	Name   string `yaml:"name"`
	Number int    `yaml:"number"`
}

type CortexRuleInfo struct {
	Description   string `yaml:"description"`
	EffectiveFrom string `yaml:"effectiveFrom"`
	Identifier    string `yaml:"identifier"`
	LevelName     string `yaml:"levelName"`
	Title         string `yaml:"title"`
	Weight        int    `yaml:"weight"`

	// Not in the API response, but used to enrich the data
	LevelNumber int `yaml:"-"`
}

// Response elements for the /scorecards/{tag}/scores endpoint
type CortexScorecardScoreResponse struct {
	ScorecardName string                `yaml:"scorecardName"`
	ScorecardTag  string                `yaml:"scorecardTag"`
	ServiceScores []*CortexServiceScore `yaml:"serviceScores"`
	Page          int                   `yaml:"page"`
	TotalPages    int                   `yaml:"totalPages"`
	Total         int                   `yaml:"total"`
}

type CortexServiceScore struct {
	LastEvaluated string               `yaml:"lastEvaluated"`
	Service       *CortexEntityElement `yaml:"service"`
	Score         *CortexScore         `yaml:"score"`
}

type CortexScore struct {
	Rules []*CortexRuleScore `yaml:"rules"`
}

type CortexRuleScore struct {
	Expression string `yaml:"expression"`
	Identifier string `yaml:"identifier"`
	Score      int    `yaml:"score"`
}

// Used to represent the data we want to return in the table
type CortexScorecardScoreRow struct {
	ScorecardName string
	ScorecardTag  string
	LastEvaluated string
	Service       *CortexEntityElement
	RuleScore     *CortexRuleScore
	RuleInfo      *CortexRuleInfo
}

func (r *CortexScorecardScoreRow) IsRulePass() bool {
	return r.RuleScore.Score == r.RuleInfo.Weight
}

func tableCortexScorecardScore() *plugin.Table {
	return &plugin.Table{
		Name:        "cortex_scorecard_score",
		Description: "Cortex scorecard score api.",
		List: &plugin.ListConfig{
			Hydrate: listScorecardScores,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "scorecard_tag", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			{Name: "scorecard_tag", Type: proto.ColumnType_STRING, Description: "Scorecard tag."},
			{Name: "scorecard_name", Type: proto.ColumnType_STRING, Description: "Scorecard name."},
			{Name: "service_tag", Type: proto.ColumnType_STRING, Description: "Service type.", Transform: transform.FromField("Service.Tag")},
			{Name: "service_name", Type: proto.ColumnType_STRING, Description: "Service name.", Transform: transform.FromField("Service.Name")},
			{Name: "service_groups", Type: proto.ColumnType_JSON, Description: "Service groups.", Transform: transform.FromField("Service.Groups")},
			{Name: "last_evaluated", Type: proto.ColumnType_STRING, Description: "Last evaluated."},
			{Name: "rule_identifier", Type: proto.ColumnType_STRING, Description: "Rule identifier.", Transform: transform.FromField("RuleScore.Identifier")},
			{Name: "rule_title", Type: proto.ColumnType_STRING, Description: "Rule title.", Transform: transform.FromField("RuleInfo.Title")},
			{Name: "rule_description", Type: proto.ColumnType_STRING, Description: "Rule description.", Transform: transform.FromField("RuleInfo.Description")},
			{Name: "rule_expression", Type: proto.ColumnType_STRING, Description: "Rule expression.", Transform: transform.FromField("RuleScore.Expression")},
			{Name: "rule_effective_from", Type: proto.ColumnType_STRING, Description: "Rule effective from.", Transform: transform.FromField("RuleInfo.EffectiveFrom")},
			{Name: "rule_level_name", Type: proto.ColumnType_STRING, Description: "Rule level name.", Transform: transform.FromField("RuleInfo.LevelName")},
			{Name: "rule_level_number", Type: proto.ColumnType_INT, Description: "Rule level number.", Transform: transform.FromField("RuleInfo.LevelNumber")},
			{Name: "rule_weight", Type: proto.ColumnType_INT, Description: "Rule weight.", Transform: transform.FromField("RuleInfo.Weight")},
			{Name: "rule_score", Type: proto.ColumnType_INT, Description: "Rule score.", Transform: transform.FromField("RuleScore.Score")},
			{Name: "rule_pass", Type: proto.ColumnType_BOOL, Description: "Rule pass.", Transform: transform.FromP(transform.MethodValue, "IsRulePass")},
		},
	}
}

func listScorecardScores(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	config := GetConfig(d.Connection)
	client := CortexHTTPClient(ctx, config)

	var scorecardTag string = d.EqualsQuals["scorecard_tag"].GetStringValue()

	// Get information about the scorecard to enrich the data
	var scorecardResponse CortexScorecardResponse
	err := client.
		Get("/api/v1/scorecards/{tag}").
		SetPathParam("tag", scorecardTag).
		Do(ctx).
		Into(&scorecardResponse)
	if err != nil {
		logger.Error("listScorecardScores", "Error", err)
		return nil, err
	}
	// Make a map of rule identifier to CortexRule
	rules := make(map[string]*CortexRuleInfo)
	for _, rule := range scorecardResponse.Scorecard.Rules {
		rules[rule.Identifier] = rule
		// add level number to the rule
		for _, level := range scorecardResponse.Scorecard.Levels {
			if level.Level.Name == rule.LevelName {
				rule.LevelNumber = level.Level.Number
			}
		}
	}

	// Get the scores for the scorecard
	var response CortexScorecardScoreResponse
	var page int = 0
	for {
		err := client.
			Get("/api/v1/scorecards/{tag}/scores").
			SetPathParam("tag", scorecardTag).
			// Pagination
			SetQueryParam("pageSize", "1000").
			SetQueryParam("page", strconv.Itoa(page)).
			Do(ctx).
			Into(&response)
		if err != nil {
			logger.Error("listScorecardScores", "Error", err)
			return nil, err
		}
		for _, result := range response.ServiceScores {
			for _, ruleScore := range result.Score.Rules {
				// Get the rule info
				ruleInfo, ok := rules[ruleScore.Identifier]
				if !ok {
					continue
				}
				row := CortexScorecardScoreRow{
					ScorecardName: response.ScorecardName,
					ScorecardTag:  response.ScorecardTag,
					LastEvaluated: result.LastEvaluated,
					Service:       result.Service,
					RuleScore:     ruleScore,
					RuleInfo:      ruleInfo,
				}
				// send the item to steampipe
				d.StreamListItem(ctx, row)
				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.RowsRemaining(ctx) == 0 {
					return nil, nil
				}
			}
		}
		page++
		if page >= response.TotalPages {
			break
		}
	}
	return nil, nil
}
