terraform {
  required_providers {
    border0 = {
      source  = "border0.com/border0/border0"
      version = "0.1.0" # keep this as 0.1.0, it's a dummy version for local build
    }
  }
}

provider "border0" {
  api_url             = "https://api.staging.border0.com/api/v1"
  token               = "REPLACE ME WITH TOKEN"
  http_client_timeout = "5m"
}

