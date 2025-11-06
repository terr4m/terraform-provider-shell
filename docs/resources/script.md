---
page_title: "shell_script (Resource) - terraform-provider-shell"
subcategory: ""
description: |-
  The Shell script resource (shell_script) allows you to execute arbitrary commands as part of a Terraform lifecycle. All commands must output a JSON string to the file defined by the TF_SCRIPT_OUTPUT environment variable and the file must be consistent on re-reading. If a script exits with a non-zero code the provider will ready any text from the file defined by the TF_SCRIPT_ERROR environment variable and return it as part of the error diagnostics.
---

# shell_script (Resource)

The _Shell_ script resource (`shell_script`) allows you to execute arbitrary commands as part of a _Terraform_ lifecycle. All commands must output a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable and the file must be consistent on re-reading. If a script exits with a non-zero code the provider will ready any text from the file defined by the `TF_SCRIPT_ERROR` environment variable and return it as part of the error diagnostics.

## Environment Variables

The following environment variables provide the shell script integration with the provider.

| **Name** | **Description** |
| :--- | :--- |
| `TF_SCRIPT_LIFECYCLE` | The current lifecycle that triggered the script; this can be one of `create`, `read`, `update`, or `delete`. |
| `TF_SCRIPT_INPUTS` | The values passed into the data source `inputs` as JSON. |
| `TF_SCRIPT_OUTPUT` | Path to the file where the script output must be written; the output must be valid JSON. |
| `TF_SCRIPT_ERROR` | Path to a file which will be read as the error diagnostics if the scripts exits with a non-zero code. |
| `TF_SCRIPT_STATE_OUTPUT` | The current value of `output` in the state file, as JSON. |

## Capabilities

This resource supports the following capabilities to bring script functionality closer to native Terraform resources.

### JSON Inputs

Scripts receive input parameters as JSON via the `TF_SCRIPT_INPUTS` environment variable, simplifying data handling.

### JSON Outputs

Scripts must write their output as JSON to the file specified by the `TF_SCRIPT_OUTPUT` environment variable, ensuring structured data exchange. There is a special `__meta` key that can be used to provide additional metadata back to the provider.

### Plan Customization

Scripts can customize the plan phase of the Terraform lifecycle by providing a plan command configuration, allowing for more dynamic resource management. If the plan returns a `"???"` value for any `output` key this will be treated as a dynamic value that requires re-evaluation during the apply phase. This allows scripts to signal that certain outputs cannot be determined until the resource is actually created or updated.

### Output Drift Detection

If an update is required to correct the state of the `output` values, the read script can set the `output.__meta.output_drift_detected` key to `true`. This will allow the provider to set the `output_drift` attribute to `true` and trigger an update during the next apply.

### State Awareness

Scripts can access the current state output via the `TF_SCRIPT_STATE_OUTPUT` environment variable, allowing for more informed operations during updates or deletions.

### Lifecycle Awareness

By inspecting the `TF_SCRIPT_LIFECYCLE` environment variable, scripts can adapt their behavior based on the current lifecycle phase.

## Example Usage

```terraform
resource "shell_script" "example" {
  inputs = {
    file_name = "foo"
  }

  os_commands = {
    default = {
      plan = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      create = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          touch "$${path}"
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          if [[ -f "$${path}" ]]; then
            printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false,"path":"%s","__meta":{"output_drift_detected":true}}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          fi
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          old_path="$(jq --raw-output '.path' <<<"$${TF_SCRIPT_STATE_OUTPUT}")"
          if [[ "$${path}" != "$${old_path}" ]] && [[ -f "$${old_path}" ]]; then
            mv -f "$${old_path}" "$${path}"
          else
            touch "$${path}"
          fi
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          rm -f "$${path}"
        EOF
      }
    }
    windows = {
      plan = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      create = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          if (Test-Path $path) {
            @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          } else {
            @{exists=$false; path=$path; __meta=@{output_drift_detected=$true}} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          }
        EOF
      }
      update = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $oldPath = $state.path
          if ($path -ne $oldPath) {
            Move-Item -Path $oldPath -Destination $path -Force
          } else {
            New-Item -Path $path -ItemType File -Force
          }
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          Remove-Item -Path $path -Force -ErrorAction Ignore
        EOF
      }
    }
  }

  output_drift = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `os_commands` (Attributes Map) A map of commands to run as part of the Terraform lifecycle where the map key is the `GOOS` value or `default`; `default` must be provided. (see [below for nested schema](#nestedatt--os_commands))
- `output_drift` (Boolean) This is used by the provider to manage the output state and must always be set to false.

### Optional

- `environment` (Map of String) The environment variables to set when executing commands; to be combined with the OS environment and the provider environment.
- `inputs` (Dynamic) Inputs to be made available to the script; these can be accessed as JSON via the `TF_SCRIPT_INPUTS` environment variable.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `triggers` (Dynamic) Allows specifying values that trigger resource replacement when changed.
- `working_directory` (String) The working directory to use when executing the commands; this will default to the _Terraform_ working directory.

### Read-Only

- `output` (Dynamic) The output of the script as a structured type; this can be accessed in the read, update and delete commands as JSON via the `TF_SCRIPT_STATE_OUTPUT` environment variable.

<a id="nestedatt--os_commands"></a>
### Nested Schema for `os_commands`

Required:

- `create` (Attributes) The create command configuration. (see [below for nested schema](#nestedatt--os_commands--create))
- `delete` (Attributes) The delete command configuration. (see [below for nested schema](#nestedatt--os_commands--delete))
- `read` (Attributes) The read command configuration. (see [below for nested schema](#nestedatt--os_commands--read))
- `update` (Attributes) The update command configuration. (see [below for nested schema](#nestedatt--os_commands--update))

Optional:

- `plan` (Attributes) The plan command configuration, this can be used to customize the plan phase of the Terraform lifecycle. (see [below for nested schema](#nestedatt--os_commands--plan))

<a id="nestedatt--os_commands--create"></a>
### Nested Schema for `os_commands.create`

Required:

- `command` (String) The create command to execute.

Optional:

- `interpreter` (List of String) The interpreter to use for executing the create command; if not set the platform default interpreter will be used.


<a id="nestedatt--os_commands--delete"></a>
### Nested Schema for `os_commands.delete`

Required:

- `command` (String) The delete command to execute.

Optional:

- `interpreter` (List of String) The interpreter to use for executing the delete command; if not set the platform default interpreter will be used.


<a id="nestedatt--os_commands--read"></a>
### Nested Schema for `os_commands.read`

Required:

- `command` (String) The read command to execute.

Optional:

- `interpreter` (List of String) The interpreter to use for executing the read command; if not set the platform default interpreter will be used.


<a id="nestedatt--os_commands--update"></a>
### Nested Schema for `os_commands.update`

Required:

- `command` (String) The update command to execute.

Optional:

- `interpreter` (List of String) The interpreter to use for executing the update command; if not set the platform default interpreter will be used.


<a id="nestedatt--os_commands--plan"></a>
### Nested Schema for `os_commands.plan`

Required:

- `command` (String) The plan command to execute.

Optional:

- `interpreter` (List of String) The interpreter to use for executing the plan command; if not set the platform default interpreter will be used.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) Timeout for creating the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `delete` (String) Timeout for deleting the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `read` (String) Timeout for reading the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
- `update` (String) Timeout for updating the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).
