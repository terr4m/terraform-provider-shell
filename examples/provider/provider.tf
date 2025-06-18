provider "shell" {}

data "shell_script" "example" {
  environment = {
    "TARGET" = "my-resource"
  }

  os_commands = {
    default = {
      read = {
        command = file("${path.module}/scripts/read.sh")
      }
    }
    windows = {
      read = {
        command = file("${path.module}/scripts/read.ps1")
      }
    }
  }
}

resource "shell_script" "example" {
  environment = {
    "TARGET" = "my-resource"
  }

  os_commands = {
    default = {
      create = {
        command = file("${path.module}/scripts/create.sh")
      }
      read = {
        command = file("${path.module}/scripts/read.sh")
      }
      update = {
        command = file("${path.module}/scripts/update.sh")
      }
      delete = {
        command = file("${path.module}/scripts/delete.sh")
      }
    }
    windows = {
      create = {
        command = file("${path.module}/scripts/create.ps1")
      }
      read = {
        command = file("${path.module}/scripts/read.ps1")
      }
      update = {
        command = file("${path.module}/scripts/update.ps1")
      }
      delete = {
        command = file("${path.module}/scripts/delete.ps1")
      }
    }
  }
}
