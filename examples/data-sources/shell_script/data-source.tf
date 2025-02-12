data "shell_script" "example" {
  environment = {
    "TF_VERSION_COUNT" = "3"
  }

  command = <<-EOF
    set -euo pipefail
    curl -s https://endoflife.date/api/terraform.json | jq -rc --argjson count "$${TF_VERSION_COUNT}" '[sort_by(.releaseDate) | reverse | .[0:$count] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
  EOF
}
