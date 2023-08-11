terraform {
  required_providers {
    border0 = {
      source  = "borderzero/border0"
      version = "<version>"
    }
  }
}

provider "border0" {
  // Border0 access token. Required.
  // if not set explicitly, the provider will use the environment variable BORDER0_TOKEN
  token = "_my_access_token_"
}
