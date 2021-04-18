// Copyright 2021 William Perron. All rights reserved. MIT License.
package deploy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/wperron/terraform-deploy-provider/client"
)

// Provider for Deno Deploy
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API Token used for accessing Deno Deploy",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"deploy_project": resourceProject(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"deploy_user": dataSourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("api_token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := client.New(token)

	return c, diags
}
