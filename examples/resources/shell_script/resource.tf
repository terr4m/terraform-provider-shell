resource "shell_script" "example" {
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
