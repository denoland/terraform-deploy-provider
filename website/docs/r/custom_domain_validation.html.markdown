---
subcategory: "Project"
layout: "deploy"
page_title: "Deploy: custom domain"
description: |-
  Coordinates and waits for a custom domain to validated and provisioned with
  certificates.
---

# Resource: deploy_custom_domain_validation

This resource represents a successful validation of a custom domain name and the
complete provisioning of TLS certificates for that domain.

This would be used in concert with other resources such as [AWS's Route53][1],
[Google DNS][2] or [CloudFlare][3] to manage the whole lifecycle of creating a
new domain, adding it to Deno Deploy, creating the validation records and
provisioning TLS certificates.

* **WARNING:** This resource implements the validation workflow only. It does
  not represent a real-world resource in Deno Deploy. Changing or deleting this
  resource will not have any immediate effect.

## Example Usage

### DNS Validation using AWS Route53

```terraform
resource "deploy_project" "this" {
  name       = "my-test-project"
  source_url = "https://dash.deno.com/examples/hello.js"
}

resource "deploy_custom_domain" "this" {
  project_id  = deploy_project.this.id
  domain_name = "foo.example.org"
}

data "aws_route53_zone" "this" {
  name = "example.org."
}

resource "aws_route53_record" "example" {
  for_each = {
    for domain in deploy_custom_domain.this.records : domain.type => {
      name  = domain.domain_name
      value = domain.value
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.value]
  ttl             = 60
  type            = each.key
  zone_id         = data.aws_route53_zone.this.zone_id
}

resource "deploy_custom_domain_validation" "this" {
  project_id    = deploy_project.this.id
  custom_domain = deploy_custom_domain.this.domain_name

  # It's important to explicitely define the dependency on the DNS provider to
  # Terraform so that the resources get created in the correct order.
  depends_on = [aws_route53_record.example]
}
```

## Argument Reference

The following arguments are required:

* `project_id` - (Required) The project ID the domain name is linked to.
* `custom_domain` - (Required) The custom domain name to validate.