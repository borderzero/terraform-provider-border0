package ssh

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's SSH service configuration into terraform resource data for
// the "ssh_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.SshServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "ssh" but SSH service configuration was not present`)
	}

	data := map[string]any{
		"service_type": config.SshServiceType,
	}

	var diags diag.Diagnostics

	switch config.SshServiceType {
	case service.SshServiceTypeStandard:
		diags = standardFromUpstreamConfig(&data, config.StandardSshServiceConfiguration)
	case service.SshServiceTypeAwsEc2InstanceConnect:
		diags = awsEc2InstanceConnectFromUpstreamConfig(&data, config.AwsEc2ICSshServiceConfiguration)
	case service.SshServiceTypeAwsSsm:
		diags = awsSsmFromUpstreamConfig(&data, config.AwsSsmSshServiceConfiguration)
	case service.SshServiceTypeConnectorBuiltIn:
		diags = connectorBuiltInFromUpstreamConfig(&data, config.BuiltInSshServiceConfiguration)
	default:
		return diag.Errorf(`sockets with SSH service type "%s" not yet supported`, config.SshServiceType)
	}

	if diags.HasError() {
		return diags
	}

	if err := d.Set("ssh_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "ssh_configuration"`)
	}

	return nil
}

func standardFromUpstreamConfig(data *map[string]any, config *service.StandardSshServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with SSH service type "standard" but standard SSH service configuration was not present`)
	}

	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port
	(*data)["authentication_type"] = config.SshAuthenticationType

	switch config.SshAuthenticationType {
	case service.StandardSshServiceAuthenticationTypeBorder0Certificate:
		if config.Border0CertificateAuthConfiguration == nil {
			return diag.Errorf(`got a socket with SSH authentication type "border0_certificate" but Border0 certificate auth configuration was not present`)
		}
		(*data)["username_provider"] = config.Border0CertificateAuthConfiguration.UsernameProvider
		(*data)["username"] = config.Border0CertificateAuthConfiguration.Username
	case service.StandardSshServiceAuthenticationTypePrivateKey:
		if config.PrivateKeyAuthConfiguration == nil {
			return diag.Errorf(`got a socket with SSH authentication type "private_key" but private key auth configuration was not present`)
		}
		(*data)["username_provider"] = config.PrivateKeyAuthConfiguration.UsernameProvider
		(*data)["username"] = config.PrivateKeyAuthConfiguration.Username
		(*data)["private_key"] = config.PrivateKeyAuthConfiguration.PrivateKey
	case service.StandardSshServiceAuthenticationTypeUsernameAndPassword:
		if config.UsernameAndPasswordAuthConfiguration == nil {
			return diag.Errorf(`got a socket with SSH authentication type "username_and_password" but username and password auth configuration was not present`)
		}
		(*data)["username_provider"] = config.UsernameAndPasswordAuthConfiguration.UsernameProvider
		(*data)["username"] = config.UsernameAndPasswordAuthConfiguration.Username
		(*data)["password"] = config.UsernameAndPasswordAuthConfiguration.Password
	default:
		return diag.Errorf(`SSH authentication type "%s" is invalid`, config.SshAuthenticationType)
	}

	return nil
}

func awsEc2InstanceConnectFromUpstreamConfig(data *map[string]any, config *service.AwsEc2ICSshServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with SSH service type "aws_ec2_instance_connect" but AWS EC2 Instance Connect SSH service configuration was not present`)
	}

	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port
	(*data)["username_provider"] = config.UsernameProvider
	(*data)["username"] = config.Username
	(*data)["ec2_instance_id"] = config.Ec2InstanceId
	(*data)["ec2_instance_region"] = config.Ec2InstanceRegion
	(*data)["aws_credentials"] = shared.FromAwsCredentials(config.AwsCredentials)

	return nil
}

func awsSsmFromUpstreamConfig(data *map[string]any, config *service.AwsSsmSshServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with SSH service type "aws_ssm" but AWS SSM SSH service configuration was not present`)
	}

	(*data)["ssm_target_type"] = config.SsmTargetType

	switch config.SsmTargetType {
	case service.SsmTargetTypeEc2:
		if config.AwsSsmEc2TargetConfiguration == nil {
			return diag.Errorf(`got a socket with AWS SSM target type "ec2" but AWS SSM EC2 target configuration was not present`)
		}
		(*data)["ec2_instance_id"] = config.AwsSsmEc2TargetConfiguration.Ec2InstanceId
		(*data)["ec2_instance_region"] = config.AwsSsmEc2TargetConfiguration.Ec2InstanceRegion
		(*data)["aws_credentials"] = shared.FromAwsCredentials(config.AwsSsmEc2TargetConfiguration.AwsCredentials)
	case service.SsmTargetTypeEcs:
		if config.AwsSsmEcsTargetConfiguration == nil {
			return diag.Errorf(`got a socket with AWS SSM target type "ecs" but AWS SSM ECS target configuration was not present`)
		}
		(*data)["ecs_cluster_region"] = config.AwsSsmEcsTargetConfiguration.EcsClusterRegion
		(*data)["ecs_cluster_name"] = config.AwsSsmEcsTargetConfiguration.EcsClusterName
		(*data)["ecs_service_name"] = config.AwsSsmEcsTargetConfiguration.EcsServiceName
		(*data)["aws_credentials"] = shared.FromAwsCredentials(config.AwsSsmEcsTargetConfiguration.AwsCredentials)
	default:
		return diag.Errorf(`sockets with AWS SSM target type "%s" not yet supported`, config.SsmTargetType)
	}

	return nil
}

func connectorBuiltInFromUpstreamConfig(data *map[string]any, config *service.BuiltInSshServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with SSH service type "connector_built_in" but built-in connector SSH service configuration was not present`)
	}

	(*data)["username_provider"] = config.UsernameProvider
	(*data)["username"] = config.Username

	return nil
}
