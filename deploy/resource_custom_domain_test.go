// Copyright 2021 William Perron. All rights reserved. MIT License.
package deploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/wperron/terraform-deploy-provider/client"
)

func TestAccCustomDomain_basic(t *testing.T) {
	randomID := acctest.RandStringFromCharSet(4, acctest.CharSetAlphaNum)

	config := fmt.Sprintf(testAccCustomDomainConfig_basic, randomID, randomID)
	randomDomain := fmt.Sprintf("foo-%s.example.org", randomID)

	var project client.Project
	var domain client.Domain
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCustomDomainCheckDestroy(&project, &domain),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCustomDomainCheckExists("deploy_custom_domain.test", &project, &domain),
					resource.TestCheckResourceAttrSet(
						"deploy_custom_domain.test", "project_id",
					),
					resource.TestCheckResourceAttr(
						"deploy_custom_domain.test", "domain_name", randomDomain,
					),
					resource.TestCheckResourceAttr(
						"deploy_custom_domain.test", "is_validated", "false",
					),
					// TODO(wperron) figure out how to do this check properly
					// resource.TestCheckTypeSetElemNestedAttrs(
					// 	"deploy_custom_domain.test", "records.*", map[string]string{
					// 		"domain_name": randomDomain,
					// 		"type":        "*",
					// 		"value":       "*",
					// 	},
					// ),
				),
			},
		},
	})
}

func testAccCustomDomainCheckExists(rn string, p *client.Project, d *client.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}
		domain := string(rs.Primary.ID)

		projectID, ok := rs.Primary.Attributes["project_id"]
		if !ok {
			return fmt.Errorf("property `project_id` is not set")
		}
		client := testAccProvider.Meta().(*client.Client)
		project, err := client.GetProject(projectID)
		if err != nil {
			return fmt.Errorf("failed to get project info: %s", err)
		}
		customDomain, err := client.GetDomain(project.ID, domain)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}
		if customDomain.Domain == "" {
			return fmt.Errorf("custom domain has no domain name")
		}
		*d = customDomain
		return nil
	}
}

func testAccCustomDomainCheckDestroy(p *client.Project, d *client.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		domain, err := client.GetDomain(p.ID, d.Domain)
		if err == nil && domain.Domain != "" {
			return fmt.Errorf("project still exists")
		}
		return nil
	}
}

const testAccCustomDomainConfig_basic = `
resource "deploy_project" "test" {
  name       = "terraform-test-%s"
  source_url = "https://dash.deno.com/examples/hello.js"
}

resource "deploy_custom_domain" "test" {
  project_id = deploy_project.test.id
  domain_name = "foo-%s.example.org"
}
`
