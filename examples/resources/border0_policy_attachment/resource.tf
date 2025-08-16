// first, create an http socket
resource "border0_socket" "example" {
  name                   = "my-example-http"
  socket_type            = "http"
  upstream_type          = "https"
  upstream_http_hostname = "www.bbc.com"
}

// and then create a policy
resource "border0_policy" "example" {
  name        = "my-example-policy"
  description = "My first policy"
  version     = "v2"
  policy_data = jsonencode({
    "permissions" : {
      "database" : {},
      "http" : {},
      "kubernetes" : {},
      "network" : {},
      "rdp" : {},
      "ssh" : {
        "shell" : {},
        "exec" : {},
        "sftp" : {},
        "tcp_forwarding" : {},
        "kubectl_exec" : {},
        "docker_exec" : {}
      },
      "tls" : {},
      "vnc" : {}
    },
    "condition" : {
      "who" : {
        "email" : [], # your email goes here
        "group" : [],
        "service_account" : []
      },
      "where" : {
        "allowed_ip" : ["0.0.0.0/0", "::/0"],
        "country" : ["NL", "CA", "US", "BR", "FR"],
        "country_not" : ["BE"]
      },
      "when" : {
        "after" : "2022-10-13T05:12:26Z",
        "before" : null,
        "time_of_day_after" : "00:00 UTC",
        "time_of_day_before" : "23:59 UTC"
      }
    }
  })
}

// finally, attach the policy to the socket
resource "border0_policy_attachment" "example" {
  policy_id = border0_policy.example.id
  socket_id = border0_socket.example.id
}
