package main

import (
	"github.com/smirl/steampipe-plugin-cortex/cortex"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: cortex.Plugin})
}
