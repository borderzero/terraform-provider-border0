# variable for the user account email

data "border0_group_names_to_ids" "b0-dev-groups" {
  names = ["dev-ops", "sys-ops", "db-ops"]
}

data "border0_group_names_to_ids" "dev-ops-group" {
  names = ["dev-ops"]
}

data "border0_group_names_to_ids" "sys-ops-group" {
  names = ["sys-ops"]
}

data "border0_group_names_to_ids" "db-ops-group" {
  names = ["db-ops"]
}

# create policy
resource "border0_policy" "my-tf-access-policy" {
  name        = "my-tf-access-policy"
  description = "My terraform managed access policy"
  version     = "v2"
  policy_data = jsonencode({
    "permissions" : {
      "network" : {},
      "ssh" : {
      "shell" : {},
      "exec" : {},
      "sftp" : {},
      "tcp_forwarding" : {},
      "kubectl_exec" : {},
      "docker_exec" : {}
      },
      "database" : {
        "allowed_databases" : [
          {
            "database" : "books",
            "allowed_query_types" : [
              "ReadOnly"
            ]
          }
        ]
      },
      "http" : {},
      "tls" : {},
      "vnc" : {},
      "rdp" : {},
      "vpn" : {},
      "kubernetes" : {},
    },
    "condition" : {
      "who" : {
        "email" : [],
        "group" : data.border0_group_names_to_ids.b0-dev-groups.ids,
        "service_account" : []
      },
      "where" : {
        "allowed_ip" : [
          "0.0.0.0/0",
          "::/0"
        ],
        "country" : null,
        "country_not" : null
      },
      "when" : {
        "after" : "2022-02-02T22:22:22Z",
        "before" : null,
        "time_of_day_after" : "",
        "time_of_day_before" : ""
      }
    }
  })
}


resource "border0_policy" "tf-dev-ops-policy" {
  name        = "tf-dev-ops-policy"
  description = "My terraform managed dev-ops access policy"
  version     = "v2"
  policy_data = jsonencode({
    "permissions" : {
      "network" : {},
      "ssh" : {
        "shell" : {},
        "exec" : {},
        "sftp" : {},
        "tcp_forwarding" : {},
        "kubectl_exec" : {},
        "docker_exec" : {}
      },
      "database" : {
        "allowed_databases" : [
          {
            "database" : "books",
            "allowed_query_types" : [
              "ReadOnly"
            ]
          }
        ]
      },
      "http" : {},
      "tls" : {},
      "vnc" : {},
      "rdp" : {},
      "vpn" : {},
      "kubernetes" : {},
    },
    "condition" : {
      "who" : {
        "email" : [],
        "group" : data.border0_group_names_to_ids.dev-ops-group.ids,
        "service_account" : []
      },
      "where" : {
        "allowed_ip" : [
          "0.0.0.0/0",
          "::/0"
        ],
        "country" : null,
        "country_not" : null
      },
      "when" : {
        "after" : "2022-02-02T22:22:22Z",
        "before" : null,
        "time_of_day_after" : "",
        "time_of_day_before" : ""
      }
    }
  })
}

resource "border0_policy" "tf-sys-ops-policy" {
  name        = "tf-sys-ops-policy"
  description = "My terraform managed sys-ops access policy"
  version     = "v2"
  policy_data = jsonencode({
    "permissions" : {
      "network" : {},
      "ssh" : {
        "shell" : {},
        "exec" : {},
        "sftp" : {},
        "tcp_forwarding" : {},
        "kubectl_exec" : {},
        "docker_exec" : {}
      },
      "database" : {
        "allowed_databases" : [
          {
            "database" : "books",
            "allowed_query_types" : [
              "ReadOnly"
            ]
          }
        ]
      },
      "http" : {},
      "tls" : {},
      "vnc" : {},
      "rdp" : {},
      "vpn" : {},
      "kubernetes" : {},
    },
    "condition" : {
      "who" : {
        "email" : [],
        "group" : data.border0_group_names_to_ids.sys-ops-group.ids,
        "service_account" : []
      },
      "where" : {
        "allowed_ip" : [
          "0.0.0.0/0",
          "::/0"
        ],
        "country" : null,
        "country_not" : null
      },
      "when" : {
        "after" : "2022-02-02T22:22:22Z",
        "before" : null,
        "time_of_day_after" : "",
        "time_of_day_before" : ""
      }
    }
  })
}


resource "border0_policy" "tf-db-ops-policy" {
  name        = "tf-db-ops-policy"
  description = "My terraform managed db-ops access policy"
  version     = "v2"
  policy_data = jsonencode({
    "permissions" : {
      "network" : {},
      "ssh" : {
        "shell" : {},
        "exec" : {},
        "sftp" : {},
        "tcp_forwarding" : {},
        "kubectl_exec" : {},
        "docker_exec" : {}
      },
      "database" : {
        "allowed_databases" : [
          {
            "database" : "books",
            "allowed_query_types" : [
              "ReadOnly"
            ]
          }
        ]
      },
      "http" : {},
      "tls" : {},
      "vnc" : {},
      "rdp" : {},
      "vpn" : {},
      "kubernetes" : {},
    },
    "condition" : {
      "who" : {
        "email" : [],
        "group" : data.border0_group_names_to_ids.db-ops-group.ids,
        "service_account" : []
      },
      "where" : {
        "allowed_ip" : [
          "0.0.0.0/0",
          "::/0"
        ],
        "country" : null,
        "country_not" : null
      },
      "when" : {
        "after" : "2022-02-02T22:22:22Z",
        "before" : null,
        "time_of_day_after" : "",
        "time_of_day_before" : ""
      }
    }
  })
}

# those are empty policies for use with socket tags
resource "border0_policy" "tf-ops-tag-policy" {
  name        = "tf-ops-tag-policy"
  description = "My terraform managed ops-tag access policy"
  version     = "v2"
  policy_data = jsonencode({
    "permissions" : {
      "database" : {},
      "http" : {},
      "network" : {},
      "ssh" : {},
    },
    "condition" : {
      "who" : {
      },
      "where" : {
        "allowed_ip" : ["0.0.0.0/0", "::/0"],
      },
      "when" : {
        "after" : "2022-10-13T05:12:26Z",
        "before" : null,
      }
    }
  })
}