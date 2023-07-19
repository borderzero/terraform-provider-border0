// ===========================================================================
// Create a policy using a json string
// ===========================================================================

resource "border0_policy" "my_policy" {
  name = "my-policy"
  description = "my policy"
  org_wide = false
  policy = <<EOF
{
  "version": "v1",
  "action": [ "database", "ssh", "http", "tls" ],
  "condition": {
    "who": {
      "email": [ "johndoe@example.com" ],
      "domain": [ "example.com" ]
    },
    "where": {
      "allowed_ip": [ "0.0.0.0/0", "::/0" ],
      "country": [ "NL", "CA", "US", "BR", "FR" ],
      "country_not": [ "BE" ]
    },
    "when": {
      "after": "2022-10-13",
      "before": null,
      "time_of_day_after": "00:00:00 UTC",
      "time_of_day_before": "23:59:59 UTC"
    }
  }
}
EOF
}

// ===========================================================================
// Create another policy using a data source with easier to read syntax
// ===========================================================================

data "border0_policy_document" "another_policy" {
  version = "v1"
  action = [ "database", "ssh", "http", "tls" ]
  condition = {
    who = {
      email = [ "johndoe@example.com" ]
      domain = [ "example.com" ]
    }
    where = {
      allowed_ip = [ "0.0.0.0/0", "::/0" ]
      country = [ "NL", "CA", "US", "BR", "FR" ]
      country_not = [ "BE" ]
    }
    when = {
      after = "2022-10-13"
      before = null // should this be omitted?
      time_of_day_after = "00:00:00 UTC"
      time_of_day_before = "23:59:59 UTC"
    }
  }
}

resource "border0_policy" "another_policy" {
  name = "my-policy"
  description = "my policy"
  org_wide = false
  policy = data.border0_policy_document.another_policy.json
}
