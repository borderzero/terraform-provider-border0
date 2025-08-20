
data "border0_policy_v2_document" "example" {
  permissions {
    ssh {
      allowed = true
      shell { allowed = true }
      exec { allowed = true }
      sftp { allowed = true }
      tcp_forwarding { allowed = true }
      kubectl_exec { allowed = true }
      docker_exec { allowed = true }
    }

    database {
      allowed = true
      allowed_databases {
        database            = "books"
        allowed_query_types = ["ReadOnly"]
      }
    }

    http { allowed = true }
    tls { allowed = true }
    vnc { allowed = true }
    rdp { allowed = true }
    network { allowed = true }
    kubernetes { allowed = true }
  }

  condition {
    who {
      email           = ["example@example.com"]
      group           = ["38bf96d7-1fa2-4a71-b010-361a64355ba7"]
      service_account = ["example-service-account"]
    }
    where {
      allowed_ip = ["0.0.0.0/0", "::/0"]
      country     = ["NL", "CA", "US", "BR", "FR"]
      country_not = ["RU"]
    }
    when {
      after              = "2022-02-02T22:22:22Z"
      before             = null
      time_of_day_after  = null
      time_of_day_before = null
    }
  }
}

# Using the data source to create/manage a new policy
resource "border0_policy" "example" {
  name        = "example-policy"
  policy_data = data.border0_policy_v2_document.example.json
}
