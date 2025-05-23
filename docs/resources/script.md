---
page_title: "shell_script (Resource) - terraform-provider-shell"
subcategory: ""
description: |-
  The Shell script resource (shell_script) allows you to execute arbitrary commands as part of a Terraform lifecycle. All commands must output a JSON string to the file defined by the TF_SCRIPT_OUTPUT environment variable and the file must be consistent on re-reading. You can access the output value in state in the read, update and delete commands via the TF_STATE_OUTPUT environment variable.
---

# shell_script (Resource)

The _Shell_ script resource (`shell_script`) allows you to execute arbitrary commands as part of a _Terraform_ lifecycle. All commands must output a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable and the file must be consistent on re-reading. You can access the output value in state in the read, update and delete commands via the `TF_STATE_OUTPUT` environment variable.

## Example Usage

```terraform
resource "shell_script" "example" {
  environment = {
    "OLD_TARGET_FILE" = "my-resource"
    "TARGET_FILE"     = "my-resource-new"
  }

  commands = {
    create = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      touch "$${path}"
      printf '{"exists": true, "path": "%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
    EOF

    read = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      if [[ -f "$${path}" ]]; then
        printf '{"exists": true, "path": "%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
      else
        printf '{"exists": false, "path": "%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
      fi
    EOF

    update = <<-EOF
      set -euo pipefail
      old_path="/tmp/$${OLD_TARGET_FILE}"
      if [[ -f "$${old_path}" ]]; then
        rm -f "$${old_path}"
      fi
      path="/tmp/$${TARGET_FILE}"
      touch "$${path}"
      printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
    EOF

    delete = <<-EOF
      set -euo pipefail
      rm -f "/tmp/$${TARGET_FILE}"
    EOF
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `commands` (Attributes) The commands to run as part of the _Terraform_ lifecycle. All commands must write a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable. (see [below for nested schema](#nestedatt--commands))

### Optional

- `environment` (Map of String) The environment variables to set when executing commands; to be combined with the OS environment and the provider environment.
- `interpreter` (List of String) The interpreter to use for executing the commands; if not set the provider interpreter will be used.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `working_directory` (String) The working directory to use when executing the commands; this will default to the _Terraform_ working directory..

### Read-Only

- `output` (Dynamic) The output of the script as a structured type.

<a id="nestedatt--commands"></a>
### Nested Schema for `commands`

Required:

- `create` (String) The command to execute when creating the resource.
- `delete` (String) The command to execute when deleting the resource.
- `read` (String) The command to execute when reading the resource.
- `update` (String) The command to execute when updating the resource.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) Timeout for creating the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `delete` (String) Timeout for deleting the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `read` (String) Timeout for reading the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `update` (String) Timeout for updating the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration] (https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
