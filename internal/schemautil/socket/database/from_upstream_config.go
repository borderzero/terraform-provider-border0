package database

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FromUpstreamConfig(d *schema.ResourceData, config *service.DatabaseServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "database" but database service configuration was not present`)
	}

	data := map[string]any{
		"service_type": config.DatabaseServiceType,
	}

	var diags diag.Diagnostics

	switch config.DatabaseServiceType {
	case service.DatabaseServiceTypeStandard:
		diags = standardFromUpstreamConfig(&data, config.Standard)
	case service.DatabaseServiceTypeAwsRds:
		diags = awsRdsFromUpstreamConfig(&data, config.AwsRds)
	case service.DatabaseServiceTypeGcpCloudSql:
		diags = gcpCloudSqlFromUpstreamConfig(&data, config.GcpCloudSql)
	default:
		return diag.Errorf(`sockets with database service type "%s" not yet supported`, config.DatabaseServiceType)
	}

	if diags.HasError() {
		return diags
	}

	if err := d.Set("database_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "database_configuration"`)
	}

	return nil
}

func standardFromUpstreamConfig(data *map[string]any, config *service.StandardDatabaseServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with database service type "standard" but standard database service configuration was not present`)
	}

	(*data)["authentication_type"] = config.AuthenticationType
	(*data)["protocol"] = config.DatabaseProtocol
	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port

	switch config.AuthenticationType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "username_and_password" but username and password auth configuration was not present`)
		}
		(*data)["username"] = config.UsernameAndPasswordAuth.Username
		(*data)["password"] = config.UsernameAndPasswordAuth.Password
	case service.DatabaseAuthenticationTypeTls:
		if config.TlsAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "tls" but TLS auth configuration was not present`)
		}
		(*data)["username"] = config.TlsAuth.Username
		(*data)["password"] = config.TlsAuth.Password
		(*data)["certificate"] = config.TlsAuth.Certificate
		(*data)["private_key"] = config.TlsAuth.Key
		(*data)["ca_certificate"] = config.TlsAuth.CaCertificate
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, config.AuthenticationType)
	}

	return nil
}

func awsRdsFromUpstreamConfig(data *map[string]any, config *service.AwsRdsDatabaseServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with database service type "aws_rds" but AWS RDS database service configuration was not present`)
	}

	(*data)["authentication_type"] = config.AuthenticationType
	(*data)["protocol"] = config.DatabaseProtocol
	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port

	switch config.AuthenticationType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "username_and_password" but username and password auth configuration was not present`)
		}
		(*data)["username"] = config.UsernameAndPasswordAuth.Username
		(*data)["password"] = config.UsernameAndPasswordAuth.Password
		(*data)["ca_certificate"] = config.UsernameAndPasswordAuth.CaCertificate
	case service.DatabaseAuthenticationTypeIam:
		if config.IamAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "iam" but IAM auth configuration was not present`)
		}
		(*data)["username"] = config.IamAuth.Username
		(*data)["rds_instance_region"] = config.IamAuth.RdsInstanceRegion
		(*data)["ca_certificate"] = config.IamAuth.CaCertificate
		(*data)["aws_credentials"] = shared.FromAwsCredentials(config.IamAuth.AwsCredentials)
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, config.AuthenticationType)
	}

	return nil
}

func gcpCloudSqlFromUpstreamConfig(data *map[string]any, config *service.GcpCloudSqlDatabaseServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with database service type "gcp_cloudsql" but GCP Cloud SQL database service configuration was not present`)
	}

	if config.CloudSqlConnectorEnabled {
		return gcpCloudSqlConnectorFromUpstreamConfig(data, config.Connector)
	}
	return gcpCloudSqlStandardFromUpstreamConfig(data, config.Standard)
}

func gcpCloudSqlConnectorFromUpstreamConfig(data *map[string]any, config *service.GcpCloudSqlConnectorConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with database service type "gcp_cloudsql" and Cloud SQL Connector enabled but Cloud SQL Connector configuration was not present`)
	}

	(*data)["authentication_type"] = config.AuthenticationType
	(*data)["protocol"] = config.DatabaseProtocol

	switch config.AuthenticationType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "username_and_password" but username and password auth configuration was not present`)
		}
		(*data)["username"] = config.UsernameAndPasswordAuth.Username
		(*data)["password"] = config.UsernameAndPasswordAuth.Password
		(*data)["cloudsql_instance_id"] = config.UsernameAndPasswordAuth.InstanceId
		(*data)["gcp_credentials"] = config.UsernameAndPasswordAuth.GcpCredentialsJson
	case service.DatabaseAuthenticationTypeIam:
		if config.IamAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "iam" but IAM auth configuration was not present`)
		}
		(*data)["username"] = config.IamAuth.Username
		(*data)["cloudsql_instance_id"] = config.IamAuth.InstanceId
		(*data)["gcp_credentials"] = config.IamAuth.GcpCredentialsJson
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, config.AuthenticationType)
	}

	return nil
}

func gcpCloudSqlStandardFromUpstreamConfig(data *map[string]any, config *service.GcpCloudSqlStandardConfiguration) diag.Diagnostics {
	if data == nil {
		return diag.Errorf(`got a socket with database service type "gcp_cloudsql" but standard Cloud SQL database service configuration was not present`)
	}

	(*data)["authentication_type"] = config.AuthenticationType
	(*data)["protocol"] = config.DatabaseProtocol
	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port

	switch config.AuthenticationType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "username_and_password" but username and password auth configuration was not present`)
		}
		(*data)["username"] = config.UsernameAndPasswordAuth.Username
		(*data)["password"] = config.UsernameAndPasswordAuth.Password
	case service.DatabaseAuthenticationTypeTls:
		if config.TlsAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "tls" but TLS auth configuration was not present`)
		}
		(*data)["username"] = config.TlsAuth.Username
		(*data)["password"] = config.TlsAuth.Password
		(*data)["certificate"] = config.TlsAuth.Certificate
		(*data)["private_key"] = config.TlsAuth.Key
		(*data)["ca_certificate"] = config.TlsAuth.CaCertificate
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, config.AuthenticationType)
	}

	return nil
}
