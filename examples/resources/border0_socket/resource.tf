resource "border0_socket" "test_http" {
  name = "test-http"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
  upstream_type = "https"
}

// create an SSH socket and link it to a connector that was created outside of Terraform
resource "border0_socket" "test_ssh" {
  name = "test-ssh"
  recording_enabled = true
  socket_type = "ssh"
  connector_id = "a7de4cc3-d977-4c4b-82e7-dedb6e7b74a1"
  upstream_hostname = "127.0.0.1"
  upstream_port = 22
  upstream_username = "test_user"
  upstream_connection_type = "ssh"
  upstream_authentication_type = "border0_cert"
}

// create another SSH socket and link it to a connector that was created in Terraform
resource "border0_socket" "test_another_ssh" {
  name = "test-another_ssh"
  recording_enabled = true
  socket_type = "ssh"
  connector_id = border0_connector.test_connector.id
  upstream_hostname = "127.0.0.1"
  upstream_port = 22
  upstream_username = "test_user"
  upstream_connection_type = "ssh"
  upstream_authentication_type = "username_password"
}
