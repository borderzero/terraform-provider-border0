// use the provider from local path (for development)
terraform {
  required_providers {
    border0 = {
      source  = "border0.com/border0/border0"
      version = "0.1.0" # keep this as 0.1.0, it's a dummy version for local build
    }
  }
}

variable "token" { type = string }
variable "api_url" { type = string }

provider "border0" {
  token   = var.token
  api_url = var.api_url
}

resource "border0_connector" "test_tf_connector" {
  name                         = "test-tf-connector-1"
  description                  = "test connector from terraform"
  built_in_ssh_service_enabled = true
}

resource "border0_connector" "test_tf_connector_2" {
  name                         = "test-tf-connector-2"
  description                  = "test connector from terraform"
  built_in_ssh_service_enabled = true
}

resource "border0_connector_token" "test_tf_connector_token_never_expires" {
  connector_id = border0_connector.test_tf_connector.id
  name         = "test-tf-connector-token-never-expires"

  provisioner "local-exec" {
    command = "echo 'token: ${self.token}' > ./border0.yaml"
  }
}

resource "border0_connector_token" "test_tf_connector_token_expires" {
  connector_id = border0_connector.test_tf_connector.id
  name         = "test-tf-connector-token-expires"
  expires_at   = "2035-12-31T23:59:59Z"

  provisioner "local-exec" {
    command = "echo 'token: ${self.token}' > ./border0-connector-token-expires.yaml"
  }
}

resource "border0_policy" "test_tf_policy" {
  name        = "test-tf-policy"
  description = "test policy from terraform"
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

data "border0_policy_v2_document" "test_tf_policy_document" {
  permissions {
    network { allowed = true }
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
    kubernetes { allowed = true }
  }
  condition {
    who {
      email = [] # your email goes here
      group = []
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

resource "border0_policy" "another_test_tf_policy" {
  name        = "another-test-tf-policy"
  description = "another test policy from terraform"
  policy_data = data.border0_policy_v2_document.test_tf_policy_document.json
}

resource "border0_socket" "test_tf_http" {
  name        = "test-tf-http"
  socket_type = "http"
  connector_ids = [
    border0_connector.test_tf_connector.id,
    border0_connector.test_tf_connector_2.id,
  ]

  http_configuration {
    upstream_url = "https://www.bbc.com"
    header {
      key    = "X-Custom-Header"
      values = ["custom-value"]
    }
    header {
      key    = "X-Another-Header"
      values = ["another-value", "second-value"]
    }
  }
  tags = {
    "test_key_1" = "test_value_1"
  }
}


resource "border0_socket" "test_tf_ssh" {
  name              = "test-tf-ssh"
  recording_enabled = true
  socket_type       = "ssh"
  connector_ids     = [border0_connector.test_tf_connector.id]

  ssh_configuration {
    hostname            = "127.0.0.1"
    port                = 22
    username            = "test_user"
    authentication_type = "border0_certificate"
  }
}

resource "border0_policy_attachment" "test_tf_policy_builtin_ssh" {
  policy_id = border0_policy.test_tf_policy.id
  socket_id = border0_connector.test_tf_connector.built_in_ssh_service_id
}

resource "border0_policy_attachment" "another_test_tf_policy_builtin_ssh" {
  policy_id = border0_policy.another_test_tf_policy.id
  socket_id = border0_connector.test_tf_connector.built_in_ssh_service_id
}

resource "border0_policy_attachment" "test_tf_policy_http_socket" {
  policy_id = border0_policy.test_tf_policy.id
  socket_id = border0_socket.test_tf_http.id
}

resource "border0_policy_attachment" "test_tf_policy_ssh_socket" {
  policy_id = border0_policy.test_tf_policy.id
  socket_id = border0_socket.test_tf_ssh.id
}

resource "border0_policy_attachment" "another_test_tf_policy_ssh_socket" {
  policy_id = border0_policy.another_test_tf_policy.id
  socket_id = border0_socket.test_tf_ssh.id
}

resource "border0_socket" "test_tf_mysql" {
  name              = "test-tf-MySQL"
  recording_enabled = true
  socket_type       = "database"
  connector_ids     = [border0_connector.test_tf_connector.id]

  database_configuration {
    protocol      = "mysql"
    hostname      = "127.0.0.1"
    port          = 3307
    database_name = "test_db"
    username      = "root"
    password      = "test"
  }
}

resource "border0_socket" "test_tf_tls" {
  name          = "test-tf-tls"
  socket_type   = "tls"
  connector_ids = [border0_connector.test_tf_connector.id]

  tls_configuration {
    hostname = "127.0.0.1"
    port     = 4242
  }
}

resource "border0_socket" "test_tf_vnc" {
  name          = "test-tf-vnc"
  socket_type   = "vnc"
  connector_ids = [border0_connector.test_tf_connector.id]

  vnc_configuration {
    hostname = "127.0.0.1"
    port     = 5900
  }
}

resource "border0_socket" "test_tf_rdp" {
  name          = "test-tf-rdp"
  socket_type   = "rdp"
  connector_ids = [border0_connector.test_tf_connector.id]

  rdp_configuration {
    hostname = "127.0.0.1"
    port     = 3389
  }
}

output "managed_resources" {
  value = {
    connector = {
      id                           = border0_connector.test_tf_connector.id
      name                         = border0_connector.test_tf_connector.name
      built_in_ssh_service_enabled = border0_connector.test_tf_connector.built_in_ssh_service_enabled
      built_in_ssh_service_id      = border0_connector.test_tf_connector.built_in_ssh_service_id
    }
    connector_token = {
      id           = border0_connector_token.test_tf_connector_token_never_expires.id
      connector_id = border0_connector_token.test_tf_connector_token_never_expires.connector_id
      name         = border0_connector_token.test_tf_connector_token_never_expires.name
      expires_at   = border0_connector_token.test_tf_connector_token_never_expires.expires_at
    }
    another_connector_token = {
      id           = border0_connector_token.test_tf_connector_token_expires.id
      connector_id = border0_connector_token.test_tf_connector_token_expires.connector_id
      name         = border0_connector_token.test_tf_connector_token_expires.name
      expires_at   = border0_connector_token.test_tf_connector_token_expires.expires_at
    }
    policy = {
      id   = border0_policy.test_tf_policy.id
      name = border0_policy.test_tf_policy.name
    }
    another_policy = {
      id   = border0_policy.another_test_tf_policy.id
      name = border0_policy.another_test_tf_policy.name
    }
    http_socket = {
      id   = border0_socket.test_tf_http.id
      name = border0_socket.test_tf_http.name
      type = border0_socket.test_tf_http.socket_type
    }
    ssh_socket = {
      id   = border0_socket.test_tf_ssh.id
      name = border0_socket.test_tf_ssh.name
      type = border0_socket.test_tf_ssh.socket_type
    }
  }
}
