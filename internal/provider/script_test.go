package provider

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccScriptResource(t *testing.T) {
	t.Parallel()

	t.Run("create_default", func(t *testing.T) {
		t.Parallel()

		if runtime.GOOS == "windows" {
			t.Skip("skipping test on windows platform")
		}

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
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"run": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("create_os", func(t *testing.T) {
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
          value="$(echo "$${TF_SCRIPT_INPUTS}" | jq -r '.value')"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          value="$(echo "$${TF_SCRIPT_INPUTS}" | jq -r '.value')"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          value="$(echo "$${TF_SCRIPT_INPUTS}" | jq -r '.value')"
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
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"value": knownvalue.StringExact("my-value")})),
					},
				},
			},
		})
	})

	t.Run("create_update", func(t *testing.T) {
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
          file_name="$(echo "$${TF_SCRIPT_INPUTS}" | jq -r '.file_name')"
          path="/tmp/$${file_name}"
          touch "$${path}"
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(echo "$${TF_SCRIPT_INPUTS}" | jq -r '.file_name')"
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
          file_name="$(echo "$${TF_SCRIPT_INPUTS}" | jq -r '.file_name')"
          path="/tmp/$${file_name}"
          old_path="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
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
          file_name="$(echo "$${TF_SCRIPT_INPUTS}" | jq -r '.file_name')"
          path="/tmp/$${file_name}"
          rm -f "$${path}"
        EOF
      }
    }
    windows = {
      create = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $path = "$env:TEMP\$inputs.file_name"
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $path = "$env:TEMP\$inputs.file_name"
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
          $path = "$env:TEMP\$inputs.file_name"
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
          $path = "$env:TEMP\$inputs.file_name"
          Remove-Item -Path $path -Force
        EOF
      }
    }
  }
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

	t.Run("log_output", func(t *testing.T) {
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
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"run": knownvalue.Bool(true)})),
					},
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
}
`,
					ExpectError: regexp.MustCompile(`my-error`),
				},
			},
		})
	})
}
