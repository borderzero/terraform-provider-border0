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


// Using the data source to create/manage a tag based policy
// More info: https://docs.border0.com/docs/product-updates#march-2025
resource "border0_policy" "example" {
  name        = "example-policy"
  description = "My tag based policy"
  version     = "v2"
  policy_data = data.border0_policy_v2_document.example.json
  // This results in (example-tag1:example-value1 AND example_tag2:more-example-value OR example-tag:example-value)
  tag_rules = [
    { example_tag1 = "example-value1", example_tag2 = "more-example-value" },
    { example_tag = "example-value"},
  ]
}
