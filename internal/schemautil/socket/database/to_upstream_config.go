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

	case service.DatabaseServiceTypeGcpCloudSql:
		if config.GcpCloudSql == nil {
			config.GcpCloudSql = new(service.GcpCloudSqlDatabaseServiceConfiguration)
		}
		return gcpCloudSqlToUpstreamConfig(data, config.GcpCloudSql)

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
			shared.ToAwsCredentials(v, config.IamAuth.AwsCredentials)
		}
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, authType)
	}

	return nil
}

func gcpCloudSqlToUpstreamConfig(data map[string]any, config *service.GcpCloudSqlDatabaseServiceConfiguration) diag.Diagnostics {
	var connectorEnabled bool // default to false

	if v, ok := data["cloudsql_connector_enabled"]; ok {
		connectorEnabled = v.(bool)
	}
	config.CloudSqlConnectorEnabled = connectorEnabled

	if connectorEnabled {
		if config.Connector == nil {
			config.Connector = new(service.GcpCloudSqlConnectorConfiguration)
		}
		return gcpCloudSqlConnectorToUpstreamConfig(data, config.Connector)
	}

	if config.Standard == nil {
		config.Standard = new(service.GcpCloudSqlStandardConfiguration)
	}
	return gcpCloudSqlStandardToUpstreamConfig(data, config.Standard)
}

func gcpCloudSqlConnectorToUpstreamConfig(data map[string]any, config *service.GcpCloudSqlConnectorConfiguration) diag.Diagnostics {
	authType := service.DatabaseAuthenticationTypeUsernameAndPassword // default to "username_and_password"

	if v, ok := data["authentication_type"]; ok {
		authType = v.(string)
	}
	config.AuthenticationType = authType

	if v, ok := data["protocol"]; ok {
		config.DatabaseProtocol = v.(string)
	}

	switch authType {
	case service.DatabaseAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuth == nil {
			config.UsernameAndPasswordAuth = new(service.GcpCloudSqlUsernameAndPasswordAuthConfiguration)
		}

		if v, ok := data["username"]; ok {
			config.UsernameAndPasswordAuth.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.UsernameAndPasswordAuth.Password = v.(string)
		}
		if v, ok := data["cloudsql_instance_id"]; ok {
			config.UsernameAndPasswordAuth.InstanceId = v.(string)
		}
		if v, ok := data["gcp_credentials"]; ok {
			config.UsernameAndPasswordAuth.GcpCredentialsJson = v.(string)
		}
	case service.DatabaseAuthenticationTypeIam:
		if config.IamAuth == nil {
			config.IamAuth = new(service.GcpCloudSqlIamAuthConfiguration)
		}

		if v, ok := data["username"]; ok {
			config.IamAuth.Username = v.(string)
		}
		if v, ok := data["cloudsql_instance_id"]; ok {
			config.IamAuth.InstanceId = v.(string)
		}
		if v, ok := data["gcp_credentials"]; ok {
			config.IamAuth.GcpCredentialsJson = v.(string)
		}
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, authType)
	}

	return nil
}

func gcpCloudSqlStandardToUpstreamConfig(data map[string]any, config *service.GcpCloudSqlStandardConfiguration) diag.Diagnostics {
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
	default:
		return diag.Errorf(`database authentication type "%s" is invalid`, authType)
	}

	return nil
}
