package provider

import (
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
					Config: `resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "tf-script-test"
  }
  commands = {
    create = <<-EOF
      set -euo pipefail
      touch "/tmp/$${TARGET_FILE}"
      printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
    EOF
    read = <<-EOF
      set -euo pipefail
      if [[ -f "/tmp/$${TARGET_FILE}" ]]; then
        printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
      else
        printf '{"exists": false}' > "$${TF_SCRIPT_OUTPUT}"
      fi
    EOF
    update = <<-EOF
      set -euo pipefail
      touch "/tmp/$${TARGET_FILE}"
      printf '{"exists": true}' > "$${TF_SCRIPT_OUTPUT}"
    EOF
    delete = <<-EOF
      set -euo pipefail
      rm -f "/tmp/$${TARGET_FILE}"
    EOF
  }
}`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output").AtMapKey("exists"), knownvalue.Bool(true)),
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
					Config: `resource "shell_script" "test" {
  environment = {
    "TARGET_FILE" = "tf-script-test"
  }
  commands = {
    create = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      touch "$${path}"
      printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
    EOF
    read = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      if [[ -f "$${path}" ]]; then
        printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
      else
        printf '{"exists": false,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
      fi
    EOF
    update = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      touch "$${path}"
      printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
    EOF
    delete = <<-EOF
      set -euo pipefail
      rm -f "/tmp/$${TARGET_FILE}"
    EOF
  }
}`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output").AtMapKey("exists"), knownvalue.Bool(true)),
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output").AtMapKey("path"), knownvalue.StringExact("/tmp/tf-script-test")),
					},
				},
				{
					Config: `resource "shell_script" "test" {
  environment = {
    "OLD_TARGET_FILE" = "tf-script-test"
    "TARGET_FILE" = "tf-script-test-new"
  }
  commands = {
    create = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      touch "$${path}"
      printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
    EOF
    read = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      if [[ -f "$${path}" ]]; then
        printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
      else
        printf '{"exists": false,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
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
}`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output").AtMapKey("exists"), knownvalue.Bool(true)),
						statecheck.ExpectKnownValue("shell_script.test", tfjsonpath.New("output").AtMapKey("path"), knownvalue.StringExact("/tmp/tf-script-test-new")),
					},
				},
			},
		})
	})
}
