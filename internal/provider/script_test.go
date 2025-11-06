package provider

import (
	"fmt"
	"math/big"
	"os"
	"path"
	"regexp"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccScriptResource(t *testing.T) {
	t.Parallel()

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		cmd := `printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"`
		if runtime.GOOS == "windows" {
			cmd = `'{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8`
		}

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "shell_script" "test" {
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          %s
        EOF
      }
      read = {
        command = <<-EOF
          %[1]s
        EOF
      }
      update = {
        command = <<-EOF
          %[1]s
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`, cmd),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"run": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("create_with_os", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          exit 1
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          exit 1
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          exit 1
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
          exit 1
        EOF
      }
    }
    linux = {
      create = {
        command = <<-EOF
          set -euo pipefail
          printf '{"os": "linux"}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '{"os": "linux"}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          printf '{"os": "linux"}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
    windows = {
      create = {
        command = <<-EOF
          '{"os": "windows"}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          '{"os": "windows"}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      update = {
        command = <<-EOF
          '{"os": "windows"}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"os": knownvalue.StringExact(runtime.GOOS)})),
					},
				},
			},
		})
	})

	t.Run("create_with_environment", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  environment = {
    "MY_VALUE" = "my-value"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          printf '{"value": "%s"}' "$${MY_VALUE}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '{"value": "%s"}' "$${MY_VALUE}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          printf '{"value": "%s"}' "$${MY_VALUE}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
    windows = {
      create = {
        command = <<-EOF
          @{value=$env:MY_VALUE} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          @{value=$env:MY_VALUE} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      update = {
        command = <<-EOF
          @{value=$env:MY_VALUE} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"value": knownvalue.StringExact("my-value")})),
					},
				},
			},
		})
	})

	t.Run("create_with_inputs", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  inputs = {
    value = "my-value"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_INPUTS}")"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_INPUTS}")"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_INPUTS}")"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
    windows = {
      create = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          @{value=$inputs.value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          @{value=$inputs.value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      update = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          @{value=$inputs.value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"value": knownvalue.StringExact("my-value")})),
					},
				},
			},
		})
	})

	t.Run("create_with_triggers", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  triggers = {
    gen = 1
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
    windows = {
      create = {
        command = <<-EOF
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      update = {
        command = <<-EOF
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"run": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("create_with_plan", func(t *testing.T) {
		t.Parallel()

		v := "my-value"

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					ConfigVariables: map[string]config.Variable{
						"value": config.StringVariable(v),
					},
					Config: `
variable "value" {
  type = string
}
resource "shell_script" "test" {
  inputs = {
    value = var.value
  }
  os_commands = {
    default = {
      plan = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_INPUTS}")"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      create = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_INPUTS}")"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_INPUTS}")"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_INPUTS}")"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
    windows = {
      plan = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          @{value=$inputs.value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      create = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          @{value=$inputs.value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          @{value=$inputs.value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      update = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          @{value=$inputs.value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`,

					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"value": knownvalue.StringExact(v)})),
						},
					},
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"value": knownvalue.StringExact(v)})),
					},
				},
			},
		})
	})

	t.Run("create_with_plan_unknown", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  os_commands = {
    default = {
      plan = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_STATE_OUTPUT:-}")"
          if [[ -z "$${value}" ]]; then
            printf '{"value": "???"}' > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"value": %d}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
          fi
        EOF
      }
      create = {
        command = <<-EOF
          set -euo pipefail
          printf '{"value": %d}' "$(date +%s)" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_STATE_OUTPUT}")"
          printf '{"value": %d}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          printf '{"value": %d}' "$(date +%s)" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
    windows = {
      plan = {
        command = <<-EOF
          if ($env:TF_SCRIPT_STATE_OUTPUT) {
            $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
            $value = $state.value
            @{value=$value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          } else {
            @{value="???"} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          }
        EOF
      }
      create = {
        command = <<-EOF
          $value = [int64] (Get-Date -UFormat "%s")
          @{value=$value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $value = $state.value
          @{value=$value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      update = {
        command = <<-EOF
          $value = [int64] (Get-Date -UFormat "%s")
          @{value=$value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`,

					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectUnknownValue("shell_script.test", tfjsonpath.New("output").AtMapKey("value")),
						},
					},
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"value": knownvalue.NumberFunc(func(v *big.Float) error { return nil })})),
					},
				},
			},
		})
	})

	t.Run("create_with_log_output", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          echo "[INFO] create"
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          echo "[INFO] read"
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          echo "[INFO] update"
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
          echo "[INFO] delete"
        EOF
      }
    }
    windows = {
      create = {
        command = <<-EOF
          Write-Output "[INFO] create"
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          Write-Output "[INFO] read"
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      update = {
        command = <<-EOF
          Write-Output "[INFO] update"
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = <<-EOF
          Write-Output "[INFO] delete"
        EOF
      }
    }
  }
  output_drift = false
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"run": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Parallel()

		file := acctest.RandomWithPrefix("tf-script-test")
		newFile := acctest.RandomWithPrefix("tf-script-test-new")
		config := `
resource "shell_script" "test" {
  inputs = {
    file_name = "%s"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          touch "$${path}"
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          if [[ -f "$${path}" ]]; then
            printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
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
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
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
            @{exists=$false; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
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
          Remove-Item -Path $path -Force
        EOF
      }
    }
  }
  output_drift = false
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(config, file),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringRegexp(regexp.MustCompile(".+tf-script-test")), "exists": knownvalue.Bool(true)})),
					},
				},
				{
					Config: fmt.Sprintf(config, newFile),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringRegexp(regexp.MustCompile(".+tf-script-test-new")), "exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("update_external", func(t *testing.T) {
		t.Parallel()

		file := acctest.RandomWithPrefix("tf-script-test")
		config := fmt.Sprintf(`
resource "shell_script" "test" {
  inputs = {
    file_name = "%s"
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
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          if [[ -f "$${path}" ]]; then
            printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false,"path":"%%s","__meta":{"output_drift_detected":true}}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
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
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
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
`, file)

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: config,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringRegexp(regexp.MustCompile(fmt.Sprintf(".+%s$", file))), "exists": knownvalue.Bool(true)})),
					},
				},
				{
					PreConfig: func() {
						if err := os.Remove(path.Join(os.TempDir(), file)); err != nil {
							t.Fatalf("failed to remove file: %s", err)
						}
					},
					Config: config,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectNonEmptyPlan(),
						},
					},
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringRegexp(regexp.MustCompile(fmt.Sprintf(".+%s$", file))), "exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("update_with_triggers", func(t *testing.T) {
		t.Parallel()

		config := `
resource "shell_script" "test" {
  triggers = {
    gen = %d
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
    windows = {
      create = {
        command = <<-EOF
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      update = {
        command = <<-EOF
          '{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(config, 0),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"run": knownvalue.Bool(true)})),
					},
				},
				{
					Config: fmt.Sprintf(config, 1),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectResourceAction("shell_script.test", plancheck.ResourceActionReplace),
						},
					},
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"run": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("error_no_output_drift", func(t *testing.T) {
		t.Parallel()

		cmd := `printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"`
		if runtime.GOOS == "windows" {
			cmd = `'{"run": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8`
		}

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "shell_script" "test" {
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          %s
        EOF
      }
      read = {
        command = <<-EOF
          %[1]s
        EOF
      }
      update = {
        command = <<-EOF
          %[1]s
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
}
`, cmd),
					ExpectError: regexp.MustCompile(`Missing required argument`),
				},
			},
		})
	})

	t.Run("error_no_default_commands", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  os_commands = {
    linux = {
      create = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`,
					ExpectError: regexp.MustCompile(`Default commands are required`),
				},
			},
		})
	})

	t.Run("error_no_json", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  os_commands = {
    default = {
      create = {
        command = ""
      }
      read = {
        command = ""
      }
      update = {
        command = ""
      }
      delete = {
        command = ""
      }
    }
  }
  output_drift = false
}
`,
					ExpectError: regexp.MustCompile(`Failed to read output file`),
				},
			},
		})
	})

	t.Run("error_exit_code", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  os_commands = {
    default = {
      create = {
        command = "exit 1"
      }
      read = {
        command = "exit 1"
      }
      update = {
        command = "exit 1"
      }
      delete = {
        command = "exit 1"
      }
    }
  }
  output_drift = false
}
`,
					ExpectError: regexp.MustCompile(`Command failed with exit code: 1`),
				},
			},
		})
	})

	t.Run("error_message", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          printf 'my-error' > "$${TF_SCRIPT_ERROR}"
          exit 1
        EOF
      }
      read = {
        command = <<-EOF
          printf 'my-error' > "$${TF_SCRIPT_ERROR}"
          exit 1
        EOF
      }
      update = {
        command = <<-EOF
          printf 'my-error' > "$${TF_SCRIPT_ERROR}"
          exit 1
        EOF
      }
      delete = {
        command = <<-EOF
          printf 'my-error' > "$${TF_SCRIPT_ERROR}"
          exit 1
        EOF
      }
    }
    windows = {
      create = {
        command = <<-EOF
          'my-error' | Out-File -FilePath $env:TF_SCRIPT_ERROR -Encoding utf8
          exit 1
        EOF
      }
      read = {
        command = <<-EOF
          'my-error' | Out-File -FilePath $env:TF_SCRIPT_ERROR -Encoding utf8
          exit 1
        EOF
      }
      update = {
        command = <<-EOF
          'my-error' | Out-File -FilePath $env:TF_SCRIPT_ERROR -Encoding utf8
          exit 1
        EOF
      }
      delete = {
        command = <<-EOF
          'my-error' | Out-File -FilePath $env:TF_SCRIPT_ERROR -Encoding utf8
          exit 1
        EOF
      }
    }
  }
  output_drift = false
}
`,
					ExpectError: regexp.MustCompile(`my-error`),
				},
			},
		})
	})
}
