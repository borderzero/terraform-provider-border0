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
      "ssh" : {},
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
