package pkg

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const DefaultBaseURL = "https://api.getcortexapp.com"

type SteampipeConfig struct {
	ApiKey  *string `cty:"api_key"`
	BaseURL *string `cty:"base_url"`
}

func NewSteampipeConfig(token, url string) *SteampipeConfig {
	return &SteampipeConfig{ApiKey: &token, BaseURL: &url}
}

func GetConfig(connection *plugin.Connection) *SteampipeConfig {
	if connection == nil || connection.Config == nil {
		return NewSteampipeConfig("", DefaultBaseURL)
	}
	config, _ := connection.Config.(SteampipeConfig)

	// Read the API key from the environment and override the value in the config
	token, ok := os.LookupEnv("CORTEX_API_KEY")
	if ok {
		config.ApiKey = &token
	}

	// Read the base URL from the environment and override the value in the config
	baseURL, ok := os.LookupEnv("CORTEX_BASE_URL")
	if ok {
		config.BaseURL = &baseURL
	}

	return &config
}

func SteampipePlugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-cortex",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: func() interface{} { return NewSteampipeConfig("", DefaultBaseURL) },
			Schema:      map[string]*schema.Attribute{"api_key": {Type: schema.TypeString}},
		},
		TableMap: map[string]*plugin.Table{
			"cortex_descriptor": tableCortexDescriptor(),
			"cortex_entity":     tableCortexEntity(),
		},
	}
	return p
}
