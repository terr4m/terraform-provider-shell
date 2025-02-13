package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccScriptDataSource(t *testing.T) {
	t.Run("read", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `data "shell_script" "test" {
  command = <<-EOF
    set -euo pipefail
    curl -s https://endoflife.date/api/terraform.json | jq -rc '[sort_by(.releaseDate) | reverse | .[0:4] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
  EOF
}`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ListSizeExact(4)),
					},
				},
			},
		})
	})

	t.Run("read_with_ignored_output", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `data "shell_script" "test" {
  command = <<-EOF
    set -euo pipefail
    echo "Test..."
    echo "Test..." >&2
    curl -s https://endoflife.date/api/terraform.json | jq -rc '[sort_by(.releaseDate) | reverse | .[0:4] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
  EOF
}`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ListSizeExact(4)),
					},
				},
			},
		})
	})

	t.Run("read_with_interpreter", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `data "shell_script" "test" {
  interpreter = ["/bin/sh", "-c"]
  command = <<-EOF
    set -eu
    curl -s https://endoflife.date/api/terraform.json | jq -rc '[sort_by(.releaseDate) | reverse | .[0:4] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
  EOF
}`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ListSizeExact(4)),
					},
				},
			},
		})
	})

	t.Run("read_with_environment", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `data "shell_script" "test" {
  environment = {
    "TF_VERSION_COUNT" = "3"
  }
  command = <<-EOF
    set -euo pipefail
    curl -s https://endoflife.date/api/terraform.json | jq -rc --argjson count "$${TF_VERSION_COUNT}" '[sort_by(.releaseDate) | reverse | .[0:$count] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
  EOF
}`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.NotNull()),
						statecheck.ExpectKnownValue("data.shell_script.test", tfjsonpath.New("output"), knownvalue.ListSizeExact(3)),
					},
				},
			},
		})
	})
}
