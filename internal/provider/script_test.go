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

		file := acctest.RandomWithPrefix("tf-script-test")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "%s"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          touch "/tmp/$${TARGET_FILE}"
          printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          if [[ -f "/tmp/$${TARGET_FILE}" ]]; then
            printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false}' > "$${TF_SCRIPT_OUTPUT}"
          fi
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          touch "/tmp/$${TARGET_FILE}"
          printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
          rm -f "/tmp/$${TARGET_FILE}"
        EOF
      }
    }
  }
}
`, file),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("create_os", func(t *testing.T) {
		t.Parallel()

		file := acctest.RandomWithPrefix("tf-script-test")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "%s"
  }
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
          touch "/tmp/$${TARGET_FILE}"
          printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          if [[ -f "/tmp/$${TARGET_FILE}" ]]; then
            printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false}' > "$${TF_SCRIPT_OUTPUT}"
          fi
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          touch "/tmp/$${TARGET_FILE}"
          printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
          rm -f "/tmp/$${TARGET_FILE}"
        EOF
      }
    }
    windows = {
      create = {
        command = <<-EOF
          New-Item -Path "$env:TEMP\$env:TARGET_FILE" -ItemType File -Force
          '{"exists": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          if (Test-Path "$env:TEMP\$env:TARGET_FILE") {
            '{"exists": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          } else {
            '{"exists": false}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          }
        EOF
      }
      update = {
        command = <<-EOF
          New-Item -Path "$env:TEMP\$env:TARGET_FILE" -ItemType File -Force
          '{"exists": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = <<-EOF
          Remove-Item -Path "$env:TEMP\$env:TARGET_FILE" -Force -ErrorAction SilentlyContinue
        EOF
      }
    }
  }
}
`, file),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("create_update", func(t *testing.T) {
		t.Parallel()

		file := acctest.RandomWithPrefix("tf-script-test")
		newFile := acctest.RandomWithPrefix("tf-script-test-new")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "%s"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
          touch "$${path}"
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
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
          path="/tmp/$${TARGET_FILE}"
          old_path="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
          if [[ "$${path}" != "$${old_path}" ]] && [[ -f "$${old_path}" ]]; then
            rm -f "$${old_path}"
          fi
          touch "$${path}"
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
					path="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
          rm -f "$${path}"
        EOF
      }
    }
    windows = {
      create = {
        command = <<-EOF
          $path = "$env:TEMP\$env:TARGET_FILE"
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $path = "$env:TEMP\$env:TARGET_FILE"
          if (Test-Path $path) {
            @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          } else {
            @{exists=$false; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          }
        EOF
      }
      update = {
        command = <<-EOF
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $path = "$env:TEMP\$env:TARGET_FILE"
          $oldPath = $state.path
          if ($path -ne $oldPath) {
            Remove-Item -Path $oldPath -Force
          }
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = <<-EOF
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $path = $state.path
          Remove-Item -Path $path -Force
        EOF
      }
    }
  }
}
`, file),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringRegexp(regexp.MustCompile(".+tf-script-test")), "exists": knownvalue.Bool(true)})),
					},
				},
				{
					Config: fmt.Sprintf(`
resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "%s"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
          touch "$${path}"
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
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
          path="/tmp/$${TARGET_FILE}"
          old_path="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
          if [[ "$${path}" != "$${old_path}" ]] && [[ -f "$${old_path}" ]]; then
            rm -f "$${old_path}"
          fi
          touch "$${path}"
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
					path="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
          rm -f "$${path}"
        EOF
      }
    }
    windows = {
      create = {
        command = <<-EOF
          $path = "$env:TEMP\$env:TARGET_FILE"
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $path = "$env:TEMP\$env:TARGET_FILE"
          if (Test-Path $path) {
            @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          } else {
            @{exists=$false; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          }
        EOF
      }
      update = {
        command = <<-EOF
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $path = "$env:TEMP\$env:TARGET_FILE"
          $oldPath = $state.path
          if ($path -ne $oldPath) {
            Remove-Item -Path $oldPath -Force
          }
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = <<-EOF
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $path = $state.path
          Remove-Item -Path $path -Force
        EOF
      }
    }
  }
}
`, newFile),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringRegexp(regexp.MustCompile(".+tf-script-test-new")), "exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("log_output", func(t *testing.T) {
		t.Parallel()

		file := acctest.RandomWithPrefix("tf-script-test")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "%s"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          touch "/tmp/$${TARGET_FILE}"
          printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
          echo "[INFO] create"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          if [[ -f "/tmp/$${TARGET_FILE}" ]]; then
            printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false}' > "$${TF_SCRIPT_OUTPUT}"
          fi
          echo "[INFO] read"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          touch "/tmp/$${TARGET_FILE}"
          printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
          echo "[INFO] update"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
          rm -f "/tmp/$${TARGET_FILE}"
          echo "[INFO] delete"
        EOF
      }
    }
    windows = {
      create = {
        command = <<-EOF
          New-Item -Path "$env:TEMP\$env:TARGET_FILE" -ItemType File -Force
          '{"exists": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          Write-Output "[INFO] create"
        EOF
      }
      read = {
        command = <<-EOF
          if (Test-Path "$env:TEMP\$env:TARGET_FILE") {
            '{"exists": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          } else {
            '{"exists": false}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          }
          Write-Output "[INFO] read"
        EOF
      }
      update = {
        command = <<-EOF
          New-Item -Path "$env:TEMP\$env:TARGET_FILE" -ItemType File -Force
          '{"exists": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          Write-Output "[INFO] update"
        EOF
      }
      delete = {
        command = <<-EOF
          Remove-Item -Path "$env:TEMP\$env:TARGET_FILE" -Force -ErrorAction SilentlyContinue
          Write-Output "[INFO] delete"
        EOF
      }
    }
  }
}
`, file),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"exists": knownvalue.Bool(true)})),
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
