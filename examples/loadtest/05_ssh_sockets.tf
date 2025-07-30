resource "border0_socket" "playground_ssh" {
  depends_on = [border0_policy.my-tf-access-policy]
  for_each   = toset([for i in range(3000, 4000) : tostring(i)])

  recording_enabled = false
  name              = "${each.value}-playground-ssh-${each.value}"
  socket_type       = "ssh"
  connector_ids     = [border0_connector.first-connector.id, border0_connector.second-connector.id]

  description = "My playground SSH socket ${each.value}"

  ssh_configuration {
    hostname            = "ssh.playground.border0.io"
    port                = 22
    authentication_type = "username_and_password"
    username            = "border0"
    password            = "Border0<3Ssh"
  }
  tags = {
    "border0_client_subcategory" = "The Cloud"
    "border0_client_category"    = "Playground"
    "instance"                   = each.value
    "socket_type"                = "ssh"
  }
}

resource "border0_policy_attachment" "playground_ssh-policy" {
  depends_on = [border0_policy.my-tf-access-policy, border0_policy.tf-sys-ops-policy]
  for_each   = border0_socket.playground_ssh
  policy_id  = border0_policy.tf-sys-ops-policy.id
  socket_id  = each.value.id
}


# resource "border0_socket" "playground_ssh_second" {
#   for_each = toset([for i in range(4000, 4100) : tostring(i)])

#   recording_enabled = false
#   name              = "playground-ssh-second-${each.value}"
#   socket_type       = "ssh"
#   connector_ids     = [border0_connector.first-connector.id, border0_connector.second-connector.id]

#   description = "My playground SSH socket second ${each.value}"

#   ssh_configuration {
#     hostname            = "ssh.playground.border0.io"
#     port                = 22
#     authentication_type = "username_and_password"
#     username            = "border0"
#     password            = "Border0<3Ssh"
#   }
#   tags = {
#     "border0_client_subcategory" = "The Cloud"
#     "border0_client_category"    = "Playground"
#     "instance"                   = each.value
#   }
# }
