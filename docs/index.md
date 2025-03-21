---
page_title: "Shell Provider - terraform-provider-shell"
subcategory: ""
description: |-
  The Shell provider allows you to execute arbitrary shell scripts and parse their JSON output for use in your Terraform configurations. This is particularly useful for running scripts that interact with external APIs, or other systems that don't have a native Terraform provider, or for performing complex data transformations.
---

# Shell Provider

The _Shell_ provider allows you to execute arbitrary shell scripts and parse their JSON output for use in your _Terraform_ configurations. This is particularly useful for running scripts that interact with external APIs, or other systems that don't have a native _Terraform_ provider, or for performing complex data transformations.

## Example Usage

```terraform
provider "shell" {
  interpreter = ["/bin/bash", "-c"]
}

data "shell_script" "example" {
  environment = {
    "TARGET" = "my-resource"
  }

  command = file("${path.module}/scripts/read.sh")
}

resource "shell_script" "example" {
  environment = {
    "TARGET" = "my-resource"
  }

  commands = {
    create = file("${path.module}/scripts/create.sh")
    read   = file("${path.module}/scripts/read.sh")
    update = file("${path.module}/scripts/update.sh")
    delete = file("${path.module}/scripts/delete.sh")
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `environment` (Map of String) The environment variables to set when executing scripts.
- `interpreter` (List of String) The interpreter to use for executing scripts if not provided by the resource or data source. This defaults to `["/bin/bash", "-c"]` or `["pwsh", "-c"]` on Windows.
- `log_output` (Boolean) If `true`, lines output by the script will be logged at the appropriate level if they start with the `[<LEVEL>]` pattern where `<LEVEL>` can be one of `ERROR`, `WARN`, `INFO`, `DEBUG` & `TRACE`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) Timeout for resource creation; defaults to `10m`. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `delete` (String) Timeout for resource deletion; defaults to `10m`. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `read` (String) Timeout for resource or data source reads; defaults to `10m`. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `update` (String) Timeout for resource update; defaults to `10m`. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
