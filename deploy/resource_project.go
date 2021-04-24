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
			// TODO(wperron) `git` should map to the Git struct of the client
			"git": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"production_deployment": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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

	d.SetId(project.ID)
	return readProject(d, meta)
}

func readProject(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*client.Client)
	project, err := c.GetProject(d.Id())
	if err != nil {
		return err
	}

	if project.Git != nil {
		if err := d.Set("git", fmt.Sprint(project.Git.Repository.ID)); err != nil {
			return err
		}
	}
	if project.ProductionDeployment != nil {
		if err := d.Set("production_deployment", project.ProductionDeployment.ID); err != nil {
			return err
		}
	}
	if err := d.Set("has_production_deployment", project.HasProductionDeployment); err != nil {
		return err
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
