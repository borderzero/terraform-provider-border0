// use the provider from terraform registry
terraform {
  required_providers {
    border0 = {
      source = "borderzero/border0"
      version = "0.1.1"
    }
  }
}

// use the provider from local path (for development)
// terraform {
//   required_providers {
//     border0 = {
//       source = "border0.com/border0/border0"
//       version = "0.1.0"
//     }
//   }
// }

resource "border0_socket" "test_tf_http" {
  name = "test-tf-http"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
  upstream_type = "https"
}

resource "border0_socket" "test_tf_ssh" {
  name = "test-tf-ssh"
  recording_enabled = true
  socket_type = "ssh"
  connector_id = "a7de4cc3-d977-4c4b-82e7-dedb6e7b74a1"
  upstream_hostname = "127.0.0.1"
  upstream_port = 22
  upstream_authentication_type = "border0_cert"
}

output "managed_socket" {
  value = {
    http = {
      id = border0_socket.test_tf_http.id
      name = border0_socket.test_tf_http.name
      type = border0_socket.test_tf_http.socket_type
      upstream_type = border0_socket.test_tf_http.upstream_type
      tags = border0_socket.test_tf_http.tags
    }
    ssh = {
      id = border0_socket.test_tf_ssh.id
      name = border0_socket.test_tf_ssh.name
      type = border0_socket.test_tf_ssh.socket_type
      recording_enabled = border0_socket.test_tf_ssh.recording_enabled
      connector_id = border0_socket.test_tf_ssh.connector_id
      upstream_hostname = border0_socket.test_tf_ssh.upstream_hostname
      upstream_port = border0_socket.test_tf_ssh.upstream_port
      upstream_authentication_type = border0_socket.test_tf_ssh.upstream_authentication_type
    }
  }
}
