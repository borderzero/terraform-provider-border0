terraform {
  required_providers {
    border0 = {
      source = "border0.com/border0/border0"
      version = "0.1.0"
    }
  }
}

resource "border0_socket" "test_tf_http" {
  name = "test-tf-http"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
}

resource "border0_socket" "test_tf_ssh" {
  name = "test-tf-ssh"
  recording_enabled = true
  socket_type = "ssh"
}

output "managed_socket" {
  value = {
    http = border0_socket.test_tf_http
    ssh = border0_socket.test_tf_ssh
  }
}
