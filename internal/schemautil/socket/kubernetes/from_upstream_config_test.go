package kubernetes

import (
	"testing"

	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FromUpstreamConfig(t *testing.T) {
	tests := []struct {
		name string

		// raw is what the practitioner wrote, config is what the api returned
		raw    map[string]any
		config *service.KubernetesServiceConfiguration

		// wantAttrs is checked one key at a time, so a case only lists what it
		// cares about.
		wantAttrs map[string]any

		// wantNoBlock means the read path must write no block at all
		wantNoBlock bool
		wantErr     string
	}{
		{
			// this is the bug. defaults from the api must not reach state
			name: "no block and no upstream config, nothing is read back",
			raw:  map[string]any{},
			config: &service.KubernetesServiceConfiguration{
				KubernetesServiceType:                  service.KubernetesServiceTypeStandard,
				StandardKubernetesServiceConfiguration: &service.StandardKubernetesServiceConfiguration{},
			},
			wantNoBlock: true,
		},
		{
			// this is the import case, and how real drift gets noticed
			name: "no block but upstream has config, block is read back",
			raw:  map[string]any{},
			config: &service.KubernetesServiceConfiguration{
				KubernetesServiceType: service.KubernetesServiceTypeStandard,
				StandardKubernetesServiceConfiguration: &service.StandardKubernetesServiceConfiguration{
					Server: "https://10.0.0.1:6443",
				},
			},
			wantAttrs: map[string]any{"server": "https://10.0.0.1:6443"},
		},
		{
			name: "no block but impersonation enabled, block is read back",
			raw:  map[string]any{},
			config: &service.KubernetesServiceConfiguration{
				KubernetesServiceType:                  service.KubernetesServiceTypeStandard,
				ImpersonationEnabled:                   true,
				StandardKubernetesServiceConfiguration: &service.StandardKubernetesServiceConfiguration{},
			},
			wantAttrs: map[string]any{"impersonation_enabled": true},
		},
		{
			// an empty block says nothing, but it was written on purpose, so it
			// stays. otherwise it would come back as a diff on every plan.
			name: "empty block written, block is read back",
			raw: map[string]any{
				"kubernetes_configuration": []any{
					map[string]any{"service_type": service.KubernetesServiceTypeStandard},
				},
			},
			config: &service.KubernetesServiceConfiguration{
				KubernetesServiceType:                  service.KubernetesServiceTypeStandard,
				StandardKubernetesServiceConfiguration: &service.StandardKubernetesServiceConfiguration{},
			},
			wantAttrs: map[string]any{
				"service_type":          service.KubernetesServiceTypeStandard,
				"impersonation_enabled": false,
			},
		},
		{
			name: "standard service type, all fields",
			raw:  map[string]any{},
			config: &service.KubernetesServiceConfiguration{
				KubernetesServiceType: service.KubernetesServiceTypeStandard,
				ImpersonationEnabled:  true,
				StandardKubernetesServiceConfiguration: &service.StandardKubernetesServiceConfiguration{
					KubeconfigPath:           "/root/.kube/config",
					Context:                  "test-context",
					Server:                   "https://10.0.0.1:6443",
					CertificateAuthority:     "/ca.crt",
					CertificateAuthorityData: "Y2EtZGF0YQ==",
					ClientCertificate:        "/client.crt",
					ClientCertificateData:    "Y2xpZW50LWRhdGE=",
					ClientKey:                "/client.key",
					ClientKeyData:            "a2V5LWRhdGE=",
					Token:                    "test-token",
					TokenFile:                "/token",
				},
			},
			wantAttrs: map[string]any{
				"service_type":               service.KubernetesServiceTypeStandard,
				"impersonation_enabled":      true,
				"kubeconfig_path":            "/root/.kube/config",
				"context":                    "test-context",
				"server":                     "https://10.0.0.1:6443",
				"certificate_authority":      "/ca.crt",
				"certificate_authority_data": "Y2EtZGF0YQ==",
				"client_certificate":         "/client.crt",
				"client_certificate_data":    "Y2xpZW50LWRhdGE=",
				"client_key":                 "/client.key",
				"client_key_data":            "a2V5LWRhdGE=",
				"token":                      "test-token",
				"token_file":                 "/token",
			},
		},
		{
			name: "aws eks service type",
			raw:  map[string]any{},
			config: &service.KubernetesServiceConfiguration{
				KubernetesServiceType: service.KubernetesServiceTypeAwsEks,
				AwsEksKubernetesServiceConfiguration: &service.AwsEksKubernetesServiceConfiguration{
					EksClusterName:   "test-cluster",
					EksClusterRegion: "us-east-2",
				},
			},
			wantAttrs: map[string]any{
				"service_type":       service.KubernetesServiceTypeAwsEks,
				"eks_cluster_name":   "test-cluster",
				"eks_cluster_region": "us-east-2",
			},
		},
		{
			name:    "nil config",
			raw:     map[string]any{},
			config:  nil,
			wantErr: "not present",
		},
		{
			name:    "unsupported service type",
			raw:     map[string]any{},
			config:  &service.KubernetesServiceConfiguration{KubernetesServiceType: "bogus"},
			wantErr: "not yet supported",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := testResourceData(t, test.raw)

			diags := FromUpstreamConfig(d, test.config)

			if test.wantErr != "" {
				require.True(t, diags.HasError())
				assert.Contains(t, diags[0].Summary, test.wantErr)
				return
			}
			require.False(t, diags.HasError())

			configs := d.Get("kubernetes_configuration").([]any)
			if test.wantNoBlock {
				assert.Empty(t, configs)
				return
			}
			require.Len(t, configs, 1)
			for key, want := range test.wantAttrs {
				assert.Equal(t, want, configs[0].(map[string]any)[key], key)
			}
		})
	}
}

// testResourceData builds resource data with the same "kubernetes_configuration"
// schema as the border0_socket resource. raw is what the practitioner wrote.
func testResourceData(t *testing.T, raw map[string]any) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"kubernetes_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_type":               {Type: schema.TypeString, Optional: true},
					"kubeconfig_path":            {Type: schema.TypeString, Optional: true},
					"context":                    {Type: schema.TypeString, Optional: true},
					"server":                     {Type: schema.TypeString, Optional: true},
					"certificate_authority":      {Type: schema.TypeString, Optional: true},
					"certificate_authority_data": {Type: schema.TypeString, Optional: true},
					"client_certificate":         {Type: schema.TypeString, Optional: true},
					"client_certificate_data":    {Type: schema.TypeString, Optional: true},
					"client_key":                 {Type: schema.TypeString, Optional: true},
					"client_key_data":            {Type: schema.TypeString, Optional: true},
					"token":                      {Type: schema.TypeString, Optional: true},
					"token_file":                 {Type: schema.TypeString, Optional: true},
					"impersonation_enabled":      {Type: schema.TypeBool, Optional: true},
					"eks_cluster_name":           {Type: schema.TypeString, Optional: true},
					"eks_cluster_region":         {Type: schema.TypeString, Optional: true},
					"aws_credentials":            shared.AwsCredentialsSchema,
				},
			},
		},
	}, raw)
}
