// Copyright 2021 Deno Land Inc. All rights reserved. MIT License.
package deploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/wperron/terraform-deploy-provider/client"
)

func TestAccUser_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployUserConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeployCurrentUser("data.deploy_user.current"),
				),
			},
		},
	})
}

func testAccCheckDeployCurrentUser(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find current user: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("user ID resource ID not set.")
		}

		expected, err := testAccProvider.Meta().(*client.Client).CurrentUser()
		if err != nil {
			return fmt.Errorf("failed to get current user: %s", err)
		}
		if rs.Primary.Attributes["id"] != expected.ID {
			return fmt.Errorf("incorrect user ID: expected %q, got %q", expected.ID, rs.Primary.Attributes["id"])
		}

		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("username expected to not be nil")
		}

		if rs.Primary.Attributes["github_id"] == "" {
			return fmt.Errorf("GitHub ID expected to not be nil")
		}

		return nil
	}
}

const testAccDeployUserConfig_basic = `
data "deploy_user" "current" {}
`
