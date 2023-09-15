terraform {
  required_providers {
    border0 = {
      source = "borderzero/border0"
    }
  }
}

provider "border0" {
  // Border0 access token. Required.
  // If not set explicitly, the provider will use the env var BORDER0_TOKEN.
  //
  // You can generate a Border0 access token one by going to:
  // portal.border0.com -> Organization Settings -> Access Tokens
  // and then create a token in Admin permission groups.
  token = "_my_access_token_"
}
