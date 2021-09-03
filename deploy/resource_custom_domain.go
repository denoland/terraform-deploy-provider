// Copyright 2021 Deno Land Inc. All rights reserved. MIT License.
package deploy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/wperron/terraform-deploy-provider/client"
)

func resourceCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: createCustomDomain,
		Read:   readCustomDomain,
		Delete: deleteCustomDomain,
		// TODO(wperron) implement Importer
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"is_validated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func createCustomDomain(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	project := d.Get("project_id").(string)
	domain := d.Get("domain_name").(string)

	if _, err := c.AddDomain(project, client.Domain{
		Domain: domain,
	}); err != nil {
		return err
	}

	d.SetId(domain)
	return readCustomDomain(d, meta)
}

func readCustomDomain(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	domain, err := c.GetDomain(d.Get("project_id").(string), d.Get("domain_name").(string))
	if err != nil {
		return err
	}

	if err := d.Set("records", []map[string]interface{}{
		{
			"domain_name": domain.Domain,
			"type":        "A",
			"value":       "34.120.54.55",
		},
		{
			"domain_name": domain.Domain,
			"type":        "AAAA",
			"value":       "2600:1901:0:6d85::",
		},
		{
			"domain_name": domain.Domain,
			"type":        "TXT",
			"value":       fmt.Sprintf("deno-com-validation=%s", domain.Token),
		},
	}); err != nil {
		return err
	}

	if err := d.Set("is_validated", domain.IsValidated); err != nil {
		return err
	}
	return nil
}

func deleteCustomDomain(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	return c.DeleteDomain(d.Get("project_id").(string), d.Id())
}
