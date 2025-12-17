package database

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "database_configuration" attribute into
// a socket's database service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.DatabaseServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("database_configuration"); ok {
		if databaseConfigsList := v.([]any); len(databaseConfigsList) > 0 {
			data = databaseConfigsList[0].(map[string]any)
		}
	}

	databaseServiceType := service.DatabaseServiceTypeStandard // default to "standard"

	if v, ok := data["service_type"]; ok {
		databaseServiceType = v.(string)
	}
	config.DatabaseServiceType = databaseServiceType

	switch databaseServiceType {
	case service.DatabaseServiceTypeStandard:
		if config.Standard == nil {
			config.Standard = new(service.StandardDatabaseServiceConfiguration)
		}
		return standardToUpstreamConfig(data, config.Standard)

	case service.DatabaseServiceTypeAwsRds:
		if config.AwsRds == nil {
			config.AwsRds = new(service.AwsRdsDatabaseServiceConfiguration)
		}
		return awsRdsToUpstreamConfig(data, config.AwsRds)

	case service.DatabaseServiceTypeAwsDocumentDB:
		if config.AwsDocumentDB == nil {
			config.AwsDocumentDB = new(service.AwsDocumentDBDatabaseServiceConfiguration)
		}
		return awsDocumentDBToUpstreamConfig(data, config.AwsDocumentDB)

	case service.DatabaseServiceTypeGcpCloudSql:
		if config.GcpCloudSql == nil {
			config.GcpCloudSql = new(service.GcpCloudSqlDatabaseServiceConfiguration)
		}
		return gcpCloudSqlToUpstreamConfig(data, config.GcpCloudSql)

	case service.DatabaseServiceTypeAzureSql:
		if config.AzureSql == nil {
			config.AzureSql = new(service.AzureSqlDatabaseServiceConfiguration)
		}
		return azureSqlToUpstreamConfig(data, config.AzureSql)

	case service.DatabaseServiceTypeMongoDBAtlas:
		if config.MongoDBAtlas == nil {
			config.MongoDBAtlas = new(service.MongoDBAtlasDatabaseServiceConfiguration)
		}
		return mongoDBAtlasToUpstreamConfig(data, config.MongoDBAtlas)

	default:
		return diag.Errorf(`sockets with database service type "%s" not yet supported`, databaseServiceType)
	}
}

func standardToUpstreamConfig(data map[string]any, config *service.StandardDatabaseServiceConfiguration) diag.Diagnostics {
	authType := service.DatabaseAuthenticationTypeUsernameAndPassword // default to "username_and_password"

	if v, ok := data["authentication_type"]; ok {
		authType = v.(string)
	}
	config.AuthenticationType = authType

	if v, ok := data["protocol"]; ok {
		config.DatabaseProtocol = v.(string)
	}
	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}
	if v, ok := data["database_name"]; ok {
		config.DatabaseName = v.(string)
	}

	switch authType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			config.UsernameAndPasswordAuth = new(service.DatabaseUsernameAndPasswordAuthConfiguration)
		}

		if v, ok := data["username"]; ok {
			config.UsernameAndPasswordAuth.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.UsernameAndPasswordAuth.Password = v.(string)
		}
	case service.DatabaseAuthenticationTypeTls:
		if config.TlsAuth == nil {
			config.TlsAuth = new(service.DatabaseTlsAuthConfiguration)
		}

		if v, ok := data["username"]; ok {
			config.TlsAuth.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.TlsAuth.Password = v.(string)
		}
		if v, ok := data["certificate"]; ok {
			config.TlsAuth.Certificate = v.(string)
		}
		if v, ok := data["private_key"]; ok {
			config.TlsAuth.Key = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			config.TlsAuth.CaCertificate = v.(string)
		}
	case service.DatabaseAuthenticationTypeNoAuth:
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, authType)
	}

	return nil
}

func awsRdsToUpstreamConfig(data map[string]any, config *service.AwsRdsDatabaseServiceConfiguration) diag.Diagnostics {
	authType := service.DatabaseAuthenticationTypeUsernameAndPassword // default to "username_and_password"

	if v, ok := data["authentication_type"]; ok {
		authType = v.(string)
	}
	config.AuthenticationType = authType

	if v, ok := data["protocol"]; ok {
		config.DatabaseProtocol = v.(string)
	}
	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}
	if v, ok := data["database_name"]; ok {
		config.DatabaseName = v.(string)
	}

	switch authType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			config.UsernameAndPasswordAuth = new(service.AwsRdsUsernameAndPasswordAuthConfiguration)
		}

		if v, ok := data["username"]; ok {
			config.UsernameAndPasswordAuth.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.UsernameAndPasswordAuth.Password = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			config.UsernameAndPasswordAuth.CaCertificate = v.(string)
		}
	case service.DatabaseAuthenticationTypeIam:
		if config.IamAuth == nil {
			config.IamAuth = new(service.AwsRdsIamAuthConfiguration)
		}

		if v, ok := data["username"]; ok {
			config.IamAuth.Username = v.(string)
		}
		if v, ok := data["rds_instance_region"]; ok {
			config.IamAuth.RdsInstanceRegion = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			config.IamAuth.CaCertificate = v.(string)
		}
		if v, ok := data["aws_credentials"]; ok {
			config.IamAuth.AwsCredentials = shared.ToAwsCredentials(v)
		}
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, authType)
	}

	return nil
}

func awsDocumentDBToUpstreamConfig(data map[string]any, config *service.AwsDocumentDBDatabaseServiceConfiguration) diag.Diagnostics {
	authType := service.DatabaseAuthenticationTypeUsernameAndPassword // default to "username_and_password"

	if v, ok := data["authentication_type"]; ok {
		authType = v.(string)
	}
	config.AuthenticationType = authType

	if v, ok := data["protocol"]; ok {
		config.DatabaseProtocol = v.(string)
	}
	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}
	if v, ok := data["database_name"]; ok {
		config.DatabaseName = v.(string)
	}

	switch authType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			config.UsernameAndPasswordAuth = new(service.UsernamePasswordCaAuthConfiguration)
		}

		if v, ok := data["username"]; ok {
			config.UsernameAndPasswordAuth.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.UsernameAndPasswordAuth.Password = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			config.UsernameAndPasswordAuth.CaCertificate = v.(string)
		}
	case service.DatabaseAuthenticationTypeIam:
		if config.IamAuth == nil {
			config.IamAuth = new(service.MongoAWSAuthConfiguration)
		}

		if v, ok := data["cluster_region"]; ok {
			config.IamAuth.ClusterRegion = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			config.IamAuth.CaCertificate = v.(string)
		}
		if v, ok := data["aws_credentials"]; ok {
			config.IamAuth.AwsCredentials = shared.ToAwsCredentials(v)
		}
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, authType)
	}

	return nil
}

func mongoDBAtlasToUpstreamConfig(data map[string]any, config *service.MongoDBAtlasDatabaseServiceConfiguration) diag.Diagnostics {
	authType := service.DatabaseAuthenticationTypeUsernameAndPassword // default to "username_and_password"

	if v, ok := data["authentication_type"]; ok {
		authType = v.(string)
	}
	config.AuthenticationType = authType

	// MongoDB Atlas uses mongodb protocol
	config.DatabaseProtocol = "mongodb"

	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	// Note: MongoDB Atlas uses SRV records, no port needed
	if v, ok := data["database_name"]; ok {
		config.DatabaseName = v.(string)
	}

	switch authType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			config.UsernameAndPasswordAuth = new(service.UsernamePasswordCaAuthConfiguration)
		}

		if v, ok := data["username"]; ok {
			config.UsernameAndPasswordAuth.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.UsernameAndPasswordAuth.Password = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			config.UsernameAndPasswordAuth.CaCertificate = v.(string)
		}
	case service.DatabaseAuthenticationTypeIam:
		if config.IamAuth == nil {
			config.IamAuth = new(service.MongoAWSAuthConfiguration)
		}

		if v, ok := data["cluster_region"]; ok {
			config.IamAuth.ClusterRegion = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			config.IamAuth.CaCertificate = v.(string)
		}
		if v, ok := data["aws_credentials"]; ok {
			config.IamAuth.AwsCredentials = shared.ToAwsCredentials(v)
		}
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, authType)
	}

	return nil
}

func gcpCloudSqlToUpstreamConfig(data map[string]any, config *service.GcpCloudSqlDatabaseServiceConfiguration) diag.Diagnostics {
	tlsAuthEnabled := false
	cloudSqlConnectorEnabled := false
	cloudSqlIamEnabled := false

	if v, ok := data["tls_auth"]; ok {
		tlsAuthEnabled = v.(bool)
	}
	if v, ok := data["cloudsql_connector_enabled"]; ok {
		cloudSqlConnectorEnabled = v.(bool)
	}
	if v, ok := data["cloudsql_iam_auth"]; ok {
		cloudSqlIamEnabled = v.(bool)
	}

	if tlsAuthEnabled {
		if config.TlsAuth == nil {
			config.TlsAuth = new(service.DatabaseTlsAuthConfiguration)
		}
		if v, ok := data["username"]; ok {
			config.TlsAuth.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.TlsAuth.Password = v.(string)
		}
		if v, ok := data["certificate"]; ok {
			config.TlsAuth.Certificate = v.(string)
		}
		if v, ok := data["private_key"]; ok {
			config.TlsAuth.Key = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			if ca_cert := v.(string); ca_cert != "" {
				config.TlsAuth.CaCertificate = ca_cert
			}
		}
		return nil
	}

	if cloudSqlConnectorEnabled {
		if config.GcpCloudSQLConnectorAuth == nil {
			config.GcpCloudSQLConnectorAuth = new(service.GcpCloudSqlConnectorAuthConfiguration)
		}
		if v, ok := data["username"]; ok {
			config.GcpCloudSQLConnectorAuth.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.GcpCloudSQLConnectorAuth.Password = v.(string)
		}
		if v, ok := data["cloudsql_instance_id"]; ok {
			config.GcpCloudSQLConnectorAuth.InstanceId = v.(string)
		}
		if v, ok := data["gcp_credentials"]; ok {
			config.GcpCloudSQLConnectorAuth.GcpCredentialsJson = v.(string)
		}
		return nil
	}

	if cloudSqlIamEnabled {
		if config.GcpCloudSQLConnectorIAMAuth == nil {
			config.GcpCloudSQLConnectorIAMAuth = new(service.GcpCloudSqlConnectorIamAuthConfiguration)
		}
		if v, ok := data["username"]; ok {
			config.GcpCloudSQLConnectorIAMAuth.Username = v.(string)
		}
		if v, ok := data["certificate"]; ok {
			config.TlsAuth.Certificate = v.(string)
		}
		if v, ok := data["private_key"]; ok {
			config.TlsAuth.Key = v.(string)
		}
		if v, ok := data["ca_certificate"]; ok {
			if ca_cert := v.(string); ca_cert != "" {
				config.TlsAuth.CaCertificate = ca_cert
			}
		}
		return nil
	}

	if config.UsernameAndPasswordAuth == nil {
		config.UsernameAndPasswordAuth = new(service.DatabaseUsernameAndPasswordAuthConfiguration)
	}
	if v, ok := data["username"]; ok {
		config.UsernameAndPasswordAuth.Username = v.(string)
	}
	if v, ok := data["password"]; ok {
		config.UsernameAndPasswordAuth.Password = v.(string)
	}
	return nil
}

func azureSqlToUpstreamConfig(data map[string]any, config *service.AzureSqlDatabaseServiceConfiguration) diag.Diagnostics {
	if v, ok := data["protocol"]; ok {
		config.DatabaseProtocol = v.(string)
	}
	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}

	integratedAuth := false
	azureAdAuth := false
	kerberosAuth := false
	sqlAuth := false

	if v, ok := data["azure_ad_integrated"]; ok {
		integratedAuth = v.(bool)
	}
	if v, ok := data["azure_ad_auth"]; ok {
		azureAdAuth = v.(bool)
	}
	if v, ok := data["kerberos_auth"]; ok {
		kerberosAuth = v.(bool)
	}
	if v, ok := data["sql_auth"]; ok {
		sqlAuth = v.(bool)
	}

	switch {
	case integratedAuth:
		if config.AzureActiveDirectoryIntegrated == nil {
			config.AzureActiveDirectoryIntegrated = &struct{}{}
		}

	case azureAdAuth:
		if config.AzureActiveDirectoryPassword == nil {
			config.AzureActiveDirectoryPassword = new(service.DatabaseUsernameAndPasswordAuthConfiguration)
		}
		if v, ok := data["username"]; ok {
			config.AzureActiveDirectoryPassword.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.AzureActiveDirectoryPassword.Password = v.(string)
		}

	case kerberosAuth:
		if config.Kerberos == nil {
			config.Kerberos = new(service.DatabaseKerberosAuthConfiguration)
		}
		if v, ok := data["username"]; ok {
			config.Kerberos.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.Kerberos.Password = v.(string)
		}

	case sqlAuth:
		if config.SqlAuthentication == nil {
			config.SqlAuthentication = new(service.DatabaseSqlAuthConfiguration)
		}
		if v, ok := data["username"]; ok {
			config.SqlAuthentication.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.SqlAuthentication.Password = v.(string)
		}

	default:
		return diag.Errorf(`database with no authentication configuration`)
	}

	return nil
}
