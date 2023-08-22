// use the provider from local path (for development)
terraform {
  required_providers {
    border0 = {
      source  = "border0.com/border0/border0"
      version = "0.1.0" # keep this as 0.1.0, it's a dummy version for local build
    }
  }
}

variable "token" { type = string }
variable "api_url" { type = string }

provider "border0" {
  token   = var.token
  api_url = var.api_url
}

resource "border0_connector" "test_tf_connector" {
  name                         = "test-tf-connector-1"
  description                  = "test connector from terraform"
  built_in_ssh_service_enabled = true
}

resource "border0_connector_token" "test_tf_connector_token_never_expires" {
  connector_id = border0_connector.test_tf_connector.id
  name         = "test-tf-connector-token-never-expires"
}

resource "border0_connector_token" "test_tf_connector_token_expires" {
  connector_id = border0_connector.test_tf_connector.id
  name         = "test-tf-connector-token-never-expires"
  expires_at   = "2023-12-31T23:59:59Z"
}

resource "border0_policy" "test_tf_policy" {
  name        = "test-tf-policy"
  description = "test policy from terraform"
  policy_data = jsonencode({
    "version" : "v1",
    "action" : ["database", "ssh", "http", "tls"],
    "condition" : {
      "who" : {
        "email" : ["johndoe@example.com"],
        "domain" : ["example.com"]
      },
      "where" : {
        "allowed_ip" : ["0.0.0.0/0", "::/0"],
        "country" : ["NL", "CA", "US", "BR", "FR"],
        "country_not" : ["BE"]
      },
      "when" : {
        "after" : "2022-10-13T05:12:26Z",
        "before" : null,
        "time_of_day_after" : "00:00 UTC",
        "time_of_day_before" : "23:59 UTC"
      }
    }
  })
}

data "border0_policy_document" "test_tf_policy_document" {
  version = "v1"
  action  = ["database", "ssh", "http", "tls"]
  condition {
    who {
      email  = ["johndoe@example.com"]
      domain = ["example.com"]
    }
    where {
      allowed_ip  = ["0.0.0.0/0", "::/0"]
      country     = ["NL", "CA", "US", "BR", "FR"]
      country_not = ["BE"]
    }
    when {
      after              = "2022-10-13T05:12:27Z"
      time_of_day_after  = "00:00 UTC"
      time_of_day_before = "23:59 UTC"
    }
  }
}

resource "border0_policy" "another_test_tf_policy" {
  name        = "another-test-tf-policy"
  description = "another test policy from terraform"
  policy_data = data.border0_policy_document.test_tf_policy_document.json
}

resource "border0_socket" "test_tf_http" {
  name        = "test-tf-http"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
  upstream_type = "https"
}

resource "border0_socket" "test_tf_ssh" {
  name                         = "test-tf-ssh"
  recording_enabled            = true
  socket_type                  = "ssh"
  connector_id                 = border0_connector.test_tf_connector.id
  upstream_hostname            = "127.0.0.1"
  upstream_port                = 22
  upstream_username            = "test_user"
  upstream_authentication_type = "border0_certificate"
}

resource "border0_policy_attachment" "test_tf_policy_http_socket" {
  policy_id = border0_policy.test_tf_policy.id
  socket_id = border0_socket.test_tf_http.id
}

resource "border0_policy_attachment" "test_tf_policy_ssh_socket" {
  policy_id = border0_policy.test_tf_policy.id
  socket_id = border0_socket.test_tf_ssh.id
}

resource "border0_policy_attachment" "another_test_tf_policy_ssh_socket" {
  policy_id = border0_policy.another_test_tf_policy.id
  socket_id = border0_socket.test_tf_ssh.id
}

output "managed_resources" {
  value = {
    connector = {
      id   = border0_connector.test_tf_connector.id
      name = border0_connector.test_tf_connector.name
    }
    connector_token = {
      id           = border0_connector_token.test_tf_connector_token_never_expires.id
      connector_id = border0_connector_token.test_tf_connector_token_never_expires.connector_id
      name         = border0_connector_token.test_tf_connector_token_never_expires.name
      expires_at   = border0_connector_token.test_tf_connector_token_never_expires.expires_at
    }
    another_connector_token = {
      id           = border0_connector_token.test_tf_connector_token_expires.id
      connector_id = border0_connector_token.test_tf_connector_token_expires.connector_id
      name         = border0_connector_token.test_tf_connector_token_expires.name
      expires_at   = border0_connector_token.test_tf_connector_token_expires.expires_at
    }
    policy = {
      id   = border0_policy.test_tf_policy.id
      name = border0_policy.test_tf_policy.name
    }
    another_policy = {
      id   = border0_policy.another_test_tf_policy.id
      name = border0_policy.another_test_tf_policy.name
    }
    http = {
      id   = border0_socket.test_tf_http.id
      name = border0_socket.test_tf_http.name
      type = border0_socket.test_tf_http.socket_type
    }
    ssh = {
      id   = border0_socket.test_tf_ssh.id
      name = border0_socket.test_tf_ssh.name
      type = border0_socket.test_tf_ssh.socket_type
    }
  }
}
