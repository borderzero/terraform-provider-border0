// ===========================================================================
// HTTP Sockets
// ===========================================================================

resource "border0_socket" "https_with_remote_host" {
  name = "https-with-remote-host"
  socket_type = "http"
  connector_id = border0_connector.my_connector.id

  // optional fileds
  description = "https socket with remote host"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  upstream_type = "https"
  upstream_http_hostname = "www.bbc.com"
}

resource "border0_socket" "https_with_builtin_web_server" {
  name = "https-with-builtin-web-server"
  socket_type = "http"
  connector_id = border0_connector.my_connector.id

  // optional fileds
  description = "https socket with builtin web server"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  // need to discuss how to handle this
  builtin_web_server = true
}

// ===========================================================================
// SSH Sockets
// ===========================================================================

resource "border0_socket" "ssh_with_builtin_ssh_server" {
  name = "ssh-with-builtin-ssh-server"
  socket_type = "ssh"
  connector_id = border0_connector.my_connector.id

  // optional fileds
  description = "ssh socket with builtin ssh server"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  // need to discuss how to handle this
  builtin_ssh_server = true
}

resource "border0_socket" "ssh_with_username_and_password" {
  name = "ssh-with-username-and-password"
  socket_type = "ssh"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "ssh socket with username and password"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  upstream_username = "root"
  upstream_password = "password"
  upstream_hostname = "127.0.0.1" // change from host to upstream_hostname
  upstream_port = 22 // chanage from port to upstream_port
  // should we add another field for upstream auth type?
  upstream_authentication_type = "username_password"
}

resource "border0_socket" "ssh_with_username_and_private_key" {
  name = "ssh-with-username-and-private-key"
  socket_type = "ssh"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "ssh socket with username and private key"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  upstream_username = "root"
  upstream_identity_file = "~/.ssh/id_rsa"
  upstream_hostname = "127.0.0.1" // change from host to upstream_hostname
  upstream_port = 22 // chanage from port to upstream_port
  // should we add another field for upstream auth type?
  upstream_authentication_type = "ssh_private_key"
}

resource "border0_socket" "ssh_with_certificate" {
  name = "ssh-with-certificate"
  socket_type = "ssh"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "ssh socket with certificate"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  // upstream_username = "root"
  upstream_hostname = "127.0.0.1" // change from host to upstream_hostname
  upstream_port = 22 // chanage from port to upstream_port
  // should we add another field for upstream auth type?
  upstream_authentication_type = "border0_cert"
}

resource "border0_socket" "ssh_with_ec2_instance_connect" {
  name = "ssh-with-ec2-instance-connect"
  socket_type = "ssh"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "ssh socket with ec2 instance connect"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  aws_ec2_target = "some_ec2_instance_id"
  upstream_hostname = "127.0.0.1" // change from host to upstream_hostname
  upstream_port = 22 // chanage from port to upstream_port
  // should we add another field for upstream auth type?
  upstream_authentication_type = "aws_ec2_instance_connect"
}

resource "border0_socket" "ssh_with_ssm" {
  name = "ssh-with-ssm"
  socket_type = "ssh"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "ssh socket with ssm"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  aws_ec2_target = "some_ec2_instance_id"
  // should we add another field for upstream auth type?
  upstream_authentication_type = "aws_ssm"
}

// ==========================================
// Database Sockets
// ==========================================

resource "border0_socket" "mysql" {
  name = "mysql"
  socket_type = "database"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "mysql database socket"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  upstream_username = "root" // change from upstream_user to upstream_username
  upstream_password = "password"
  upstream_tls = false // optional
  upstream_hostname = "my-rds-instance.cluster-giberish.us-east-2.rds.amazonaws.com" // change from host to upstream_hostname
  upstream_port = 3306 // change from port to upstream_port
  upstream_type = "mysql"
}

resource "border0_socket" "postgres" {
  name = "postgres"
  socket_type = "database"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "postgres database socket"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  upstream_username = "root" // change from upstream_user to upstream_username
  upstream_password = "password"
  upstream_tls = false // optional
  upstream_hostname = "my-rds-instance.cluster-giberish.us-east-2.rds.amazonaws.com" // change from host to upstream_hostname
  upstream_port = 5432 // change from port to upstream_port
  upstream_type = "postgres"
}

resource "border0_socket" "google_cloudsql_with_iam" {
  name = "google-cloudsql-with-iam"
  socket_type = "database"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "google cloudsql database socket with iam"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  upstream_username = "sql-263" // change from upstream_user to upstream_username
  upstream_type = "mysql"
  upstream_port = 3306 // change from port to upstream_port
  google_cloudsql_connector = true // was cloudsql_connector
  google_cloudsql_instance = "experiments-377315:europe-west4:mysqltest" // was cloudsql_instance
  google_credentials_file = "/path/to/credentials/experiments-377315-bef1eec1b496.json"
  // should we add another field for upstream auth type?
  upstream_authentication_type = "googe_cloudsql_iam" // was cloudsql_iam_auth
}

resource "border0_socket" "google_cloudsql_with_connector_and_username_password" {
  name = "google-cloudsql-with-connector-and-username-password"
  socket_type = "database"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "google cloudsql database socket with connector and username + password"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  upstream_username = "sql-263" // change from upstream_user to upstream_username
  upstream_password = "mysqltest"
  upstream_type = "mysql"
  upstream_port = 3306 // change from port to upstream_port
  google_cloudsql_connector = true // was cloudsql_connector
  google_cloudsql_instance = "experiments-377315:europe-west4:mysqltest" // was cloudsql_instance
  google_credentials_file = "/path/to/credentials/experiments-377315-bef1eec1b496.json"
  // should we add another field for upstream auth type?
  upstream_authentication_type = "password" // optional, default is "password"
}

resource "border0_socket" "aws_rds_with_iam" {
  name = "aws-rds-with-iam"
  socket_type = "database"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "aws rds database socket with iam"
  session_recording = true
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  upstream_username = "Border0ConnectorUser" // change from upstream_user to upstream_username
  upstream_hostname = "my-rds-instance.cluster-giberish.us-east-2.rds.amazonaws.com" // change from host to upstream_hostname
  upstream_port = 3306 // change from port to upstream_port
  upstream_type = "mysql"
  aws_region = "eu-central-1"
  // should we add another field for upstream auth type?
  upstream_authentication_type = "aws_rds_iam" // was rds_iam_auth
}

// ===========================================================================
// TLS Socket
// ===========================================================================

resource "border0_socket" "tls_vnc" {
  name = "tls-vnc"
  socket_type = "tls"
  connector_id = border0_connector.my_connector.id

  // optional fields
  description = "my vnc server"
  policy_ids = [border0_policy.my_policy.id, border0_policy.another_policy.id]
  tags = {
    env = "dev"
    team = "infra"
  }

  // optional and use case related
  upstream_hostname = "localhost" // change from host to upstream_hostname
  upstream_port = "5900" // change from port to upstream_port
}
