package cortex

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type Package struct {
	DateCreated string `yaml:"dateCreated"`
	Id          string `yaml:"id"`
	Name        string `yaml:"name"`
	PackageType string `yaml:"packageType"`
	Version     string `yaml:"version"`
}

type CortexPackageRow struct {
	PackageTag  string
	PackageType string
	Name        string
	Version     string
}

func tableCortexPackage() *plugin.Table {
	return &plugin.Table{
		Name:        "cortex_packages",
		Description: "Cortex list package api.",
		List: &plugin.ListConfig{
			Hydrate: listPackagesHydrator,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "package_tag", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			{Name: "package_tag", Type: proto.ColumnType_STRING, Description: "Package Tag."},
			{Name: "package_type", Type: proto.ColumnType_STRING, Description: "Package Type."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Package Name."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Version."},
		},
	}
}

func listPackagesHydrator(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	config := GetConfig(d.Connection)
	client := CortexHTTPClient(ctx, config)
	hydratorWriter := QueryDataWriter{d}
	packageTag := d.EqualsQuals["package_tag"].GetStringValue()
	logger.Info("listPackageHydrator", "packageTag", packageTag)
	return nil, listPackages(ctx, client, &hydratorWriter, packageTag)
}

func listPackages(ctx context.Context, client *req.Client, writer HydratorWriter, packageTag string) error {
	logger := plugin.Logger(ctx)
	var response []Package
	logger.Debug("listPackages")
	resp := client.
		Get("/api/v1/catalog/{tag}/packages").
		SetPathParam("tag", packageTag).
		// Options
		SetQueryParam("yaml", "false").
		Do(ctx)

	// Check for HTTP errors
	if resp.IsErrorState() {
		logger.Error("listPackages", "Status", resp.Status, "Body", resp.String())
		return fmt.Errorf("error from cortex API %s: %s", resp.Status, resp.String())
	}
	logger.Error("listPackages url", resp.Request.URL.String())

	// Unmarshal the response and check for unmarshal errors
	err := resp.Into(&response)
	if err != nil {
		logger.Error("listPackages", "Error", err)
		return err
	}

	// Stream each row from the response, stop if we hit the limit
	for _, result := range response {
		// send the item to steampipe
		row := CortexPackageRow{
			Name:        result.Name,
			PackageType: result.PackageType,
			PackageTag:  packageTag,
			Version:     result.Version,
		}
		writer.StreamListItem(ctx, row)
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if writer.RowsRemaining(ctx) == 0 {
			return nil
		}
	}

	return nil
}
