// Copyright 2021 William Perron. All rights reserved. MIT License.
package deploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var providerConfig string

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"deploy": testAccProvider,
	}

	providerConfig = fmt.Sprintf(`
  provider "deploy" {
    api_token = "%s"
  }
`, testToken)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestAccProviderConfigure(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() {},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
