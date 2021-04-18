// Copyright 2021 William Perron. All rights reserved. MIT License.
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/wperron/terraform-deploy-provider/deploy"
)

func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: deploy.Provider,
	}

	plugin.Serve(opts)
}
