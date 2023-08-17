// create an HTTP socket with an HTTPS upstream and add few tags to the socket
// this socket will not be linked to any connector
resource "border0_socket" "test_http" {
  name        = "test-http"
  socket_type = "http"
  tags = {
    "user"        = "Bilbo Baggins"
    "project"     = "The Hobbit"
    "region"      = "The Shire"
    "environment" = "dev"
  }
  upstream_type          = "https"
  upstream_http_hostname = "www.bbc.com"
}

// create an SSH socket and link it to a connector that's not managed by Terraform
resource "border0_socket" "test_ssh" {
  name                         = "test-ssh"
  recording_enabled            = true
  socket_type                  = "ssh"
  connector_id                 = "a7de4cc3-d977-4c4b-82e7-dedb6e7b74a1"
  upstream_hostname            = "127.0.0.1"
  upstream_port                = 22
  upstream_username            = "test_user"
  upstream_connection_type     = "ssh"
  upstream_authentication_type = "border0_cert"
}

// create another SSH socket and link it to a connector that was created with Terraform
resource "border0_socket" "test_another_ssh" {
  name                         = "test-another-ssh"
  recording_enabled            = true
  socket_type                  = "ssh"
  connector_id                 = border0_connector.test_connector.id
  upstream_hostname            = "127.0.0.1"
  upstream_port                = 22
  upstream_username            = "test_user"
  upstream_password            = "test_password"
  upstream_connection_type     = "ssh"
  upstream_authentication_type = "username_password"
}
