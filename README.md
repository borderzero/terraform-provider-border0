# terraform-provider-border0

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
      version = "0.1.4"
    }
  }
}

provider "border0" {
  // Border0 access token. Required.
  // if not set explicitly, the provider will use the environment variable BORDER0_TOKEN
  token = "_my_access_token_"
  // Border0 API URL. Optional and Defaults to https://api.border0.io/api/v1
  // if not set explicitly, the provider will use the environment variable BORDER0_API
  api_url = "https://api.border0.io/api/v1"
}
```

And explore the [examples](./examples) folder for additional use cases and examples.

## Terraform docs

To (re)generate docs, install `tfplugindocs` first:

```shell
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
```

And then run:

```shell
make docs
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
