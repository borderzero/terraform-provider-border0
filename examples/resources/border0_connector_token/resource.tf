// first, create a connector
resource "border0_connector" "example" {
  name        = "example-connector"
  description = "My first connector created from terraform"
}

// next, create a token for the connector, and ensure that it never expires.
resource "border0_connector_token" "example_token_never_expires" {
  connector_id = border0_connector.example.id
  name         = "example-connector-token-never-expires"
}

// and create another connector token that expires at a specific date
resource "border0_connector_token" "example_token_expires" {
  connector_id = border0_connector.example.id
  name         = "example-connector-token-never-expires"
  expires_at   = "2023-12-31T23:59:59Z"
}
