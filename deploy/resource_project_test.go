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

func TestAccProject_basic(t *testing.T) {
	randomID := acctest.RandStringFromCharSet(4, acctest.CharSetAlphaNum)

	config := fmt.Sprintf(testAccProjectConfig_basic, randomID)
	updated := fmt.Sprintf(testAccProjectConfig_update, randomID)

	var project client.Project
	source := "https://dash.deno.com/examples/hello.js"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccProjectCheckDestroy(&project),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("deploy_project.test", &project),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "name", fmt.Sprintf("terraform-test-%s", randomID),
					),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "has_production_deployment", "false",
					),
					resource.TestCheckResourceAttrSet(
						"deploy_project.test", "id",
					),
				),
			},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("deploy_project.test", &project),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "name", fmt.Sprintf("terraform-test-%s", randomID),
					),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "has_production_deployment", "true",
					),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "source_url", "https://dash.deno.com/examples/hello.js",
					),
					resource.TestCheckResourceAttrSet(
						"deploy_project.test", "id",
					),
					testAccProjectDeployment(&project, testProductionDeployment{
						SourceUrl: &source,
						EnvVars:   make(client.NewEnvVars),
					}),
				),
			},
		},
	})
}

func TestAccProject_envVars(t *testing.T) {
	randomID := acctest.RandStringFromCharSet(4, acctest.CharSetAlphaNum)

	config := fmt.Sprintf(testAccProjectConfig_envVars, randomID)

	var project client.Project
	source := "https://dash.deno.com/examples/hello.js"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccProjectCheckDestroy(&project),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("deploy_project.test", &project),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "name", fmt.Sprintf("terraform-test-%s", randomID),
					),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "has_production_deployment", "true",
					),
					resource.TestCheckResourceAttrSet(
						"deploy_project.test", "id",
					),
					testAccProjectDeployment(&project, testProductionDeployment{
						SourceUrl: &source,
						EnvVars: client.NewEnvVars{
							"foo":   "bar",
							"fruit": "banana",
						},
					}),
				),
			},
		},
	})
}

func TestAccProject_github(t *testing.T) {
	randomID := acctest.RandStringFromCharSet(4, acctest.CharSetAlphaNum)

	config := fmt.Sprintf(testAccProjectConfig_github, randomID)

	var project client.Project
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccProjectCheckDestroy(&project),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("deploy_project.test", &project),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "name", fmt.Sprintf("terraform-test-%s", randomID),
					),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "has_production_deployment", "true",
					),
					resource.TestCheckResourceAttrSet(
						"deploy_project.test", "id",
					),
					testAccProjectDeployment(&project, testProductionDeployment{
						EnvVars: make(client.NewEnvVars),
						GitHub: &testGitHub{
							Org:        "wperron",
							Repo:       "terraform-deploy-provider",
							Entrypoint: "/deploy/testdata/main.ts",
						},
					}),
				),
			},
		},
	})
}

func TestAccProject_linkAndUnlink(t *testing.T) {
	randomID := acctest.RandStringFromCharSet(4, acctest.CharSetAlphaNum)

	config := fmt.Sprintf(testAccProjectConfig_update, randomID)
	linked := fmt.Sprintf(testAccProjectConfig_github, randomID)

	var project client.Project
	source := "https://dash.deno.com/examples/hello.js"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccProjectCheckDestroy(&project),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("deploy_project.test", &project),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "name", fmt.Sprintf("terraform-test-%s", randomID),
					),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "has_production_deployment", "true",
					),
					resource.TestCheckResourceAttrSet(
						"deploy_project.test", "id",
					),
					testAccProjectDeployment(&project, testProductionDeployment{
						SourceUrl: &source,
						EnvVars:   make(client.NewEnvVars),
					}),
				),
			},
			{
				Config: linked,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("deploy_project.test", &project),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "name", fmt.Sprintf("terraform-test-%s", randomID),
					),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "has_production_deployment", "true",
					),
					resource.TestCheckResourceAttrSet(
						"deploy_project.test", "id",
					),
					testAccProjectDeployment(&project, testProductionDeployment{
						EnvVars: make(client.NewEnvVars),
						GitHub: &testGitHub{
							Org:        "wperron",
							Repo:       "terraform-deploy-provider",
							Entrypoint: "/deploy/testdata/main.ts",
						},
					}),
				),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccProjectCheckExists("deploy_project.test", &project),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "name", fmt.Sprintf("terraform-test-%s", randomID),
					),
					resource.TestCheckResourceAttr(
						"deploy_project.test", "has_production_deployment", "true",
					),
					resource.TestCheckResourceAttrSet(
						"deploy_project.test", "id",
					),
					testAccProjectDeployment(&project, testProductionDeployment{
						SourceUrl: &source,
						EnvVars:   make(client.NewEnvVars),
						GitHub:    nil,
					}),
				),
			},
		},
	})
}

func testAccProjectCheckExists(rn string, p *client.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}
		id := string(rs.Primary.ID)
		client := testAccProvider.Meta().(*client.Client)
		project, err := client.GetProject(id)
		if err != nil {
			return fmt.Errorf("error getting data source: %s", err)
		}
		if project.ID == "" {
			return fmt.Errorf("project has no ID")
		}
		*p = project
		return nil
	}
}

func testAccProjectCheckDestroy(p *client.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*client.Client)
		project, err := client.GetProject(p.ID)
		if err == nil && project.Name != "" {
			return fmt.Errorf("project still exists")
		}
		return nil
	}
}

type testProductionDeployment struct {
	SourceUrl *string
	EnvVars   client.NewEnvVars
	GitHub    *testGitHub
}

type testGitHub struct {
	Org        string
	Repo       string
	Entrypoint string
}

func testAccProjectDeployment(p *client.Project, expected testProductionDeployment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if p == nil {
			return fmt.Errorf("cannot check production deployment: project does not exist")
		}
		if p.ProductionDeployment == nil {
			return fmt.Errorf("no production deployment found")
		}
		if !includes(p.ProductionDeployment.EnvVars, "DENO_DEPLOYMENT_ID") {
			return fmt.Errorf("production deployment doesn't have a `DENO_DEPLOYMENT_ID` environment variable")
		}
		if expected.SourceUrl != nil && *expected.SourceUrl != p.ProductionDeployment.URL {
			return fmt.Errorf("expected production deployment with a source url %s, found %s", *expected.SourceUrl, p.ProductionDeployment.URL)
		}
		if len(expected.EnvVars) > 0 {
			for k := range expected.EnvVars {
				if !includes(p.ProductionDeployment.EnvVars, k) {
					return fmt.Errorf("could not find the expected %s environment variable", k)
				}
			}
		}
		if expected.GitHub != nil {
			if p.Git == nil {
				return fmt.Errorf("expected github %s/%s%s to be linked, found none", expected.GitHub.Org, expected.GitHub.Repo, expected.GitHub.Entrypoint)
			}
			if p.ProductionDeployment.RelatedCommit == nil {
				return fmt.Errorf("expected production deployment to have a related commit, found none")
			}
		}
		return nil
	}
}

const testAccProjectConfig_basic = `
resource "deploy_project" "test" {
  name = "terraform-test-%s"
}
`

const testAccProjectConfig_update = `
resource "deploy_project" "test" {
  name       = "terraform-test-%s"
  source_url = "https://dash.deno.com/examples/hello.js"
}
`

const testAccProjectConfig_envVars = `
resource "deploy_project" "test" {
  name       = "terraform-test-%s"
  source_url = "https://dash.deno.com/examples/hello.js"

  env_var {
    key   = "foo"
    value = "bar"
  }

  env_var {
    key   = "fruit"
    value = "banana"
  }
}
`

const testAccProjectConfig_github = `
resource "deploy_project" "test" {
  name = "terraform-test-%s"
  github_link {
    organization = "wperron"
    repo         = "terraform-deploy-provider"
    entrypoint   = "/deploy/testdata/main.ts"
  }
}
`

func includes(l []string, s string) bool {
	for _, e := range l {
		if e == s {
			return true
		}
	}
	return false
}
