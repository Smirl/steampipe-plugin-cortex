package cortex

import (
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"gopkg.in/yaml.v3"
)

func prepareScorecardResponse(t *testing.T, scorecard CortexScorecard) []byte {
	t.Helper()
	response := CortexScorecardResponse{
		Scorecard: scorecard,
	}
	responseBytes, err := yaml.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	return responseBytes
}

func prepareScorecardScoresResponse(t *testing.T, scores []*CortexServiceScore, page, totalPages, total int) []byte {
	t.Helper()
	response := CortexScorecardScoreResponse{
		ServiceScores: scores,
		Page:          page,
		TotalPages:    totalPages,
		Total:         total,
	}
	responseBytes, err := yaml.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	return responseBytes
}

func TestListScorecardScoresSinglePage(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	scorecard := CortexScorecard{
		Rules: []*CortexRuleInfo{
			{Identifier: "rule1", Title: "Rule 1", LevelName: "Level 1", Weight: 10},
		},
		Levels: []*CortexScorecardLevel{
			{Level: CortexLevel{Name: "Level 1", Number: 1}},
		},
	}
	scorecardResponseBytes := prepareScorecardResponse(t, scorecard)

	scores := []*CortexServiceScore{
		{
			LastEvaluated: "2025-05-02T12:00:00Z",
			Service:       &CortexEntityElement{Name: "Service 1", Tag: "service1"},
			Score: &CortexScore{
				Rules: []*CortexRuleScore{
					{Identifier: "rule1", Score: 10},
				},
			},
		},
	}
	scoresResponseBytes := prepareScorecardScoresResponse(t, scores, 0, 1, 1)

	ctx, server, client := setupTestServerAndClient(t,
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/scorecards/tag1"),
			gh.RespondWith(http.StatusOK, scorecardResponseBytes, nil),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/scorecards/tag1/scores"),
			gh.RespondWith(http.StatusOK, scoresResponseBytes, nil),
		),
	)
	defer server.Close()

	writer := NewSliceWriter[CortexScorecardScoreRow](100)

	err := listScorecardScores(ctx, client, writer, "tag1")
	g.Expect(err).To(BeNil())

	g.Expect(writer.Items).To(HaveLen(1))
	g.Expect(writer.Items[0].Service.Name).To(Equal("Service 1"))
	g.Expect(writer.Items[0].RuleScore.Identifier).To(Equal("rule1"))
	g.Expect(writer.Items[0].RuleScore.Score).To(Equal(10))
}

func TestListScorecardScoresError(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	ctx, server, client := setupTestServerAndClient(t,
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/scorecards/tag1"),
			gh.RespondWith(http.StatusInternalServerError, "{\"details\": \"fake error on scorecard\"}", nil),
		),
	)
	defer server.Close()

	writer := NewSliceWriter[CortexScorecardScoreRow](100)

	err := listScorecardScores(ctx, client, writer, "tag1")
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("error from cortex API 500 Internal Server Error: {\"details\": \"fake error on scorecard\"}"))
}
