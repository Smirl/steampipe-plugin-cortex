package cortex

import (
	"testing"
	_ "unsafe"

	. "github.com/onsi/gomega"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func TestGetConfig(t *testing.T) {
	g := NewWithT(t)
	testApiKey := "test_api_key"
	testBaseURL := "https://test-url.com"
	connection := &plugin.Connection{
		Config: SteampipeConfig{
			ApiKey:  &testApiKey,
			BaseURL: &testBaseURL,
		},
	}

	config := GetConfig(connection)

	g.Expect(*config.ApiKey).To(Equal("test_api_key"))
	g.Expect(*config.BaseURL).To(Equal("https://test-url.com"))
}

func TestGetConfigWithEnvVars(t *testing.T) {
	g := NewWithT(t)
	connection := &plugin.Connection{
		Config: SteampipeConfig{},
	}
	t.Setenv("CORTEX_API_KEY", "env_api_key")
	t.Setenv("CORTEX_BASE_URL", "https://env-url.com")

	config := GetConfig(connection)

	g.Expect(*config.ApiKey).To(Equal("env_api_key"))
	g.Expect(*config.BaseURL).To(Equal("https://env-url.com"))
}
