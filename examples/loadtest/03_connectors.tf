# create connector 
resource "border0_connector" "first-connector" {
  name                         = "first-connector"
  description                  = "My first connector created from terraform"
}

# create connector token
resource "border0_connector_token" "first-connector_token" {
  expires_at   = "2036-12-31T23:59:59Z"
  connector_id = border0_connector.first-connector.id
  name         = "first-connector-token"
}

# generate minimal connector configuration file with the token
resource "local_file" "write_token_first-connector-token" {
  content  = "token: ${border0_connector_token.first-connector_token.token}"
  filename = "${path.module}/first-connector.yaml"
}

# second connector
resource "border0_connector" "second-connector" {
  name        = "second-connector"
  description = "My second connector created from terraform"
}

# create connector token
resource "border0_connector_token" "second-connector_token" {
  expires_at   = "2036-12-31T23:59:59Z"
  connector_id = border0_connector.second-connector.id
  name         = "second-connector-token"
}

# generate minimal connector configuration file with the token
resource "local_file" "write_token_second-connector-token" {
  content  = "token: ${border0_connector_token.second-connector_token.token} \n"
  filename = "${path.module}/second-connector.yaml"
}