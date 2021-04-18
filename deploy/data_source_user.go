// Copyright 2021 William Perron. All rights reserved. MIT License.
package deploy

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/wperron/terraform-deploy-provider/client"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"github_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	res, err := c.CurrentUser()

	if err != nil {
		return fmt.Errorf("Error getting Current User: %w", err)
	}

	log.Printf("[DEBUG] Received Caller Identity: %s %s", res.ID, res.Name)

	d.SetId(res.ID)
	if err := d.Set("id", res.ID); err != nil {
		return err
	}
	if err := d.Set("name", res.Name); err != nil {
		return err
	}
	if err := d.Set("github_id", fmt.Sprint(res.GitHubID)); err != nil {
		return err
	}

	return nil
}
