---
subcategory: "Project"
layout: "deploy"
page_title: "Deploy: Project"
description: |-
  Provides a Deno Deploy Project.
---

# Resource: deploy_project

Provides a Deno Deploy Project.

## Example Usage

### Basic Example

```terraform
resource "deploy_project" "example" {
  name       = "my-test-project"
  source_url = "https://dash.deno.com/examples/hello.js"
}
```

### Linking a GitHub Repo

```terraform
resource "deploy_project" "example" {
  name = "my-test-project"
  github_link {
    organization = "username"
    repo         = "my-repo"
    entrypoint   = "/main.ts"
  }
}
```

### Environment Variables

```terraform
resource "deploy_project" "example" {
  name = "my-test-project"
  source_url = "https://dash.deno.com/examples/hello.js"

  env_var {
    key   = "foo"
    value = "bar"
  }

  env_var {
    key   = "greeting"
    value = "Hello World!"
  }
}
```

## Argument Reference

The following arguments are required:

* `name` - (Required) Unique name for your project.

The following arguments are optional:

* `source_url` - (Optional) The URL where the entrypoint for the project is
  located. Conflicts with `github_link`
* `github_link` - (Optional) Configuration block. Described below. Conflicts
  with `source_url`
* `env_var` - (Optional) Configuration block. Described below.

### github_link

GitHub link configuration that specifies the GitHub repository to link to the
project as well as its entrypoint script.

* `organization` - (Required) The name of the organization or user the
  repository belongs to. The current user must have admin access to the
  repository.
* `repo` - (Required) The name of the repository.
* `entrypoint` - (Required) Absolute path to the entrypoint of the project. Must
  start with a forward slash (`/`).

### env_var

Environment variable to add to the project configuration. These are available
within the script through [`Deno.env`][1]

* `key` - (Required) The name of the environment variable
* `value` - (Required) The value associated with the `key`. This is a sensitive
  value and won't show up in the logs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `project_id` - The UUID of the project.
* `production_deployment` - Detailed overview of the current production
  deployment of the project.
* `has_production_deployment` - Boolean showing whether the project has a
  production deployment or not.

[1]: https://doc.deno.land/builtin/stable#Deno.env