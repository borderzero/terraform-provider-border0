package schemautil

import (
	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig translates upstream service
// config to terraform resource schema for sockets.
func FromUpstreamConfig(
	d *schema.ResourceData,
	configs *border0client.SocketUpstreamConfigs,
) diag.Diagnostics {
	// no-op if upstream config is not set
	if len(configs.List) == 0 {
		return nil
	}

	firstConfig := configs.List[0]
	config := firstConfig.Config

	switch config.ServiceType {
	case service.ServiceTypeSsh:
		if config.SshServiceConfiguration != nil {
			return sshFromUpstreamConfig(d, config.SshServiceConfiguration)
		}
		return diag.Errorf("got a socket with service type ssh but ssh service configuration was not present")

	// TODO: support ServiceTypeDatabase, ServiceTypeHttp, ServiceTypeTls
	default:
		return diag.Errorf(`sockets with service type "%s" not yet supported`, config.ServiceType)
	}
}

func sshFromUpstreamConfig(
	d *schema.ResourceData,
	config *service.SshServiceConfiguration,
) diag.Diagnostics {
	if diagnostics := SetValues(d, map[string]any{
		"upstream_service_type": config.SshServiceType,
	}); diagnostics.HasError() {
		return diagnostics
	}

	switch config.SshServiceType {
	case service.SshServiceTypeStandard:
		if config.StandardSshServiceConfiguration != nil {
			return standardSshFromUpstreamConfig(d, config.StandardSshServiceConfiguration)
		}
		return diag.Errorf("got a socket with ssh service type standard but standard ssh service configuration was not present")

	// TODO: support SshServiceTypeAwsEc2InstanceConnect, SshServiceTypeAwsSsm, SshServiceTypeConnectorBuiltIn
	default:
		return diag.Errorf(`sockets with ssh service type "%s" not yet supported`, config.SshServiceType)
	}
}

func standardSshFromUpstreamConfig(
	d *schema.ResourceData,
	config *service.StandardSshServiceConfiguration,
) diag.Diagnostics {

	if diagnostics := SetValues(d, map[string]any{
		"upstream_hostname":            config.Hostname,
		"upstream_port":                config.Port,
		"upstream_authentication_type": config.SshAuthenticationType,
	}); diagnostics.HasError() {
		return diagnostics
	}

	switch config.SshAuthenticationType {
	case service.StandardSshServiceAuthenticationTypeBorder0Certificate:
		if config.Border0CertificateAuthConfiguration != nil {
			return border0CertificateStandardSshFromUpstreamConfig(d, config.Border0CertificateAuthConfiguration)
		}
		return diag.Errorf("got a standard ssh socket with authentication type border0 certificate but border0 certificate auth configuration was not present")
	case service.StandardSshServiceAuthenticationTypePrivateKey:
		if config.PrivateKeyAuthConfiguration != nil {
			return privateKeyStandardSshFromUpstreamConfig(d, config.PrivateKeyAuthConfiguration)
		}
		return diag.Errorf("got a standard ssh socket with authentication type border0 certificate but border0 certificate auth configuration was not present")
	case service.StandardSshServiceAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuthConfiguration != nil {
			return usernameAndPasswordStandardSshFromUpstreamConfig(d, config.UsernameAndPasswordAuthConfiguration)
		}
		return diag.Errorf("got a standard ssh socket with authentication type border0 certificate but border0 certificate auth configuration was not present")
	default:
		return diag.Errorf(`ssh authentication type "%s" is invalid`, config.SshAuthenticationType)
	}
}

func border0CertificateStandardSshFromUpstreamConfig(
	d *schema.ResourceData,
	config *service.Border0CertificateAuthConfiguration,
) diag.Diagnostics {
	// TODO: support username_provider
	return SetValues(d, map[string]any{
		"upstream_username": config.Username,
	})
}

func privateKeyStandardSshFromUpstreamConfig(
	d *schema.ResourceData,
	config *service.PrivateKeyAuthConfiguration,
) diag.Diagnostics {
	// TODO: support username_provider
	return SetValues(d, map[string]any{
		"upstream_username":    config.Username,
		"upstream_private_key": config.PrivateKey,
	})
}

func usernameAndPasswordStandardSshFromUpstreamConfig(
	d *schema.ResourceData,
	config *service.UsernameAndPasswordAuthConfiguration,
) diag.Diagnostics {
	// TODO: support username_provider
	return SetValues(d, map[string]any{
		"upstream_username": config.Username,
		"upstream_password": config.Password,
	})
}

// ToUpstreamConfig translates terraform resource schema
func ToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	// noop if connector id is not set or empty
	if v, ok := d.GetOk("connector_id"); !ok || v.(string) == "" {
		return nil
	}

	// if connector id is given, upstream config is required
	if socket.UpstreamConfig == nil {
		socket.UpstreamConfig = new(service.Configuration)
	}
	socket.UpstreamConfig.ServiceType = socket.SocketType

	var diagnostics diag.Diagnostics

	switch socket.UpstreamConfig.ServiceType {
	case service.ServiceTypeSsh:
		diagnostics = sshToUpstreamConfig(d, socket)

	// TODO: support ServiceTypeDatabase, ServiceTypeHttp, ServiceTypeTls
	default:
		return diag.Errorf(`sockets with service type "%s" not yet supported`, socket.UpstreamConfig.ServiceType)
	}

	if diagnostics.HasError() {
		return diagnostics
	}

	if err := socket.UpstreamConfig.Validate(); err != nil {
		return DiagnosticsError(err, "Upstream configuration is invalid")
	}

	return nil
}

func sshToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if socket.UpstreamConfig.SshServiceConfiguration == nil {
		socket.UpstreamConfig.SshServiceConfiguration = new(service.SshServiceConfiguration)
	}

	var sshServiceType string
	if v, ok := d.GetOk("upstream_service_type"); ok {
		sshServiceType = v.(string)
	}
	if socket.UpstreamConfig.SshServiceConfiguration.SshServiceType == "" {
		sshServiceType = service.SshServiceTypeStandard
	}
	socket.UpstreamConfig.SshServiceConfiguration.SshServiceType = sshServiceType

	switch sshServiceType {
	case service.SshServiceTypeStandard:
		return sshStandardToUpstreamConfig(d, socket)

	// TODO: support SshServiceTypeAwsEc2InstanceConnect, SshServiceTypeAwsSsm, SshServiceTypeConnectorBuiltIn
	default:
		return diag.Errorf(`sockets with ssh service type "%s" not yet supported`, sshServiceType)
	}
}

func sshStandardToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration == nil {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration = new(service.StandardSshServiceConfiguration)
	}

	var authType string
	if v, ok := d.GetOk("upstream_authentication_type"); ok {
		authType = v.(string)
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.SshAuthenticationType = authType
	}
	if v, ok := d.GetOk("upstream_hostname"); ok {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.Hostname = v.(string)
	}
	if v, ok := d.GetOk("upstream_port"); ok {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.Port = uint16(v.(int))
	}

	switch authType {
	case service.StandardSshServiceAuthenticationTypeBorder0Certificate:
		return border0CertificateStandardSshToUpstreamConfig(d, socket)
	case service.StandardSshServiceAuthenticationTypePrivateKey:
		return privateKeyStandardSshToUpstreamConfig(d, socket)
	case service.StandardSshServiceAuthenticationTypeUsernameAndPassword:
		return usernameAndPasswordStandardSshToUpstreamConfig(d, socket)
	default:
		return diag.Errorf(`ssh authentication type "%s" is invalid`, authType)
	}
}

func border0CertificateStandardSshToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.Border0CertificateAuthConfiguration == nil {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.Border0CertificateAuthConfiguration = new(service.Border0CertificateAuthConfiguration)
	}
	// TODO: support username_provider
	if v, ok := d.GetOk("upstream_username"); ok {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.Border0CertificateAuthConfiguration.Username = v.(string)
	}
	return nil
}

func privateKeyStandardSshToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.PrivateKeyAuthConfiguration == nil {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.PrivateKeyAuthConfiguration = new(service.PrivateKeyAuthConfiguration)
	}
	// TODO: support username_provider
	if v, ok := d.GetOk("upstream_username"); ok {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.PrivateKeyAuthConfiguration.Username = v.(string)
	}
	if v, ok := d.GetOk("upstream_private_key"); ok {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.PrivateKeyAuthConfiguration.PrivateKey = v.(string)
	}
	return nil
}

func usernameAndPasswordStandardSshToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.UsernameAndPasswordAuthConfiguration == nil {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.UsernameAndPasswordAuthConfiguration = new(service.UsernameAndPasswordAuthConfiguration)
	}
	// TODO: support username_provider
	if v, ok := d.GetOk("upstream_username"); ok {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.UsernameAndPasswordAuthConfiguration.Username = v.(string)
	}
	if v, ok := d.GetOk("upstream_password"); ok {
		socket.UpstreamConfig.SshServiceConfiguration.StandardSshServiceConfiguration.UsernameAndPasswordAuthConfiguration.Password = v.(string)
	}
	return nil
}
