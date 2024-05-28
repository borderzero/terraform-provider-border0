data "border0_policy_document" "example" {
  version = "v1"
  action  = ["database", "ssh", "http", "tls"]
  condition {
    who {
      email  = [] # your email goes here
      group  = []
      domain = ["example.com"]
    }
    where {
      allowed_ip  = ["0.0.0.0/0", "::/0"]
      country     = ["NL", "CA", "US", "BR", "FR"]
      country_not = ["BE"]
    }
    when {
      after              = "2022-10-13T05:12:27Z"
      time_of_day_after  = "00:00 UTC"
      time_of_day_before = "23:59 UTC"
    }
  }
}
