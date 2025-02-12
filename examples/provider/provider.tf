provider "shell" {
  interpreter = ["/bin/bash", "-c"]
}

data "shell_script" "example" {
  environment = {
    "TARGET" = "my-resource"
  }

  command = file("${path.module}/scripts/read.sh")
}

resource "shell_script" "example" {
  environment = {
    "TARGET" = "my-resource"
  }

  commands = {
    create = file("${path.module}/scripts/create.sh")
    read   = file("${path.module}/scripts/read.sh")
    update = file("${path.module}/scripts/update.sh")
    delete = file("${path.module}/scripts/delete.sh")
  }
}
