// create an HTTP socket with an HTTPS upstream and add few tags to the socket
// this socket will be linked to a connector that was created with terraform
resource "border0_socket" "example_http" {
  name         = "example-http"
  socket_type  = "http"
  connector_id = border0_connector.example.id // link to a connector that was created with terraform

  http_configuration {
    hostname    = "www.bbc.com"
    port        = 443
    host_header = "www.bbc.com"
  }
  upstream_type = "https"

  tags = {
    "user"        = "Bilbo Baggins"
    "project"     = "The Hobbit"
    "region"      = "The Shire"
    "environment" = "dev"
  }
}

// create an SSH socket and link it to a connector that's not managed by terraform
resource "border0_socket" "example_ssh_password_auth" {
  name              = "example-ssh-password-auth"
  recording_enabled = true
  socket_type       = "ssh"
  connector_id      = "a7de4cc3-d977-4c4b-82e7-dedb6e7b74a1" // replace with your connector ID

  ssh_configuration {
    hostname            = "127.0.0.1"
    port                = 22
    authentication_type = "username_and_password"
    username            = "some_user"
    password            = "from:file:/path/to/password/file"
  }
}

// create another SSH socket and link it to a connector that was created with terraform
resource "border0_socket" "example_ssh_border0_certificate_auth" {
  name              = "example-ssh-border0-certificate-auth"
  recording_enabled = true
  socket_type       = "ssh"
  connector_id      = border0_connector.example.id // link to a connector that was created with terraform

  ssh_configuration {
    hostname            = "127.0.0.1"
    port                = 22
    authentication_type = "border0_certificate"
    username            = "some_user"
  }
}

// create a database socket and link it to a connector that was created with terraform
// this socket will be used to connect to an AWS RDS instance with IAM authentication
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html
resource "border0_socket" "example_aws_rds_with_iam_auth" {
  name              = "example-aws-rds-with-iam-auth"
  recording_enabled = true
  socket_type       = "database"
  connector_id      = border0_connector.example.id // link to a connector that was created with terraform

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

// create an SSH socket and link it to a connector that was created with terraform
// this socket will be used to connect to an AWS EC2 instance with EC2 Instance Connect
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-connect-set-up.html
resource "border0_socket" "example_aws_ec2_instance_connect" {
  name              = "example-ec2-instance-connect"
  recording_enabled = true
  socket_type       = "ssh"
  connector_id      = border0_connector.example.id

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

// create an SSH socket and link it to a connector that was created with terraform
// this socket will be used to connect to an AWS ECS service with SSM Session Manager
// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-exec.html
// https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager.html
resource "border0_socket" "example_connect_to_ecs_with_ssm" {
  name              = "example-connect-to-ecs-with-ssm"
  recording_enabled = true
  socket_type       = "ssh"
  connector_id      = border0_connector.example.id // link to a connector that was created with terraform

  ssh_configuration {
    service_type       = "aws_ssm"
    ssm_target_type    = "ecs"
    ecs_cluster_region = "eu-west-1"
    ecs_cluster_name   = "some-ecs-cluster-name"
    ecs_service_name   = "some-ecs-service-name"
  }
}
