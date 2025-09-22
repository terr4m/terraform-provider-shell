resource "shell_script" "example" {
  environment = {
    "TARGET_FILE" = "foo"
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
          old_path="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
          if [[ "$${path}" != "$${old_path}" ]] && [[ -f "$${old_path}" ]]; then
            rm -f "$${old_path}"
          fi
          touch "$${path}"
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
					path="$(echo "$${TF_SCRIPT_STATE_OUTPUT}" | jq -r '.path')"
          rm -f "/tmp/$${path}"
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
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $path = "$env:TEMP\$env:TARGET_FILE"
          $oldPath = $state.path
          if ($path -ne $oldPath) {
            Remove-Item -Path $oldPath -Force
          }
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = <<-EOF
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $path = $state.path
          Remove-Item -Path $file -Force
        EOF
      }
    }
  }
}
