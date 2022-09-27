package main

import (
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/incentro-dc/terraform-provider-commercelayer/commercelayer"
)

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: commercelayer.Provider}

	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "registry.terraform.io/incentro-dc/commercelayer"
	}

	plugin.Serve(opts)
}
