---
page_title: "shell_script (Resource) - terraform-provider-shell"
subcategory: ""
description: |-
  The Shell script resource (shell_script) runs configured commands during the Terraform lifecycle.
---

# shell_script (Resource)

The _Shell_ script resource (`shell_script`) runs configured commands during the _Terraform_ lifecycle.

The `create`, `read`, and `update` commands must write one valid JSON value to `TF_SCRIPT_OUTPUT` on every successful run. The optional `plan` command must also write JSON when configured. Output is stored in Terraform state; the `read` command should return a consistent representation of current external state. If the object is removed externally, report it using `__meta.output_drift_detected = true`. A top-level `__meta` object is reserved and removed from `output`. The `delete` command does not return output. Non-zero command exits include text from `TF_SCRIPT_ERROR` in diagnostics.

## Environment Variables

The following environment variables provide the shell script integration with the provider.

| **Name** | **Description** |
| :--- | :--- |
| `TF_SCRIPT_LIFECYCLE` | The current lifecycle stage: `plan`, `create`, `read`, `update`, or `delete`. |
| `TF_SCRIPT_INPUTS` | The values passed into the resource `inputs` as JSON. |
| `TF_SCRIPT_OUTPUT` | Path where `create`, `read`, and `update` must write valid JSON. A configured `plan` command must also write JSON when it runs. `delete` output is ignored. |
| `TF_SCRIPT_ERROR` | Path to a file whose contents are included in error diagnostics when a command exits with a non-zero code. |
| `TF_SCRIPT_STATE_OUTPUT` | The prior `output` value as JSON, available to `read`, `update`, and `delete`; also available to `plan` when prior state exists. |

## Capabilities

This resource supports the following capabilities to bring script functionality closer to native Terraform resources.

### JSON Inputs

Scripts receive input parameters as JSON via the `TF_SCRIPT_INPUTS` environment variable, simplifying data handling. Inputs may be stored in Terraform state; do not pass secrets unless you understand the state and logging exposure.

### JSON Outputs

The `plan`, `create`, `read`, and `update` commands must write one valid JSON value to the file specified by `TF_SCRIPT_OUTPUT` on every successful run. The provider parses that value as the resource's `output`; stdout is not used as output. The `delete` command does not need to write `TF_SCRIPT_OUTPUT`, and the provider ignores any output it writes.

A top-level `__meta` object is reserved for provider metadata and is removed before the remaining JSON becomes `output`. See [Output Drift Detection](#output-drift-detection).

### Plan Customization

Each OS command configuration can provide an optional `plan` command. When configured, Terraform runs it during planning and uses its JSON output as the planned `output`. It receives `TF_SCRIPT_LIFECYCLE=plan`, `TF_SCRIPT_INPUTS`, and (when prior state is available) `TF_SCRIPT_STATE_OUTPUT`. Without a plan command, the provider leaves output planning to Terraform, except that detected output drift makes `output` unknown.

A string value of `"???"` anywhere in plan output is interpreted as unknown and re-evaluated during apply. This exact string cannot be represented as a known output string from the plan command.

### Output Drift Detection

If an external change means the resource needs reconciliation, the `read` command can include a top-level `__meta` object with `"output_drift_detected": true`, for example `{"exists": false, "__meta": {"output_drift_detected": true}}`. The provider removes `__meta` from stored `output` and sets the computed `output_drift` attribute to `true`. At the next plan, a configured plan command calculates output; without one, the provider marks `output` unknown. `output_drift` does not itself guarantee that an update command will repair the external object; implement repair in your plan/create/update commands.

### State Awareness

The `read`, `update`, and `delete` commands can access the prior `output` as JSON via `TF_SCRIPT_STATE_OUTPUT`. A configured `plan` command also receives it when prior state is available.

### Lifecycle Awareness

By inspecting the `TF_SCRIPT_LIFECYCLE` environment variable, scripts can adapt their behavior based on the current lifecycle phase.

Commands execute on the machine running Terraform, as the Terraform process user, with its inherited environment and filesystem access. The default interpreter is `/bin/bash -c` on non-Windows systems and `pwsh -c` on Windows; command configurations can override it. Required interpreters and utilities must be installed on that machine. The examples use Bash, `jq`, temporary files, or PowerShell as appropriate and may need adaptation. Treat command strings as executable code and use only trusted configurations and script sources. Inputs and outputs may be persisted in Terraform state; avoid putting secrets in them or in logged output unless their exposure is understood. Provider `log_output` sends lines with recognized level prefixes to Terraform logs. Optional plan commands execute during planning.

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
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `os_commands` (Attributes Map) Commands to run during the Terraform lifecycle, keyed by `GOOS` value or `default`; `default` is required. (see [below for nested schema](#nestedatt--os_commands))

### Optional

- `environment` (Map of String) The environment variables to set when executing commands; to be combined with the OS environment and the provider environment.
- `inputs` (Dynamic) Inputs to be made available to the script; these can be accessed as JSON via the `TF_SCRIPT_INPUTS` environment variable.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `triggers` (Dynamic) Allows specifying values that trigger resource replacement when changed.
- `working_directory` (String) The working directory to use when executing the commands; this will default to the _Terraform_ working directory.

### Read-Only

- `output` (Dynamic) The JSON output from plan, create, read, or update as a structured type. The prior value is available to plan (when prior state exists), read, update, and delete commands as JSON via `TF_SCRIPT_STATE_OUTPUT`. A top-level `__meta` key is reserved for provider metadata and is removed from output.
- `output_drift` (Boolean) True when the read command reported `__meta.output_drift_detected = true`. This signals output drift; it does not by itself guarantee that an update command will repair the external object.

<a id="nestedatt--os_commands"></a>
### Nested Schema for `os_commands`

Required:

- `create` (Attributes) The create command configuration. (see [below for nested schema](#nestedatt--os_commands--create))
- `delete` (Attributes) The delete command configuration. (see [below for nested schema](#nestedatt--os_commands--delete))
- `read` (Attributes) The read command configuration. (see [below for nested schema](#nestedatt--os_commands--read))
- `update` (Attributes) The update command configuration. (see [below for nested schema](#nestedatt--os_commands--update))

Optional:

- `plan` (Attributes) Optional command run during planning. It must write valid JSON to `TF_SCRIPT_OUTPUT`; that JSON becomes the planned `output`. A string value of `"???"` anywhere in plan output is treated as unknown and re-evaluated during apply, so this exact string cannot be represented as a known string output from the plan command. (see [below for nested schema](#nestedatt--os_commands--plan))

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
