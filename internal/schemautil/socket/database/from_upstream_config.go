package database

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's database service configuration into terraform resource data for
// the "database_configuration" attribute on the "border0_socket" resource.
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
	case service.DatabaseServiceTypeAzureSql:
		diags = azureSqlFromUpstreamConfig(&data, config.AzureSql)
	case service.DatabaseServiceTypeAwsDocumentDB:
		diags = awsDocumentDBFromUpstreamConfig(&data, config.AwsDocumentDB)
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
	(*data)["database_name"] = config.DatabaseName

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
	case service.DatabaseAuthenticationTypeNoAuth:
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
	(*data)["database_name"] = config.DatabaseName

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

func awsDocumentDBFromUpstreamConfig(data *map[string]any, config *service.AwsDocumentDBDatabaseServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with database service type "aws_documentdb" but AWS DocumentDB database service configuration was not present`)
	}

	(*data)["authentication_type"] = config.AuthenticationType
	(*data)["protocol"] = config.DatabaseProtocol
	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port
	(*data)["database_name"] = config.DatabaseName

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
		(*data)["cluster_region"] = config.IamAuth.ClusterRegion
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

	(*data)["protocol"] = config.DatabaseProtocol
	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port

	switch {
	case config.UsernameAndPasswordAuth != nil:
		if config.UsernameAndPasswordAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "username_and_password" but username and password auth configuration was not present`)
		}
		(*data)["username"] = config.UsernameAndPasswordAuth.Username
		(*data)["password"] = config.UsernameAndPasswordAuth.Password

	case config.TlsAuth != nil:
		if config.TlsAuth == nil {
			return diag.Errorf(`got a socket with database authentication type "tls" but TLS auth configuration was not present`)
		}
		(*data)["username"] = config.TlsAuth.Username
		(*data)["password"] = config.TlsAuth.Password
		(*data)["certificate"] = config.TlsAuth.Certificate
		(*data)["private_key"] = config.TlsAuth.Key
		if config.TlsAuth.CaCertificate != "" {
			(*data)["ca_certificate"] = config.TlsAuth.CaCertificate
		}
		(*data)["tls_auth"] = true

	case config.GcpCloudSQLConnectorAuth != nil:
		(*data)["username"] = config.GcpCloudSQLConnectorAuth.Username
		(*data)["password"] = config.GcpCloudSQLConnectorAuth.Password
		(*data)["cloudsql_instance_id"] = config.GcpCloudSQLConnectorAuth.InstanceId
		(*data)["gcp_credentials"] = config.GcpCloudSQLConnectorAuth.GcpCredentialsJson
		(*data)["cloudsql_connector_enabled"] = true

	case config.GcpCloudSQLConnectorIAMAuth != nil:
		(*data)["username"] = config.GcpCloudSQLConnectorIAMAuth.Username
		(*data)["cloudsql_instance_id"] = config.GcpCloudSQLConnectorIAMAuth.InstanceId
		(*data)["gcp_credentials"] = config.GcpCloudSQLConnectorIAMAuth.GcpCredentialsJson
		(*data)["cloudsql_iam_auth"] = true

	default:
		return diag.Errorf("gcp cloud sql database configuration had no authentication configured")
	}

	return nil
}

func azureSqlFromUpstreamConfig(data *map[string]any, config *service.AzureSqlDatabaseServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with database service type "azure_sql" but Azure SQL database service configuration was not present`)
	}

	(*data)["protocol"] = config.DatabaseProtocol
	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port

	switch {
	case config.AzureActiveDirectoryIntegrated != nil:
		(*data)["azure_ad_integrated"] = true

	case config.AzureActiveDirectoryPassword != nil:
		(*data)["username"] = config.AzureActiveDirectoryPassword.Username
		(*data)["password"] = config.AzureActiveDirectoryPassword.Password
		(*data)["azure_ad_auth"] = true

	case config.Kerberos != nil:
		(*data)["username"] = config.Kerberos.Username
		(*data)["password"] = config.Kerberos.Password
		(*data)["kerberos_auth"] = true

	case config.SqlAuthentication != nil:
		(*data)["username"] = config.SqlAuthentication.Username
		(*data)["password"] = config.SqlAuthentication.Password
		(*data)["sql_auth"] = true

	default:
		return diag.Errorf("gcp cloud sql database configuration had no authentication configured")
	}

	return nil
}
