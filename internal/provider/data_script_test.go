package provider

import (
	"fmt"
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

		cmd := `printf '{"data": true}' > "$${TF_SCRIPT_OUTPUT}"`
		if runtime.GOOS == "windows" {
			cmd = `'{"data": true}' | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8`
		}

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(`
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          %s
        EOF
      }
    }
  }
}
`, cmd),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"data": knownvalue.Bool(true)})),
					},
				},
			},
		})
	})

	t.Run("read_os", func(t *testing.T) {
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

	t.Run("read_with_ignored_stdout", func(t *testing.T) {
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
          printf '{"os": "linux"}\n' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          Write-Host "Test..."
          Write-Error "Test..."
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
          printf '{"os": "linux"}\n' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        interpreter = ["pwsh", "-c"]
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
    "MY_VALUE" = "my-value"
  }
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '{"value": "%s"}' "$${MY_VALUE}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          @{value=$env:MY_VALUE} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
}
`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"value": knownvalue.StringExact("my-value")})),
					},
				},
			},
		})
	})

	t.Run("read_with_inputs", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  inputs = {
    value = "my-value"
  }
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          value="$(jq --raw-output '.value' <<<"$${TF_SCRIPT_INPUTS}")"
          printf '{"value": "%s"}' "$${value}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          @{value=$inputs.value} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
}
`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ObjectExact(map[string]knownvalue.Check{"value": knownvalue.StringExact("my-value")})),
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
    read = "30s"
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

	t.Run("error_no_default_commands", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "shell_script" "test" {
  os_commands = {
    linux = {
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '{"run": true}' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
  }
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
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
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
data "shell_script" "test" {
  os_commands = {
    default = {
      read = {
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

	t.Run("error_timeout", func(t *testing.T) {
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

	t.Run("error_message", func(t *testing.T) {
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
