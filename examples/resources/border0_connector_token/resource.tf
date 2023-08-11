resource "border0_connector" "test_connector" {
  name = "test-connector"
  description = "test creating and managing connector from terraform"
}

// create a connector token that never expires
resource "border0_connector_token" "test_connector_token_never_expires" {
  connector_id = border0_connector.test_connector.id
  name = "test-connector-token-never-expires"
}

// create another connector token that expires at a specific date
resource "border0_connector_token" "test_connector_token_expires" {
  connector_id = border0_connector.test_connector.id
  name = "test-connector-token-never-expires"
  expires_at = "2023-12-31T23:59:59Z"
}
