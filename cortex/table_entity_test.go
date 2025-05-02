package cortex

import (
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"gopkg.in/yaml.v3"
)

func prepareEntityResponse(t *testing.T, entities []CortexEntityElement, page, totalPages, total int) []byte {
	t.Helper()
	response := CortexEntityResponse{
		Entities:   entities,
		Page:       page,
		TotalPages: totalPages,
		Total:      total,
	}
	responseBytes, err := yaml.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	return responseBytes
}

func TestListEntitiesSinglePage(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	responseBytes := prepareEntityResponse(t, []CortexEntityElement{{Name: "entity1"}}, 0, 1, 1)

	ctx, server, client := setupTestServerAndClient(t,
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/catalog"),
			gh.VerifyHeaderKV("Authorization", "Bearer fake_api_key"),
			gh.RespondWith(http.StatusOK, responseBytes, nil),
		),
	)
	defer server.Close()

	writer := NewSliceWriter[CortexEntityElement](100)

	err := listEntities(ctx, client, writer, "false", "")
	g.Expect(err).To(BeNil())

	g.Expect(writer.Items).To(HaveLen(1))
	g.Expect(writer.Items[0].Name).To(Equal("entity1"))
}

func TestListEntitiesMultiPage(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	respPage0Bytes := prepareEntityResponse(t, []CortexEntityElement{
		{Name: "entity1"},
		{Name: "entity2"},
	}, 0, 2, 3)

	respPage1Bytes := prepareEntityResponse(t, []CortexEntityElement{
		{Name: "entity3"},
	}, 1, 2, 3)

	ctx, server, client := setupTestServerAndClient(t,
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/catalog"),
			gh.VerifyHeaderKV("Authorization", "Bearer fake_api_key"),
			gh.RespondWith(http.StatusOK, respPage0Bytes, nil),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/catalog"),
			gh.VerifyHeaderKV("Authorization", "Bearer fake_api_key"),
			gh.RespondWith(http.StatusOK, respPage1Bytes, nil),
		),
	)
	defer server.Close()

	writer := NewSliceWriter[CortexEntityElement](100)

	err := listEntities(ctx, client, writer, "false", "")
	g.Expect(err).To(BeNil())

	g.Expect(writer.Items).To(HaveLen(3))
	g.Expect(writer.Items[0].Name).To(Equal("entity1"))
	g.Expect(writer.Items[1].Name).To(Equal("entity2"))
	g.Expect(writer.Items[2].Name).To(Equal("entity3"))
}

func TestListEntitiesError(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	ctx, server, client := setupTestServerAndClient(t,
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/catalog"),
			gh.VerifyHeaderKV("Authorization", "Bearer fake_api_key"),
			gh.RespondWith(http.StatusInternalServerError, "{\"details\": \"fake error on page 0\"}", nil),
		),
	)
	defer server.Close()

	writer := NewSliceWriter[CortexEntityElement](100)

	err := listEntities(ctx, client, writer, "false", "")
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("error from cortex API 500 Internal Server Error: {\"details\": \"fake error on page 0\"}"))
}

func TestTableCortexEntity(t *testing.T) {
	g := NewWithT(t)
	table := tableCortexEntity()

	// Check basic table properties.
	g.Expect(table).ToNot(BeNil())
	g.Expect(table.Name).To(Equal("cortex_entity"))
	g.Expect(table.Description).To(Equal("Cortex list entities api."))

	// Check list configuration.
	g.Expect(table.List).ToNot(BeNil())
	g.Expect(table.List.Hydrate).ToNot(BeNil())
	g.Expect(table.List.KeyColumns).To(HaveLen(2))
	g.Expect(table.List.KeyColumns[0].Name).To(Equal("archived"))
	g.Expect(table.List.KeyColumns[0].Require).To(Equal(plugin.Optional))
	g.Expect(table.List.KeyColumns[1].Name).To(Equal("type"))
	g.Expect(table.List.KeyColumns[1].Require).To(Equal(plugin.Optional))

	// Define expected columns.
	expectedColumns := []struct {
		Name string
		Type proto.ColumnType
	}{
		{"name", proto.ColumnType_STRING},
		{"tag", proto.ColumnType_STRING},
		{"description", proto.ColumnType_STRING},
		{"type", proto.ColumnType_STRING},
		{"parents", proto.ColumnType_JSON},
		{"groups", proto.ColumnType_JSON},
		{"metadata", proto.ColumnType_JSON},
		{"last_updated", proto.ColumnType_TIMESTAMP},
		{"links", proto.ColumnType_JSON},
		{"archived", proto.ColumnType_BOOL},
		{"repository", proto.ColumnType_STRING},
		{"slack_channels", proto.ColumnType_JSON},
		{"owner_teams", proto.ColumnType_JSON},
		{"owner_individuals", proto.ColumnType_JSON},
	}

	// Check that the table has the expected columns.
	g.Expect(table.Columns).To(HaveLen(len(expectedColumns)))
	for i, exp := range expectedColumns {
		g.Expect(table.Columns[i].Name).To(Equal(exp.Name))
		g.Expect(table.Columns[i].Type).To(Equal(exp.Type))
	}
}
