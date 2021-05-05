---
subcategory: ""
layout: "deploy"
page_title: "Deploy: User"
description: |-
  Get information on the current user
---

# Data Source: deploy_user

This data source can be used to fetch information the owner of the token used
by the provider's configuration.

## Example Usage

```terraform
data "deploy_user" "current" {}

output "user_id" {
  value = data.deploy_user.current.id
}

output "name" {
  value = data.deploy_user.current.name
}

output "github_id" {
  value = data.deploy_user.current.github_id
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

* `id` - The UUID of the user.
* `name` - The GitHub username of the user.
* `github_id` - The GitHub numeric ID of the user.