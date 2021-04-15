package denodeploy

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	deployclient "github.com/wperron/terraform-deno-deploy-provider/deploy-client"
)

func DataSourceUser() *schema.Resource {
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
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*deployclient.Client)

	res, err := c.CurrentUser()

	if err != nil {
		return fmt.Errorf("Error getting Current User: %w", err)
	}

	log.Printf("[DEBUG] Received Caller Identity: %s %s", res.Id, res.Name)

	d.SetId(res.Id)
	d.Set("id", res.Id)
	d.Set("name", res.Name)
	d.Set("github_id", res.GitHubID)

	return nil
}
