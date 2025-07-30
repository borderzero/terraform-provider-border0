resource "border0_socket" "playground_http01" {
  # This socket is for HTTP traffic, using the first 1000 instances
  depends_on        = [border0_policy.my-tf-access-policy]
  recording_enabled = false
  for_each          = toset([for i in range(5000, 6000) : tostring(i)])
  name              = "${each.value}-playground-http-${each.value}"
  socket_type       = "http"
  connector_ids     = [border0_connector.first-connector.id, border0_connector.second-connector.id]
  description       = "My playground HTTP socket ${each.value}"

  http_configuration {
    upstream_url = "http://http.playground.border0.io"

    header {
      key    = "X-Custom-Header"
      values = ["custom-${each.value}", "another-${each.value}"]
    }
    header {
      key    = "X-Another-Header"
      values = ["yet-another-${each.value}"]
    }
  }
  tags = {
    "border0_client_subcategory" = "The Cloud"
    "border0_client_category"    = "Playground"
    "instance"                   = each.value
    "socket_type"                = "http"
  }
}

resource "border0_policy_attachment" "playground_http01-policy" {
  depends_on = [border0_policy.my-tf-access-policy, border0_policy.tf-dev-ops-policy]
  for_each   = border0_socket.playground_http01
  policy_id  = border0_policy.tf-dev-ops-policy.id
  socket_id  = each.value.id
}

# resource "border0_socket" "playground_https" {
#   recording_enabled = false
#   name              = "standalone-https-site"
#   socket_type       = "http"
#   connector_ids     = [border0_connector.first-connector.id, border0_connector.second-connector.id]
#   description       = "My playground HTTP socket default"

#   http_configuration {
#     upstream_url = "http://138.68.91.88"

#     header {
#       key    = "X-Custom-Header"
#       values = ["custom-default", "another-default"]
#     }
#     header {
#       key    = "X-Another-Header"
#       values = ["yet-another-default"]
#     }
#   }
#   tags = {
#     "border0_client_subcategory" = "The Cloud"
#     "border0_client_category"    = "Playground"
#   }
# }


# resource "border0_socket" "playground_http02" {
#   recording_enabled = false
#   for_each = toset([for i in range(6000, 6500) : tostring(i)])
#   name              = "playground-http-${each.value}"
#   socket_type       = "http"
#   connector_ids     = [border0_connector.first-connector.id, border0_connector.second-connector.id]
#   description       = "My playground HTTP socket ${each.value}"

#   http_configuration {
#     upstream_url = "http://http.playground.border0.io"

#     header {
#       key    = "X-Custom-Header"
#       values = ["custom-${each.value}", "another-${each.value}"]
#     }
#     header {
#       key    = "X-Another-Header"
#       values = ["yet-another-${each.value}"]
#     }
#   }
#   tags = {
#     "border0_client_subcategory" = "The Cloud"
#     "border0_client_category"    = "Playground"
#     "instance"                   = each.value
#   }
# }
