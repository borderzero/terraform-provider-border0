package ssh

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "ssh_configuration" attribute into
// a socket's SSH service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.SshServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("ssh_configuration"); ok {
		if sshConfigsList := v.([]any); len(sshConfigsList) > 0 {
			data = sshConfigsList[0].(map[string]any)
		}
	}

	sshServiceType := service.SshServiceTypeStandard // default to "standard"

	if v, ok := data["service_type"]; ok {
		sshServiceType = v.(string)
	}
	config.SshServiceType = sshServiceType

	switch sshServiceType {
	case service.SshServiceTypeStandard:
		if config.StandardSshServiceConfiguration == nil {
			config.StandardSshServiceConfiguration = new(service.StandardSshServiceConfiguration)
		}
		return standardToUpstreamConfig(data, config.StandardSshServiceConfiguration)

	case service.SshServiceTypeAwsEc2InstanceConnect:
		if config.AwsEc2ICSshServiceConfiguration == nil {
			config.AwsEc2ICSshServiceConfiguration = new(service.AwsEc2ICSshServiceConfiguration)
		}
		return awsEc2InstanceConnectToUpstreamConfig(data, config.AwsEc2ICSshServiceConfiguration)

	case service.SshServiceTypeAwsSsm:
		if config.AwsSsmSshServiceConfiguration == nil {
			config.AwsSsmSshServiceConfiguration = new(service.AwsSsmSshServiceConfiguration)
		}
		return awsSsmToUpstreamConfig(data, config.AwsSsmSshServiceConfiguration)

	case service.SshServiceTypeConnectorBuiltIn:
		if config.BuiltInSshServiceConfiguration == nil {
			config.BuiltInSshServiceConfiguration = new(service.BuiltInSshServiceConfiguration)
		}
		return connectorBuiltInToUpstreamConfig(data, config.BuiltInSshServiceConfiguration)

	default:
		return diag.Errorf(`sockets with ssh service type "%s" not yet supported`, sshServiceType)
	}
}

func standardToUpstreamConfig(data map[string]any, config *service.StandardSshServiceConfiguration) diag.Diagnostics {
	authType := service.StandardSshServiceAuthenticationTypeBorder0Certificate // default to "border0_certificate"
	usernameProvider := service.UsernameProviderPromptClient                   // default to "prompt_client"

	if v, ok := data["authentication_type"]; ok {
		authType = v.(string)
	}
	config.SshAuthenticationType = authType

	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}

	switch authType {
	case service.StandardSshServiceAuthenticationTypeBorder0Certificate:
		if config.Border0CertificateAuthConfiguration == nil {
			config.Border0CertificateAuthConfiguration = new(service.Border0CertificateAuthConfiguration)
		}

		if v, ok := data["username_provider"]; ok {
			usernameProvider = v.(string)
		}
		config.Border0CertificateAuthConfiguration.UsernameProvider = usernameProvider

		if v, ok := data["username"]; ok {
			config.Border0CertificateAuthConfiguration.Username = v.(string)
		}
	case service.StandardSshServiceAuthenticationTypePrivateKey:
		if config.PrivateKeyAuthConfiguration == nil {
			config.PrivateKeyAuthConfiguration = new(service.PrivateKeyAuthConfiguration)
		}

		if v, ok := data["username_provider"]; ok {
			usernameProvider = v.(string)
		}
		config.PrivateKeyAuthConfiguration.UsernameProvider = usernameProvider

		if v, ok := data["username"]; ok {
			config.PrivateKeyAuthConfiguration.Username = v.(string)
		}
		if v, ok := data["private_key"]; ok {
			config.PrivateKeyAuthConfiguration.PrivateKey = v.(string)
		}
	case service.StandardSshServiceAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuthConfiguration == nil {
			config.UsernameAndPasswordAuthConfiguration = new(service.UsernameAndPasswordAuthConfiguration)
		}

		if v, ok := data["username_provider"]; ok {
			usernameProvider = v.(string)
		}
		config.UsernameAndPasswordAuthConfiguration.UsernameProvider = usernameProvider

		if v, ok := data["username"]; ok {
			config.UsernameAndPasswordAuthConfiguration.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.UsernameAndPasswordAuthConfiguration.Password = v.(string)
		}
	default:
		return diag.Errorf(`ssh authentication type "%s" is invalid`, authType)
	}

	return nil
}

func awsEc2InstanceConnectToUpstreamConfig(data map[string]any, config *service.AwsEc2ICSshServiceConfiguration) diag.Diagnostics {
	usernameProvider := service.UsernameProviderPromptClient // default to "prompt_client"

	if v, ok := data["username_provider"]; ok {
		usernameProvider = v.(string)
	}
	config.UsernameProvider = usernameProvider

	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}
	if v, ok := data["username"]; ok {
		config.Username = v.(string)
	}
	if v, ok := data["ec2_instance_id"]; ok {
		config.Ec2InstanceId = v.(string)
	}
	if v, ok := data["ec2_instance_region"]; ok {
		config.Ec2InstanceRegion = v.(string)
	}
	if v, ok := data["aws_credentials"]; ok {
		config.AwsCredentials = shared.ToAwsCredentials(v)
	}

	return nil
}

func awsSsmToUpstreamConfig(data map[string]any, config *service.AwsSsmSshServiceConfiguration) diag.Diagnostics {
	targetType := service.SsmTargetTypeEc2 // default to "ec2"

	if v, ok := data["ssm_target_type"]; ok {
		targetType = v.(string)
	}
	config.SsmTargetType = targetType

	switch targetType {
	case service.SsmTargetTypeEc2:
		if config.AwsSsmEc2TargetConfiguration == nil {
			config.AwsSsmEc2TargetConfiguration = new(service.AwsSsmEc2TargetConfiguration)
		}
		if v, ok := data["ec2_instance_id"]; ok {
			config.AwsSsmEc2TargetConfiguration.Ec2InstanceId = v.(string)
		}
		if v, ok := data["ec2_instance_region"]; ok {
			config.AwsSsmEc2TargetConfiguration.Ec2InstanceRegion = v.(string)
		}
		if v, ok := data["aws_credentials"]; ok {
			config.AwsSsmEc2TargetConfiguration.AwsCredentials = shared.ToAwsCredentials(v)
		}
	case service.SsmTargetTypeEcs:
		if config.AwsSsmEcsTargetConfiguration == nil {
			config.AwsSsmEcsTargetConfiguration = new(service.AwsSsmEcsTargetConfiguration)
		}
		if v, ok := data["ecs_cluster_region"]; ok {
			config.AwsSsmEcsTargetConfiguration.EcsClusterRegion = v.(string)
		}
		if v, ok := data["ecs_cluster_name"]; ok {
			config.AwsSsmEcsTargetConfiguration.EcsClusterName = v.(string)
		}
		if v, ok := data["ecs_service_name"]; ok {
			config.AwsSsmEcsTargetConfiguration.EcsServiceName = v.(string)
		}
		if v, ok := data["aws_credentials"]; ok {
			config.AwsSsmEcsTargetConfiguration.AwsCredentials = shared.ToAwsCredentials(v)
		}
	default:
		return diag.Errorf(`ssm target type "%s" is invalid`, targetType)
	}

	return nil
}

func connectorBuiltInToUpstreamConfig(data map[string]any, config *service.BuiltInSshServiceConfiguration) diag.Diagnostics {
	usernameProvider := service.UsernameProviderPromptClient // default to "prompt_client"

	if v, ok := data["username_provider"]; ok {
		usernameProvider = v.(string)
	}
	config.UsernameProvider = usernameProvider

	if v, ok := data["username"]; ok {
		config.Username = v.(string)
	}

	return nil
}
