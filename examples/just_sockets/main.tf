// use the provider from terraform registry
terraform {
  required_providers {
    border0 = {
      source = "borderzero/border0"
    }
  }
}

variable "token" {
  type = string
}

provider "border0" {
  token = var.token
}

resource "border0_socket" "test_tf_http" {
  name          = "test-tf-http"
  socket_type   = "http"
  connector_ids = ["a7de4cc3-d977-4c4b-82e7-dedb6e7b74a1"] // replace with your connector ID

  http_configuration {
    upstream_url = "https://www.bbc.com"
    header {
      key    = "X-Custom-Header"
      values = ["custom-value", "another-value"]
    }
  }

  tags = {
    "test_key_1" = "test_value_1"
  }
}

resource "border0_socket" "test_tf_ssh" {
  name              = "test-tf-ssh"
  recording_enabled = true
  socket_type       = "ssh"
  connector_ids     = ["a7de4cc3-d977-4c4b-82e7-dedb6e7b74a1"] // replace with your connector ID

  ssh_configuration {
    hostname            = "127.0.0.1"
    port                = 22
    username            = "test_user"
    authentication_type = "border0_certificate"
  }
}

output "managed_socket" {
  value = {
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
