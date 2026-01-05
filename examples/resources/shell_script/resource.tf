resource "shell_script" "example" {
  inputs = {
    file_name = "foo"
  }

  os_commands = {
    default = {
      plan = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          printf '{"exists": true,"path":"%%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      create = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          touch "$${path}"
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          if [[ -f "$${path}" ]]; then
            printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          else
            printf '{"exists": false,"path":"%s","__meta":{"output_drift_detected":true}}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
          fi
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          old_path="$(jq --raw-output '.path' <<<"$${TF_SCRIPT_STATE_OUTPUT}")"
          if [[ "$${path}" != "$${old_path}" ]] && [[ -f "$${old_path}" ]]; then
            mv -f "$${old_path}" "$${path}"
          else
            touch "$${path}"
          fi
          printf '{"exists": true,"path":"%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = <<-EOF
          set -euo pipefail
          file_name="$(jq --raw-output '.file_name' <<<"$${TF_SCRIPT_INPUTS}")"
          path="/tmp/$${file_name}"
          rm -f "$${path}"
        EOF
      }
    }
    windows = {
      plan = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      create = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          New-Item -Path $path -ItemType File -Force
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      read = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          if (Test-Path $path) {
            @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          } else {
            @{exists=$false; path=$path; __meta=@{output_drift_detected=$true}} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
          }
        EOF
      }
      update = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          $state = $env:TF_SCRIPT_STATE_OUTPUT | ConvertFrom-Json
          $oldPath = $state.path
          if ($path -ne $oldPath) {
            Move-Item -Path $oldPath -Destination $path -Force
          } else {
            New-Item -Path $path -ItemType File -Force
          }
          @{exists=$true; path=$path} | ConvertTo-Json -Compress | Out-File -FilePath $env:TF_SCRIPT_OUTPUT -Encoding utf8
        EOF
      }
      delete = {
        command = <<-EOF
          $inputs = $env:TF_SCRIPT_INPUTS | ConvertFrom-Json
          $fileName = $inputs.file_name
          $path = "$env:TEMP\$fileName"
          Remove-Item -Path $path -Force -ErrorAction Ignore
        EOF
      }
    }
  }

  output_drift = false
}
