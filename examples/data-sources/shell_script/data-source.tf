data "shell_script" "example" {
  inputs = {
    version_count = 3
  }
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          version_count="$(jq --raw-output '.version_count' <<<"$${TF_SCRIPT_INPUTS}")"
          curl --fail --silent --location --retry 3 https://endoflife.date/api/terraform.json | jq --raw-output --compact-output --argjson count "$${version_count}" '[sort_by(.releaseDate) | reverse | .[0:$count] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
    windows = {
      read = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $response = Invoke-RestMethod -Uri "https://endoflife.date/api/terraform.json"
          $sorted = $response | Sort-Object releaseDate -Descending | Select-Object -First $inputs.version_count
          $latest = $sorted | ForEach-Object { $_.latest }
          $latest | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
    }
  }
}
