resource "shell_script" "example" {
  environment = {
    "OLD_TARGET_FILE" = "my-resource"
    "TARGET_FILE"     = "my-resource-new"
  }

  commands = {
    create = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      touch "$${path}"
      printf '{"exists": true, "path": "%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
    EOF

    read = <<-EOF
      set -euo pipefail
      path="/tmp/$${TARGET_FILE}"
      if [[ -f "$${path}" ]]; then
        printf '{"exists": true, "path": "%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
      else
        printf '{"exists": false, "path": "%s"}' "$${path}" > "$${TF_SCRIPT_OUTPUT}"
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
}
