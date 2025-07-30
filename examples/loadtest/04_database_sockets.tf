resource "border0_socket" "playground_mysql-first-1000" {
  depends_on = [border0_policy.my-tf-access-policy]
  for_each   = toset([for i in range(1000, 2000) : tostring(i)])

  recording_enabled = false
  name              = "${each.value}-playground-mysql-${each.value}"
  socket_type       = "database"
  #   connector_id      = border0_connector.first-connector.id
  connector_ids = [border0_connector.first-connector.id, border0_connector.second-connector.id]
  description   = "My playground MySQL socket ${each.value}"

  database_configuration {
    protocol            = "mysql"
    hostname            = "mysql.playground.border0.io"
    port                = 3306
    authentication_type = "username_and_password"
    username            = "border0"
    password            = "Border0<3MySql"
  }
  tags = {
    "border0_client_subcategory" = "The Cloud"
    "border0_client_category"    = "Playground"
    "instance"                   = each.value
    "socket_type"                = "database"
  }
}

resource "border0_policy_attachment" "playground_mysql-first-1000-policy" {
  depends_on = [border0_policy.my-tf-access-policy, border0_policy.tf-db-ops-policy]
  for_each   = border0_socket.playground_mysql-first-1000
  policy_id  = border0_policy.tf-db-ops-policy.id
  socket_id  = each.value.id
}


# resource "border0_socket" "playground_mysql-second-1000" {
#   for_each = toset([for i in range(2000, 2100) : tostring(i)])

#   recording_enabled = false
#   name              = "playground-mysql-${each.value}"
#   socket_type       = "database"
#   connector_ids     = [border0_connector.first-connector.id, border0_connector.second-connector.id]
#   description       = "My playground MySQL socket ${each.value}"

#   database_configuration {
#     protocol            = "mysql"
#     hostname            = "mysql.playground.border0.io"
#     port                = 3306
#     authentication_type = "username_and_password"
#     username            = "border0"
#     password            = "Border0<3MySql"
#   }
#   tags = {
#     "border0_client_subcategory" = "The Cloud"
#     "border0_client_category"    = "Playground"
#     "instance"                   = each.value
#   }
# }
