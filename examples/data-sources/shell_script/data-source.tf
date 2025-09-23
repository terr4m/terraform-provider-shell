data "shell_script" "example" {
  inputs = {
    version_count = 3
  }
  os_commands = {
    default = {
      read = {
        command = <<-EOF
          set -euo pipefail
          version_count="$(echo "$${TF_SCRIPT_INPUTS}" | jq -r '.version_count')"
          curl -s https://endoflife.date/api/terraform.json | jq -rc --argjson count "$${version_count}" '[sort_by(.releaseDate) | reverse | .[0:$count] | .[].latest]' > "$${TF_SCRIPT_OUTPUT}"
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
