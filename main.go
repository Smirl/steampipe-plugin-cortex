package main

import (
	"github.com/smirl/steampipe-plugin-cortex/pkg"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: pkg.SteampipePlugin})
}
