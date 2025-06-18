data "shell_script" "example" {
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
