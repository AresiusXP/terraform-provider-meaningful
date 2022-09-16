package main

import (
	"terraform-provider-meaningful/meaningful"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: meaningful.Provider})
}