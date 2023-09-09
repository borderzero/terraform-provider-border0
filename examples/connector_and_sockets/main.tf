// use the provider from terraform registry
terraform {
  required_providers {
    border0 = {
      source  = "borderzero/border0"
      version = "0.1.10"
    }
  }
}

variable "token" {
  type = string
}

provider "border0" {
  token = var.token
}

resource "border0_connector" "test_tf_connector" {
  name                         = "test-tf-connector"
  description                  = "test connector from terraform"
  built_in_ssh_service_enabled = true
}

resource "border0_connector_token" "test_tf_connector_token_never_expires" {
  connector_id = border0_connector.test_tf_connector.id
  name         = "test-tf-connector-token-never-expires"
}

resource "border0_connector_token" "test_tf_connector_token_expires" {
  connector_id = border0_connector.test_tf_connector.id
  name         = "test-tf-connector-token-never-expires"
  expires_at   = "2023-12-31T23:59:59Z"
}

resource "border0_socket" "test_tf_http" {
  name        = "test-tf-http"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
  upstream_type = "https"
}

resource "border0_socket" "test_tf_http" {
  name         = "test-tf-http"
  socket_type  = "http"
  connector_id = border0_connector.test_tf_connector.id

  http_configuration {
    hostname    = "www.bbc.com"
    port        = 443
    host_header = "www.bbc.com"
  }
  upstream_type = "https"

  tags = {
    "test_key_1" = "test_value_1"
  }
}

resource "border0_socket" "test_tf_ssh" {
  name              = "test-tf-ssh"
  recording_enabled = true
  socket_type       = "ssh"
  connector_id      = border0_connector.test_tf_connector.id

  ssh_configuration {
    hostname            = "127.0.0.1"
    port                = 22
    username            = "test_user"
    authentication_type = "border0_certificate"
  }
}

// create a database socket and link it to a connector that was created with Terraform
// this socket will be used to connect to an AWS RDS instance with IAM authentication
resource "border0_socket" "test_tf_aws_rds_with_iam_auth" {
  name              = "test-tf-aws-rds-with-iam-auth"
  recording_enabled = true
  socket_type       = "database"
  connector_id      = border0_connector.test_tf_connector.id

  database_configuration {
    protocol            = "mysql"
    hostname            = "some-aws-rds-cluster.us-west-2.rds.amazonaws.com"
    port                = 3306
    service_type        = "aws_rds"
    authentication_type = "iam"
    rds_instance_region = "us-east-2"
    username            = "some_db_iam_user_name"
  }
}

// create an SSH socket and link it to a connector that was created with Terraform
// this socket will be used to connect to an AWS EC2 instance with EC2 Instance Connect
resource "border0_socket" "test_tf_aws_ec2_instance_connect" {
  name              = "test-tf-ec2-instance-connect"
  recording_enabled = true
  socket_type       = "ssh"
  connector_id      = border0_connector.test_tf_connector.id

  ssh_configuration {
    service_type        = "aws_ec2_instance_connect"
    hostname            = "10.0.0.101"
    port                = 22
    username_provider   = "defined"
    username            = "ubuntu"
    ec2_instance_id     = "i-00000000000000001"
    ec2_instance_region = "ap-southeast-2"
  }
}

// create an SSH socket and link it to a connector that was created with Terraform
// this socket will be used to connect to an AWS ECS service with SSM Session Manager
resource "border0_socket" "test_tf_connect_to_ecs_with_ssm" {
  name              = "test-tf-connect-to-ecs-with-ssm"
  recording_enabled = true
  socket_type       = "ssh"
  connector_id      = border0_connector.test_tf_connector.id

  ssh_configuration {
    service_type       = "aws_ssm"
    ssm_target_type    = "ecs"
    ecs_cluster_region = "eu-west-1"
    ecs_cluster_name   = "some-ecs-cluster-name"
    ecs_service_name   = "some-ecs-service-name"
  }
}

output "managed_resources" {
  value = {
    connector = {
      id   = border0_connector.test_tf_connector.id
      name = border0_connector.test_tf_connector.name
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
    aws_rds_socket = {
      id   = border0_socket.test_tf_aws_rds_with_iam_auth.id
      name = border0_socket.test_tf_aws_rds_with_iam_auth.name
      type = border0_socket.test_tf_aws_rds_with_iam_auth.socket_type
    }
    aws_ec2_instance_connect_socket = {
      id   = border0_socket.test_tf_aws_ec2_instance_connect.id
      name = border0_socket.test_tf_aws_ec2_instance_connect.name
      type = border0_socket.test_tf_aws_ec2_instance_connect.socket_type
    }
    aws_ecs_ssm_socket = {
      id   = border0_socket.test_tf_connect_to_ecs_with_ssm.id
      name = border0_socket.test_tf_connect_to_ecs_with_ssm.name
      type = border0_socket.test_tf_connect_to_ecs_with_ssm.socket_type
    }
  }
}
