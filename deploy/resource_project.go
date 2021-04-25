// Copyright 2021 William Perron. All rights reserved. MIT License.
package deploy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/wperron/terraform-deploy-provider/client"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: createProject,
		Read:   readProject,
		Update: updateProject,
		Delete: deleteProject,
		Exists: existsProject,
		// TODO(wperron) implement Importer
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"github_link"},
			},
			"github_link": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"source_url"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization": {
							Type:     schema.TypeString,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"repo": {
							Type:     schema.TypeString,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"entrypoint": {
							Type:     schema.TypeString,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"production_deployment": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"domain_mappings": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeMap,
								Elem: &schema.Schema{Type: schema.TypeString},
							},
						},
						"related_commit": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hash": {
										Type:     schema.TypeString,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"message": {
										Type:     schema.TypeString,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"author_name": {
										Type:     schema.TypeString,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"author_email": {
										Type:     schema.TypeString,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"author_github_username": {
										Type:     schema.TypeString,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"url": {
										Type:     schema.TypeString,
										Required: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"env_vars": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"updated_at": {
							Type:     schema.TypeString,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"created_at": {
							Type:     schema.TypeString,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"has_production_deployment": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"env_vars": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func createProject(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	name := d.Get("name").(string)
	vars := make(client.EnvVars)
	tmp := d.Get("env_vars").(map[string]interface{})
	for k, v := range tmp {
		vars[k] = fmt.Sprint(v)
	}

	project, err := c.CreateProject(name, vars)
	if err != nil {
		return err
	}

	if source, ok := d.GetOk("source_url"); ok {
		if _, err := c.NewProjectDeployment(project.ID, client.NewDeploymentRequest{
			URL:        source.(string),
			Production: true,
		}); err != nil {
			return err
		}
	} else if gh, ok := d.GetOk("github_link"); ok {
		ghLinkList := gh.([]interface{})
		ghLink := ghLinkList[0].(map[string]interface{})
		if _, err := c.LinkProject(client.LinkProjectRequest{
			ProjectID:    d.Id(),
			Organization: ghLink["organization"].(string),
			Repo:         ghLink["repo"].(string),
			Entrypoint:   ghLink["entrypoint"].(string),
		}); err != nil {
			return err
		}
	}

	d.SetId(project.ID)
	return readProject(d, meta)
}

func readProject(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	project, err := c.GetProject(d.Id())
	if err != nil {
		return err
	}

	if project.ProductionDeployment != nil {
		if err := d.Set("production_deployment", productionDeploymentToTerraformSchema(project.ProductionDeployment)); err != nil {
			return err
		}
		if source, ok := d.GetOk("source_url"); ok && source != project.ProductionDeployment.URL {
			if err := d.Set("source_url", project.ProductionDeployment.URL); err != nil {
				return err
			}
		}
	}
	if err := d.Set("has_production_deployment", project.HasProductionDeployment); err != nil {
		return err
	}

	if project.Git != nil {
		if err := d.Set("github_link", []map[string]interface{}{
			{
				"organization": project.Git.Repository.Owner,
				"repo":         project.Git.Repository.Name,
				"entrypoint":   project.Git.Entrypoint,
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func updateProject(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	if d.HasChange("name") {
		name := d.Get("name").(string)
		if err := c.UpdateProject(d.Id(), name); err != nil {
			return err
		}
	}

	if d.HasChange("env_vars") {
		vars := make(client.EnvVars)
		tmp := d.Get("env_vars").(map[string]interface{})
		for k, v := range tmp {
			vars[k] = fmt.Sprint(v)
		}

		if err := c.UpdateEnvVars(d.Id(), vars); err != nil {
			return err
		}
	}

	if d.HasChange("source_url") {
		if source, ok := d.GetOk("source_url"); ok {
			if _, err := c.NewProjectDeployment(d.Id(), client.NewDeploymentRequest{
				URL:        source.(string),
				Production: true,
			}); err != nil {
				return err
			}
		}
	}

	if d.HasChange("github_link") {
		if gh, ok := d.GetOk("github_link"); ok {
			ghLinkList := gh.([]interface{})
			ghLink := ghLinkList[0].(map[string]interface{})
			if _, err := c.LinkProject(client.LinkProjectRequest{
				ProjectID:    d.Id(),
				Organization: ghLink["organization"].(string),
				Repo:         ghLink["repo"].(string),
				Entrypoint:   ghLink["entrypoint"].(string),
			}); err != nil {
				return err
			}
		}
	}

	return readProject(d, meta)
}

func deleteProject(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	return c.DeleteProject(d.Id())
}

func existsProject(d *schema.ResourceData, meta interface{}) (bool, error) {
	c := meta.(*client.Client)
	if _, err := c.GetProject(d.Id()); err != nil {
		return false, err
	}
	return true, nil
}

func productionDeploymentToTerraformSchema(depl *client.Deployment) []interface{} {
	if depl == nil {
		return nil
	}

	tfMap := map[string]interface{}{}
	tfMap["id"] = depl.ID
	tfMap["url"] = depl.URL
	domains := []interface{}{}
	for _, domain := range depl.DomainMappings {
		domains = append(domains, map[string]interface{}{
			"domain":     domain.Domain,
			"updated_at": domain.UpdatedAt,
			"created_at": domain.CreatedAt,
		})
	}
	tfMap["domain_mappings"] = domains

	if depl.RelatedCommit != nil {
		tfMap["related_commit"] = []map[string]interface{}{
			{
				"hash":                   depl.RelatedCommit.Hash,
				"message":                depl.RelatedCommit.Message,
				"author_name":            depl.RelatedCommit.AuthorName,
				"author_email":           depl.RelatedCommit.AuthorEmail,
				"author_github_username": depl.RelatedCommit.AuthorGitHubUsername,
				"url":                    depl.RelatedCommit.URL,
			},
		}
	}

	tfMap["env_vars"] = depl.EnvVars
	tfMap["updated_at"] = depl.UpdatedAt
	tfMap["created_at"] = depl.CreatedAt

	return []interface{}{tfMap}
}
