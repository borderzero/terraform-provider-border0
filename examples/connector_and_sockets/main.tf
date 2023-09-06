// use the provider from terraform registry
terraform {
  required_providers {
    border0 = {
      source  = "borderzero/border0"
      version = "0.1.7"
    }
  }
}

variable "token" {
  type = string
}

provider "border0" {
  token = var.token
}

resource "border0_connector" "test_tf_connector" {
  name                         = "test-tf-connector"
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

resource "border0_socket" "test_tf_http" {
  name        = "test-tf-http"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
  upstream_type = "https"
}

resource "border0_socket" "test_tf_ssh" {
  name              = "test-tf-ssh"
  recording_enabled = true
  socket_type       = "ssh"
  connector_id      = border0_connector.test_tf_connector.id

  ssh_configuration {
    hostname            = "127.0.0.1"
    port                = 22
    username            = "test_user"
    authentication_type = "border0_certificate"
  }
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
    http_socket = {
      id   = border0_socket.test_tf_http.id
      name = border0_socket.test_tf_http.name
      type = border0_socket.test_tf_http.socket_type
    }
    ssh_socket = {
      id   = border0_socket.test_tf_ssh.id
      name = border0_socket.test_tf_ssh.name
      type = border0_socket.test_tf_ssh.socket_type
    }
  }
}
