resource "border0_connector" "example" {
  name                         = "example-connector"
  description                  = "My first connector created from terraform"
  built_in_ssh_service_enabled = true
}
