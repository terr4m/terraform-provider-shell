package provider

import (
	"regexp"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccScriptDataSource(t *testing.T) {
	t.Parallel()

	t.Run("read", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          curl -s https://endoflife.date/api/terraform.json | jq -rc '[sort_by(.releaseDate) | reverse | .[0:4] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          $response = Invoke-RestMethod -Uri "https://endoflife.date/api/terraform.json"
          $sorted = $response | Sort-Object releaseDate -Descending | Select-Object -First 4
          $latest = $sorted | ForEach-Object { $_.latest }
          $latest | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ListSizeExact(4)),
					},
				},
			},
		})
	})

	t.Run("read_os_script", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          exit 1
        EOF
      }
    }
    linux = {
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '{"os": "linux"}\n' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          '{"os": "windows"}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"os": knownvalue.StringExact(runtime.GOOS)})),
					},
				},
			},
		})
	})

	t.Run("read_with_ignored_output", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          echo "Test..."
          echo "Test..." >&2
          curl -s https://endoflife.date/api/terraform.json | jq -rc '[sort_by(.releaseDate) | reverse | .[0:4] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          Write-Host "Test..."
					Write-Error "Test..."
          $response = Invoke-RestMethod -Uri "https://endoflife.date/api/terraform.json"
          $sorted = $response | Sort-Object releaseDate -Descending | Select-Object -First 4
          $latest = $sorted | ForEach-Object { $_.latest }
          $latest | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ListSizeExact(4)),
					},
				},
			},
		})
	})

	t.Run("read_with_interpreter", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
        interpreter = ["/bin/bash", "-c"]
        command = <<-EOF
          set -euo pipefail
          curl -s https://endoflife.date/api/terraform.json | jq -rc '[sort_by(.releaseDate) | reverse | .[0:4] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        interpreter = ["pwsh", "-c"]
        command = <<-EOF
          $response = Invoke-RestMethod -Uri "https://endoflife.date/api/terraform.json"
          $sorted = $response | Sort-Object releaseDate -Descending | Select-Object -First 4
          $latest = $sorted | ForEach-Object { $_.latest }
          $latest | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
}
`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ListSizeExact(4)),
					},
				},
			},
		})
	})

	t.Run("read_with_environment", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  environment = {
    "TF_VERSION_COUNT" = "3"
  }
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          curl -s https://endoflife.date/api/terraform.json | jq -rc --argjson count "$${TF_VERSION_COUNT}" '[sort_by(.releaseDate) | reverse | .[0:$count] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          $response = Invoke-RestMethod -Uri "https://endoflife.date/api/terraform.json"
          $sorted = $response | Sort-Object releaseDate -Descending | Select-Object -First $env:TF_VERSION_COUNT
          $latest = $sorted | ForEach-Object { $_.latest }
          $latest | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
}
`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ListSizeExact(3)),
					},
				},
			},
		})
	})

	t.Run("read_with_timeout", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          sleep 1s
          printf '{"success": true}\n' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          Start-Sleep -Seconds 1
          '{"success": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
  timeouts = {
    read = "10s"
  }
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"success": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("read_with_timeout_error", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          sleep 10s
          printf '{"success": true}\n' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          Start-Sleep -Seconds 10
          '{"success": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
  timeouts = {
    read = "1s"
  }
}
`,
					ExpectError: regexp.MustCompile(".+"),
				},
			},
		})
	})

	t.Run("read_with_error_message", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          printf 'my-error' > "$${TF_SCRIPT_ERROR}"
          exit 1
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          'my-error' | Out-File -FilePath $env:TF_SCRIPT_ERROR -Encoding utf8
          exit 1
        EOF
      }
    }
  }
}
`,
					ExpectError: regexp.MustCompile("my-error"),
				},
			},
		})
	})
}
