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
  policy_data = jsonencode({
    "version" : "v1",
    "action" : ["database", "ssh", "http", "tls"],
    "condition" : {
      "who" : {
        "email" : ["johndoe@example.com"],
        "domain" : ["example.com"]
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
