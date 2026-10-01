package cortex

import (
	"context"
	"net/http"
	"testing"
	_ "unsafe"

	"github.com/hashicorp/go-hclog"
	"github.com/imroc/req/v3"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"gopkg.in/yaml.v3"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/context_key"
)

func setupTestServerAndClient(t *testing.T, handlers ...http.HandlerFunc) (context.Context, *ghttp.Server, *req.Client) {
	t.Helper()

	// Create a fake upstream server and register handlers.
	server := ghttp.NewServer()
	for _, handler := range handlers {
		server.AppendHandlers(handler)
	}

	// Create a context with a logger.
	ctx := context.WithValue(context.Background(), context_key.Logger, hclog.NewNullLogger())

	// Create a testing client.
	config := NewSteampipeConfig("fake_api_key", server.URL())
	client := CortexHTTPClient(ctx, config)

	return ctx, server, client
}

func prepareDescriptorResponse(t *testing.T, descriptors []Cortex, page, totalPages, total int) []byte {
	t.Helper()
	response := CortexDescriptorsResponse{
		Descriptors: descriptors,
		Page:        page,
		TotalPages:  totalPages,
		Total:       total,
	}
	responseBytes, err := yaml.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	return responseBytes
}

func TestTableCortexDescriptor(t *testing.T) {
	g := NewWithT(t)
	table := tableCortexDescriptor()

	// Check basic table properties.
	g.Expect(table).ToNot(BeNil())
	g.Expect(table.Name).To(Equal("cortex_descriptor"))
	g.Expect(table.Description).To(Equal("Cortex openapi descriptors."))

	// Check list configuration.
	g.Expect(table.List).ToNot(BeNil())
	g.Expect(table.List.Hydrate).ToNot(BeNil())

	// Define expected columns.
	expectedColumns := []struct {
		Name string
		Type proto.ColumnType
	}{
		{"tag", proto.ColumnType_STRING},
		{"title", proto.ColumnType_STRING},
		{"description", proto.ColumnType_STRING},
		{"type", proto.ColumnType_STRING},
		{"parents", proto.ColumnType_JSON},
		{"groups", proto.ColumnType_JSON},
		{"team", proto.ColumnType_JSON},
		{"owners", proto.ColumnType_JSON},
		{"slack", proto.ColumnType_JSON},
		{"links", proto.ColumnType_JSON},
		{"metadata", proto.ColumnType_JSON},
		{"repository", proto.ColumnType_STRING},
		{"victorops", proto.ColumnType_STRING},
		{"jira", proto.ColumnType_JSON},
		{"slos", proto.ColumnType_JSON},
		{"static_analysis", proto.ColumnType_JSON},
	}

	// Check that the table has the expected columns.
	g.Expect(table.Columns).To(HaveLen(len(expectedColumns)))
	for i, exp := range expectedColumns {
		g.Expect(table.Columns[i].Name).To(Equal(exp.Name))
		g.Expect(table.Columns[i].Type).To(Equal(exp.Type))
	}
}

// --- Tests for listDescriptors ---
func TestListDescriptorsSinglePage(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	// Create the expected response using the helper function.
	responseBytes := prepareDescriptorResponse(t, []Cortex{{Info: CortexInfo{Tag: "tag1"}}}, 0, 1, 1)

	// Create a fake upstream server and client using the helper function.
	ctx, server, client := setupTestServerAndClient(t,
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/catalog/descriptors"),
			gh.VerifyHeaderKV("Authorization", "Bearer fake_api_key"),
			gh.RespondWith(http.StatusOK, responseBytes, nil),
		),
	)
	defer server.Close()

	// Create a HydratorWriter that will capture the streamed items
	writer := NewSliceWriter[CortexInfo](100)

	// h is unused so we pass nil.
	err := listDescriptors(ctx, client, writer)
	g.Expect(err).To(BeNil())

	g.Expect(writer.Items).To(HaveLen(1))
	g.Expect(writer.Items[0].Tag).To(Equal("tag1"))
}

func TestListDescriptorsMultiPage(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	// Prepare response for page 0: return 2 descriptors, with total 3 descriptors over 2 pages.
	respPage0Bytes := prepareDescriptorResponse(t, []Cortex{
		{Info: CortexInfo{Tag: "tag1"}},
		{Info: CortexInfo{Tag: "tag2"}},
	}, 0, 2, 3)

	// Prepare response for page 1: return 1 descriptor.
	respPage1Bytes := prepareDescriptorResponse(t, []Cortex{
		{Info: CortexInfo{Tag: "tag3"}},
	}, 1, 2, 3)

	// Create a fake upstream server and client using the helper function.
	ctx, server, client := setupTestServerAndClient(t,
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/catalog/descriptors"),
			gh.VerifyHeaderKV("Authorization", "Bearer fake_api_key"),
			gh.RespondWith(http.StatusOK, respPage0Bytes, nil),
		),
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/catalog/descriptors"),
			gh.VerifyHeaderKV("Authorization", "Bearer fake_api_key"),
			gh.RespondWith(http.StatusOK, respPage1Bytes, nil),
		),
	)
	defer server.Close()

	// Create a writer to capture the descriptors.
	writer := NewSliceWriter[CortexInfo](100)

	// Execute the listing of descriptors.
	err := listDescriptors(ctx, client, writer)
	g.Expect(err).To(BeNil())

	// Validate that all three descriptors were streamed.
	g.Expect(writer.Items).To(HaveLen(3))
	g.Expect(writer.Items[0].Tag).To(Equal("tag1"))
	g.Expect(writer.Items[1].Tag).To(Equal("tag2"))
	g.Expect(writer.Items[2].Tag).To(Equal("tag3"))
}

func TestListDescriptorsError(t *testing.T) {
	g := NewWithT(t)
	gh := ghttp.NewGHTTPWithGomega(g)

	// Create a fake upstream server and client using the helper function.
	ctx, server, client := setupTestServerAndClient(t,
		ghttp.CombineHandlers(
			gh.VerifyRequest("GET", "/api/v1/catalog/descriptors"),
			gh.VerifyHeaderKV("Authorization", "Bearer fake_api_key"),
			gh.RespondWith(http.StatusInternalServerError, "{\"details\": \"fake error on page 0\"}", nil),
		),
	)
	defer server.Close()

	// Create a writer; its contents are not used because we expect an error.
	writer := NewSliceWriter[CortexInfo](100)

	// Execute the listing of descriptors and expect an error.
	err := listDescriptors(ctx, client, writer)
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("error from cortex API 500 Internal Server Error: {\"details\": \"fake error on page 0\"}"))
}
