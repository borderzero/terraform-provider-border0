package kubernetes

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's kubernetes service configuration into terraform resource data for
// the "kubernetes_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.KubernetesServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "kubernetes" but kubernetes service configuration was not present`)
	}

	data := map[string]any{
		"service_type":          config.KubernetesServiceType,
		"impersonation_enabled": config.ImpersonationEnabled,
	}

	var diags diag.Diagnostics

	switch config.KubernetesServiceType {
	case service.KubernetesServiceTypeStandard:
		diags = standardFromUpstreamConfig(&data, config.StandardKubernetesServiceConfiguration)
	case service.KubernetesServiceTypeAwsEks:
		diags = awsEksFromUpstreamConfig(&data, config.AwsEksKubernetesServiceConfiguration)
	default:
		return diag.Errorf(`sockets with kubernetes service type "%s" not yet supported`, config.KubernetesServiceType)
	}

	if diags.HasError() {
		return diags
	}

	if err := d.Set("kubernetes_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "kubernetes_configuration"`)
	}

	return nil
}

func standardFromUpstreamConfig(data *map[string]any, config *service.StandardKubernetesServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with kubernetes service type "standard" but standard kubernetes service configuration was not present`)
	}

	(*data)["kubeconfig_path"] = config.KubeconfigPath
	(*data)["context"] = config.Context
	(*data)["certificate_authority"] = config.CertificateAuthority
	(*data)["certificate_authority_data"] = config.CertificateAuthorityData
	(*data)["client_certificate"] = config.ClientCertificate
	(*data)["client_certificate_data"] = config.ClientCertificateData
	(*data)["client_key"] = config.ClientKey
	(*data)["client_key_data"] = config.ClientKeyData
	(*data)["token"] = config.Token
	(*data)["token_file"] = config.TokenFile

	return nil
}

func awsEksFromUpstreamConfig(data *map[string]any, config *service.AwsEksKubernetesServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with kubernetes service type "aws_eks" but AWS EKS kubernetes service configuration was not present`)
	}

	(*data)["eks_cluster_name"] = config.EksClusterName
	(*data)["eks_cluster_region"] = config.EksClusterRegion
	(*data)["aws_credentials"] = shared.FromAwsCredentials(config.AwsCredentials)

	return nil
}
