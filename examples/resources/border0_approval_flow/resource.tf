resource "border0_approval_flow" "prod_db_access" {
  name        = "Production database access"
  description = "Requires approval from the SRE on-call before accessing production databases"

  socket_tags = {
    env  = "prod"
    type = "database"
  }

  requester_group_ids = [border0_group.engineers.id]
  approver_group_ids  = [border0_group.sre.id]

  allow_self_approval = false
}
