# terraform-provider-border0

[![Run tests](https://github.com/borderzero/terraform-provider-border0/actions/workflows/run_tests.yml/badge.svg)](https://github.com/borderzero/terraform-provider-border0/actions/workflows/run_tests.yml)
[![Terraform Registery](https://img.shields.io/badge/terraform-border0-4931de.svg?logo=terraform)](https://registry.terraform.io/providers/borderzero/border0/latest)
[![License](https://img.shields.io/github/license/borderzero/border0-go)](https://github.com/borderzero/border0-go/blob/master/LICENSE)
[![Slack](https://img.shields.io/badge/slack-community-orange.svg?logo=slack)](https://join.slack.com/t/border0community/shared_invite/zt-zbx586ls-x44z7I3POLPQfesRWnig7Q)

In this repo, you'll find the source code for the Border0 Terraform Provider. With this provider,
you can provision and manage Border0 resources, such as sockets, policies, and connectors,
using Terraform.

## Quickstart

Copy the following code to your terraform config and run `terraform init`:

```terraform
terraform {
  required_providers {
    border0 = {
      source  = "borderzero/border0"
    }
  }
}

provider "border0" {
  // Border0 token. Required.
  // If not set explicitly, the provider will use the env var BORDER0_TOKEN.
  //
  // You can generate a Border0 token by going to:
  // portal.border0.com -> Team -> Service Accounts
  // and creating a service account and a token.
  token = "_my_service_account_token_"
}
```

And explore the [examples](./examples) folder for additional use cases and examples. For a comprehensive step-by-step guide,
check out our [Quickstart Guide on docs.border0.com](https://docs.border0.com/docs/manage-border0-resources-with-terraform).

## Terraform docs generation

To (re)generate docs, install `tfplugindocs` first:

```shell
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
```

And then run:

```shell
make docs
```

## Local development

Make the changes you'd like!

Then, navigate to the [examples/development](./examples/development) directory:

```shell
cd examples/development
```

Install a border0 terraform provider (which includes your local changes):

```shell
make init
```

Log in to border0

```shell
make token
```

Apply the development setup:

```shell
make apply
```

Destroy the development setup:

```shell
make destroy
```

## Local release build

Install goreleaser:

```shell
brew install goreleaser
```

Run local release command and build terraform provider binaries

```shell
make release
```

You will find the releases in the `/dist` directory. You will need to rename the provider binary to `terraform-provider-border0` and move the binary into
[the appropriate subdirectory within the user plugins directory](https://learn.hashicorp.com/tutorials/terraform/provider-use?in=terraform/providers#install-hashicups-provider).

## Test configuration examples

Configuration examples can be tested with a local build of this terraform provider.

First, build and install the provider.

```shell
make install
```

Then, navigate to the [examples](./examples) directory.

```shell
cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```

Some different variations:

```shell
# run against prod api
BORDER0_TOKEN=_border0_auth_token_ terraform apply

# run with a local dev api that's listening to localhost:8080
BORDER0_TOKEN=_border0_auth_token_ BORDER0_API=http://localhost:8080/api/v1 terraform apply
```

## Tag a new release

Create a new tag. Assuming the next release is `v1.1.1`.

```shell
git tag v1.1.1
```

Push the tag to GitHub.

```shell
git push origin v1.1.1
```

It will trigger the GitHub Actions `release` workflow.
