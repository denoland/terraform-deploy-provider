// Copyright 2021 William Perron. All rights reserved. MIT License.
package deploy

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/wperron/terraform-deploy-provider/client"
)

func resourceCustomDomainValidation() *schema.Resource {
	return &schema.Resource{
		Create: createCustomDomainValidation,
		Read:   readCustomDomainValidation,
		Delete: deleteCustomDomainValidation,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"custom_domain": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func createCustomDomainValidation(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	projectID := d.Get("project_id").(string)
	domainName := d.Get("custom_domain").(string)
	domain, err := c.GetDomain(projectID, domainName)
	if err != nil {
		return err
	}

	if err := c.VerifyDomain(projectID, domainName); err != nil {
		return err
	}

	if err := c.ProvisionCertificateAutomatic(projectID, domainName); err != nil {
		return err
	}

	d.SetId(domain.CreatedAt)
	return readCustomDomainValidation(d, meta)
}

func readCustomDomainValidation(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	projectID := d.Get("project_id").(string)
	domainName := d.Get("custom_domain").(string)
	domain, err := c.GetDomain(projectID, domainName)
	if err != nil {
		return err
	}

	if !domain.IsValidated || len(domain.Certificates) == 0 {
		return errors.New("domain is either not validated or does not have any certificates")
	}
	return nil
}

func deleteCustomDomainValidation(d *schema.ResourceData, meta interface{}) error {
	// since this is a logical resource, there's nothing to delete.
	return nil
}
