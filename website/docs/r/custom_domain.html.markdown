---
subcategory: "Project"
layout: "deploy"
page_title: "Deploy: custom domain"
description: |-
  Assigns a custom domain name to a Deploy Project.
---

# Resource: deploy_custom_domain

Assigns a custom domain name to a Deploy Project.

## Example Usage

```terraform
resource "deploy_project" "example" {
  name       = "my-test-project"
  source_url = "https://dash.deno.com/examples/hello.js"
}

resource "deploy_custom_domain" "this" {
  project_id  = deploy_project.this.id
  domain_name = "foo.example.org"
}
```

## Argument Reference

The following arguments are required:

* `project_id` - (Required) The project ID the domain name will be associated to.
* `domain_name` - (Required) The fully-qualified domain name.

## Attributes Reference

* `records` - A list of records that must be created and the values they should
  have for the validation process.
* `is_validated` - Boolean value showing whether the domain name has been
  validated or not.