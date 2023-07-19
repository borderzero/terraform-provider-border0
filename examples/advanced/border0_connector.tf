// ===========================================================================
// Create a connector with plugins
// ===========================================================================

resource "border0_connector" "my_connector" {
  name = "my-connector"
  description = "my connector"

  aws_ec2_plugin {
    enabled = true
    scan_interval = 10 // in minutes

    // when authentication_strategy aws_profile
    authentication_strategy = "aws_profile"
    profile = "default"
    // when authentication_strategy static_credentials
    authentication_strategy = "static_credentials"
    access_key_id = "my-access-key-id"
    secret_key = "my-secret-key"

    regions = ["eu-central-1", "eu-west-1"]
    include_instance_states = ["running", "stopped"]
    include_tags = { // check how api accepts and stores include tags
      Name = "rollie"
      Environment = "production"
    }
    exclude_tags = { // check how api accepts and stores exclude tags
      Name = "rollie"
      Environment = "production"
      ToExclude = ""
    }
    enable_ssm_status_check = true
  }

  aws_ecs_plugin {
    enabled = true
    scan_interval = 10 // in minutes

    // when authentication_strategy aws_profile
    authentication_strategy = "aws_profile"
    profile = "default"
    // when authentication_strategy static_credentials
    authentication_strategy = "static_credentials"
    access_key_id = "my-access-key-id"
    secret_key = "my-secret-key"

    regions = ["eu-central-1", "eu-west-1"]
    include_tags = { // check how api accepts and stores include tags
      Name = "rollie"
      Environment = "production"
    }
    exclude_tags = { // check how api accepts and stores exclude tags
      Name = "rollie"
      Environment = "production"
      ToExclude = ""
    }
  }

  aws_rds_plugin {
    enabled = true
    scan_interval = 10 // in minutes

    // when authentication_strategy aws_profile
    authentication_strategy = "aws_profile"
    profile = "default"
    // when authentication_strategy static_credentials
    authentication_strategy = "static_credentials"
    access_key_id = "my-access-key-id"
    secret_key = "my-secret-key"

    regions = ["eu-central-1", "eu-west-1"]
    include_instance_states = ["running", "stopped"]
    include_tags = { // check how api accepts and stores include tags
      Name = "rollie"
      Environment = "production"
    }
    exclude_tags = { // check how api accepts and stores exclude tags
      Name = "rollie"
      Environment = "production"
      ToExclude = ""
    }
  }

  docker_plugin {
    enabled = true
    scan_interval = 10 // in minutes

    include_labels = { // check how api accepts and stores include labels
      Name = "rollie"
      Environment = "production"
    }
    exclude_labels = { // check how api accepts and stores exclude labels
      Name = "rollie"
      Environment = "production"
      ToExclude = ""
    }
  }

  kubernetes_plugin {
    enabled = true
    scan_interval = 10 // in minutes

    // when authentication_strategy is in_cluster_role
    authentication_strategy = "in_cluster_role"
    // when authentication_strategy kubeconfig_path
    authentication_strategy = "kubeconfig_path"
    kubeconfig_path = "/path/to/kubeconfig"

    namespaces = ["default", "kube-system"]
    include_labels = { // check how api accepts and stores include labels
      Name = "rollie"
      Environment = "production"
    }
    exclude_labels = { // check how api accepts and stores exclude labels
      Name = "rollie"
      Environment = "production"
      ToExclude = ""
    }
  }

  network_plugin {
    enabled = true
    scan_interval = 10 // in minutes

    targets = [ // see if it's possible to have a list of objects
      {
        target = "10.10.10.10" // ip, hostname or cidr
        ports = "80,443,8080" // comma separated list of ports
      }
    ]
  }
}

// ===========================================================================
// Create a connector token and associate it with a connector
// ===========================================================================

resource "border0_connector_token" "my_connector_token" {
  name = "my-connector-token"
  connector_id = border0_connector.my_connector.id
  expires_at = "2021-10-13T00:00:00Z" // optional, when not set token will never expire
}
