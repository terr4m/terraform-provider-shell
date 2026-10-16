resource "shell_sensitive_script" "create_key" {
  os_commands = {
    default = {
      create = {
        command = <<-EOF
          set -euo pipefail
          key="$(openssl rand -hex 32)"
          printf '{"key":"%s"}' "$${key}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      read = {
        command = <<-EOF
          set -euo pipefail
          printf '%s' "$${TF_SCRIPT_STATE_OUTPUT}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      update = {
        command = <<-EOF
          set -euo pipefail
          printf '%s' "$${TF_SCRIPT_STATE_OUTPUT}" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
      delete = {
        command = ""
      }
    }
  }
}

output "key" {
  value     = shell_sensitive_script.create_key.output.key
  sensitive = true
}
