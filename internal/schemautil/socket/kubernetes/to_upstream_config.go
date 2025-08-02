package kubernetes

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "kubernetes_configuration" attribute into
// a socket's kubernetes service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.KubernetesServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("kubernetes_configuration"); ok {
		if configsList := v.([]any); len(configsList) > 0 {
			data = configsList[0].(map[string]any)
		}
	}

	serviceType := service.KubernetesServiceTypeStandard // default to "standard"
	if v, ok := data["service_type"]; ok {
		serviceType = v.(string)
	}
	config.KubernetesServiceType = serviceType

	if v, ok := data["impersonation_enabled"]; ok {
		config.ImpersonationEnabled = v.(bool)
	}

	switch serviceType {
	case service.KubernetesServiceTypeStandard:
		if config.StandardKubernetesServiceConfiguration == nil {
			config.StandardKubernetesServiceConfiguration = new(service.StandardKubernetesServiceConfiguration)
		}
		return standardToUpstreamConfig(data, config.StandardKubernetesServiceConfiguration)

	case service.KubernetesServiceTypeAwsEks:
		if config.AwsEksKubernetesServiceConfiguration == nil {
			config.AwsEksKubernetesServiceConfiguration = new(service.AwsEksKubernetesServiceConfiguration)
		}
		return awsEksToUpstreamConfig(data, config.AwsEksKubernetesServiceConfiguration)

	default:
		return diag.Errorf(`sockets with kubernetes service type "%s" not yet supported`, serviceType)
	}
}

func standardToUpstreamConfig(data map[string]any, config *service.StandardKubernetesServiceConfiguration) diag.Diagnostics {
	if v, ok := data["kubeconfig_path"]; ok {
		config.KubeconfigPath = v.(string)
	}
	if v, ok := data["context"]; ok {
		config.Context = v.(string)
	}
	if v, ok := data["server"]; ok {
		config.Server = v.(string)
	}
	if v, ok := data["certificate_authority"]; ok {
		config.CertificateAuthority = v.(string)
	}
	if v, ok := data["certificate_authority_data"]; ok {
		config.CertificateAuthorityData = v.(string)
	}
	if v, ok := data["client_certificate"]; ok {
		config.ClientCertificate = v.(string)
	}
	if v, ok := data["client_certificate_data"]; ok {
		config.ClientCertificateData = v.(string)
	}
	if v, ok := data["client_key"]; ok {
		config.ClientKey = v.(string)
	}
	if v, ok := data["client_key_data"]; ok {
		config.ClientKeyData = v.(string)
	}
	if v, ok := data["token"]; ok {
		config.Token = v.(string)
	}
	if v, ok := data["token_file"]; ok {
		config.TokenFile = v.(string)
	}

	return nil
}

func awsEksToUpstreamConfig(data map[string]any, config *service.AwsEksKubernetesServiceConfiguration) diag.Diagnostics {
	if v, ok := data["eks_cluster_name"]; ok {
		config.EksClusterName = v.(string)
	}
	if v, ok := data["eks_cluster_region"]; ok {
		config.EksClusterRegion = v.(string)
	}
	if v, ok := data["aws_credentials"]; ok {
		config.AwsCredentials = shared.ToAwsCredentials(v)
	}
	return nil
}
