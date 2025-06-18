package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccScriptResource(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "tf-script-test"
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
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("create_os_script", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "tf-script-test"
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
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("create_update", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "tf-script-test"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
          touch "$${path}"
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
          if [[ -f "$${path}" ]]; then
            printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          fi
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
          touch "$${path}"
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
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
          $path = "$env:TEMP\$env:TARGET_FILE"
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
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
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringExact("/tmp/tf-script-test"), "exists": knownvalue.Bool(true)})),
					},
				},
				{
					Config: `
resource "shell_script" "test" {
  environment = {
    "OLD_TARGET_FILE" = "tf-script-test"
    "TARGET_FILE" = "tf-script-test-new"
  }
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
          touch "$${path}"
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          path="/tmp/$${TARGET_FILE}"
          if [[ -f "$${path}" ]]; then
            printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          fi
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          old_path="/tmp/$${OLD_TARGET_FILE}"
          if [[ -f "$${old_path}" ]]; then
            rm -f "$${old_path}"
          fi
          path="/tmp/$${TARGET_FILE}"
          touch "$${path}"
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
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
          $oldPath = "$env:TEMP\$env:OLD_TARGET_FILE"
          if (Test-Path $oldPath) {
            Remove-Item -Path $oldPath -Force
          }
          $path = "$env:TEMP\$env:TARGET_FILE"
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
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
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringRegexp(regexp.MustCompile(".+")), "exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("create_with_state", func(t *testing.T) {
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
          file="$(mktemp)"
          touch "$${file}"
          printf '{"path": "%s","exists": true}' "$${file}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          file="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
          if [[ -f "$${file}" ]]; then
            printf '{"path": "%s","exists": true}' "$${file}" > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"path": "%s","exists": false}' "$${file}" > "$${TF_SCRIPT_OUTPUT}"
          fi
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          printf '%s' "$${TF_SCRIPT_STATE_OUTPUT}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
          file="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
          rm -f "$${file}"
        EOF
      }
    }
    windows = {
      create = {
        command = <<-EOF
          $file = [System.IO.Path]::GetTempFileName()
          New-Item -Path $file -ItemType File -Force
          @{path=$file; exists=$true} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $file = $state.path
          if (Test-Path $file) {
            @{path=$file; exists=$true} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          } else {
            @{path=$file; exists=$false} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          }
        EOF
      }
      update = {
        command = <<-EOF
          $env:TF_SCRIPT_STATE_OUTPUT | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = <<-EOF
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $file = $state.path
          Remove-Item -Path $file -Force -ErrorAction SilentlyContinue
        EOF
      }
    }
  }
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"path": knownvalue.StringRegexp(regexp.MustCompile(".+")), "exists": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})
}
