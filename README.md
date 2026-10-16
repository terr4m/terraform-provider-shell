# Terraform Provider Shell

![Terraform](https://img.shields.io/badge/Terraform-terr4m/shell-purple?logo=terraform&link=https%3A%2F%2Fregistry.terraform.io%2Fproviders%2Fterr4m%2Fshell)
![GitHub Release (latest SemVer)](https://img.shields.io/github/v/release/terr4m/terraform-provider-shell?logo=github&label=Release&sort=semver)
![Release](https://github.com/terr4m/terraform-provider-shell/actions/workflows/release.yaml/badge.svg)
![Validate](https://github.com/terr4m/terraform-provider-shell/actions/workflows/validate.yaml/badge.svg?branch=main)

This _Terraform_ provider runs configured shell commands as part of a Terraform configuration. It provides a script resource and data source for integrating with systems that do not have a native provider. The provider is published on the [Terraform Registry](https://registry.terraform.io/providers/terr4m/shell/latest).

## Quick start

Declare the provider and a version constraint in your Terraform configuration:

```terraform
terraform {
  required_providers {
    shell = {
      source  = "terr4m/shell"
      version = "~> 1.0"
    }
  }
}

provider "shell" {}

data "shell_script" "example" {
  inputs = {
    greeting = "hello"
  }

  os_commands = {
    default = {
      read = {
        command = <<-EOF
          printf '{"greeting":"%s"}' "$(jq --raw-output '.greeting' <<<"$${TF_SCRIPT_INPUTS}")" > "$${TF_SCRIPT_OUTPUT}"
        EOF
      }
    }
  }
}
```

Run `terraform init` to install the provider. This example uses the default Bash interpreter and requires `jq` on the machine running Terraform. See the Registry docs for [provider configuration](https://registry.terraform.io/providers/terr4m/shell/latest/docs) and the complete [data source](https://registry.terraform.io/providers/terr4m/shell/latest/docs/data-sources/script) and [resource](https://registry.terraform.io/providers/terr4m/shell/latest/docs/resources/script) references.

> [!WARNING]
> This provider executes arbitrary configured commands on the machine running Terraform. Use only trusted scripts and configurations. Commands may run during planning, apply, refresh, and destroy. Inputs and outputs may be persisted in Terraform state; do not pass secrets through them unless you understand the state and logging exposure.

## Contributing

Contributions are welcome, including documentation, issues, and code. See [CONTRIBUTING](.github/CONTRIBUTING.md) for development setup, tests, generation, and linting instructions.
