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
  // Border0 API URL. Optional and Defaults to https://api.border0.io/api/v1
  // if not set explicitly, the provider will use the environment variable BORDER0_API
  api_url = "https://api.border0.io/api/v1"
}
