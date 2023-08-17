resource "border0_connector" "test_connector" {
  name                         = "test-connector"
  description                  = "test creating and managing connector from terraform"
  built_in_ssh_service_enabled = true
}