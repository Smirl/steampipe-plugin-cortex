package pkg

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type SteampipeConfig struct {
	ApiKey *string `cty:"api_key"`
}

func GetConfig(connection *plugin.Connection) SteampipeConfig {
	if connection == nil || connection.Config == nil {
		return SteampipeConfig{}
	}
	config, _ := connection.Config.(SteampipeConfig)
	return config
}

func SteampipePlugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-cortex",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: func() interface{} { return &SteampipeConfig{} },
			Schema:      map[string]*schema.Attribute{"api_key": {Type: schema.TypeString}},
		},
		TableMap: map[string]*plugin.Table{
			"cortex_descriptors": tableCortexDescriptors(),
			"cortex_entities":    tableCortexEntities(),
		},
	}
	return p
}
