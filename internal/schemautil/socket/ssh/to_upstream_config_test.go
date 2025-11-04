package ssh

import (
	"encoding/json"
	"testing"

	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKubectlExecToUpstreamConfig_Standard(t *testing.T) {
	namespaceSet := schema.NewSet(schema.HashString, []any{"default", "kube-system"})

	selectorsMap := map[string]map[string][]string{
		"default": {
			"app": []string{"nginx", "redis"},
		},
	}
	selectorsJSON, _ := json.Marshal(selectorsMap)

	data := map[string]any{
		"kubectl_exec_target_type":      service.KubectlExecTargetTypeStandard,
		"kubeconfig_path":               "/root/.kube/config",
		"master_url":                    "https://k8s.example.com",
		"namespace_allowlist":           namespaceSet,
		"namespace_selectors_allowlist": string(selectorsJSON),
	}

	config := &service.KubectlExecSshServiceConfiguration{}
	diags := kubectlExecToUpstreamConfig(data, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.KubectlExecTargetTypeStandard, config.KubectlExecTargetType)
	assert.Contains(t, config.NamespaceAllowlist, "default")
	assert.Contains(t, config.NamespaceAllowlist, "kube-system")
	assert.NotNil(t, config.StandardKubectlExecTargetConfiguration)
	assert.Equal(t, "/root/.kube/config", config.StandardKubectlExecTargetConfiguration.KubeconfigPath)
	assert.Equal(t, "https://k8s.example.com", config.StandardKubectlExecTargetConfiguration.MasterUrl)
	assert.NotNil(t, config.NamespaceSelectorsAllowlist)
	assert.Equal(t, []string{"nginx", "redis"}, config.NamespaceSelectorsAllowlist["default"]["app"])
}

func TestKubectlExecToUpstreamConfig_StandardMinimal(t *testing.T) {
	data := map[string]any{
		"kubeconfig_path": "/root/.kube/config",
	}

	config := &service.KubectlExecSshServiceConfiguration{}
	diags := kubectlExecToUpstreamConfig(data, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.KubectlExecTargetTypeStandard, config.KubectlExecTargetType)
	assert.NotNil(t, config.StandardKubectlExecTargetConfiguration)
	assert.Equal(t, "/root/.kube/config", config.StandardKubectlExecTargetConfiguration.KubeconfigPath)
	assert.Empty(t, config.StandardKubectlExecTargetConfiguration.MasterUrl)
	assert.Nil(t, config.NamespaceAllowlist)
	assert.Nil(t, config.NamespaceSelectorsAllowlist)
}

func TestKubectlExecToUpstreamConfig_EmptyTargetType(t *testing.T) {
	data := map[string]any{
		"kubectl_exec_target_type": "",
		"kubeconfig_path":          "/root/.kube/config",
	}

	config := &service.KubectlExecSshServiceConfiguration{}
	diags := kubectlExecToUpstreamConfig(data, config)

	require.False(t, diags.HasError(), "Expected no errors")
	// Empty string should default to standard
	assert.Equal(t, service.KubectlExecTargetTypeStandard, config.KubectlExecTargetType)
	assert.NotNil(t, config.StandardKubectlExecTargetConfiguration)
}

func TestKubectlExecToUpstreamConfig_AwsEks(t *testing.T) {
	awsCredsList := []any{
		map[string]any{
			"access_key_id":     "AKIAIOSFODNN7EXAMPLE",
			"secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		},
	}

	data := map[string]any{
		"kubectl_exec_target_type": service.KubectlExecTargetTypeAwsEks,
		"eks_cluster_name":         "my-cluster",
		"eks_cluster_region":       "us-west-2",
		"aws_credentials":          awsCredsList,
	}

	config := &service.KubectlExecSshServiceConfiguration{}
	diags := kubectlExecToUpstreamConfig(data, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.KubectlExecTargetTypeAwsEks, config.KubectlExecTargetType)
	assert.NotNil(t, config.AwsEksKubectlExecTargetConfiguration)
	assert.Equal(t, "my-cluster", config.AwsEksKubectlExecTargetConfiguration.EksClusterName)
	assert.Equal(t, "us-west-2", config.AwsEksKubectlExecTargetConfiguration.EksClusterRegion)
	assert.NotNil(t, config.AwsEksKubectlExecTargetConfiguration.AwsCredentials)
	assert.NotNil(t, config.AwsEksKubectlExecTargetConfiguration.AwsCredentials.AwsAccessKeyId)
	assert.Equal(t, "AKIAIOSFODNN7EXAMPLE", *config.AwsEksKubectlExecTargetConfiguration.AwsCredentials.AwsAccessKeyId)
}

func TestKubectlExecToUpstreamConfig_InvalidJSON(t *testing.T) {
	data := map[string]any{
		"namespace_selectors_allowlist": "invalid json {{{",
	}

	config := &service.KubectlExecSshServiceConfiguration{}
	diags := kubectlExecToUpstreamConfig(data, config)

	require.True(t, diags.HasError(), "Expected error for invalid JSON")
	assert.Contains(t, diags[0].Summary, "unmarshal")
}

func TestKubectlExecToUpstreamConfig_InvalidTargetType(t *testing.T) {
	data := map[string]any{
		"kubectl_exec_target_type": "invalid_type",
	}

	config := &service.KubectlExecSshServiceConfiguration{}
	diags := kubectlExecToUpstreamConfig(data, config)

	require.True(t, diags.HasError(), "Expected error for invalid target type")
	assert.Contains(t, diags[0].Summary, "invalid")
}

func TestDockerExecToUpstreamConfig_WithAllowlist(t *testing.T) {
	containerSet := schema.NewSet(schema.HashString, []any{"nginx*", "postgres", "redis*"})

	data := map[string]any{
		"container_name_allowlist": containerSet,
	}

	config := &service.DockerExecSshServiceConfiguration{}
	diags := dockerExecToUpstreamConfig(data, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Len(t, config.ContainerNameAllowlist, 3)
	assert.Contains(t, config.ContainerNameAllowlist, "nginx*")
	assert.Contains(t, config.ContainerNameAllowlist, "postgres")
	assert.Contains(t, config.ContainerNameAllowlist, "redis*")
}

func TestDockerExecToUpstreamConfig_EmptyAllowlist(t *testing.T) {
	data := map[string]any{}

	config := &service.DockerExecSshServiceConfiguration{}
	diags := dockerExecToUpstreamConfig(data, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Nil(t, config.ContainerNameAllowlist)
}

func TestDockerExecToUpstreamConfig_EmptySet(t *testing.T) {
	containerSet := schema.NewSet(schema.HashString, []any{})

	data := map[string]any{
		"container_name_allowlist": containerSet,
	}

	config := &service.DockerExecSshServiceConfiguration{}
	diags := dockerExecToUpstreamConfig(data, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Empty(t, config.ContainerNameAllowlist)
}

func TestToUpstreamConfig_KubectlExec(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"ssh_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_type": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"kubectl_exec_target_type": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"kubeconfig_path": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}, map[string]any{
		"ssh_configuration": []any{
			map[string]any{
				"service_type":             service.SshServiceTypeKubectlExec,
				"kubectl_exec_target_type": service.KubectlExecTargetTypeStandard,
				"kubeconfig_path":          "/root/.kube/config",
			},
		},
	})

	config := &service.SshServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.SshServiceTypeKubectlExec, config.SshServiceType)
	assert.NotNil(t, config.KubectlExecSshServiceConfiguration)
	assert.Equal(t, service.KubectlExecTargetTypeStandard, config.KubectlExecSshServiceConfiguration.KubectlExecTargetType)
	assert.NotNil(t, config.KubectlExecSshServiceConfiguration.StandardKubectlExecTargetConfiguration)
	assert.Equal(t, "/root/.kube/config", config.KubectlExecSshServiceConfiguration.StandardKubectlExecTargetConfiguration.KubeconfigPath)
}

func TestToUpstreamConfig_DockerExec(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"ssh_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_type": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"container_name_allowlist": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	}, map[string]any{
		"ssh_configuration": []any{
			map[string]any{
				"service_type":             service.SshServiceTypeDockerExec,
				"container_name_allowlist": []any{"nginx*", "postgres"},
			},
		},
	})

	config := &service.SshServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.SshServiceTypeDockerExec, config.SshServiceType)
	assert.NotNil(t, config.DockerExecSshServiceConfiguration)
	assert.Len(t, config.DockerExecSshServiceConfiguration.ContainerNameAllowlist, 2)
	assert.Contains(t, config.DockerExecSshServiceConfiguration.ContainerNameAllowlist, "nginx*")
	assert.Contains(t, config.DockerExecSshServiceConfiguration.ContainerNameAllowlist, "postgres")
}

func TestToUpstreamConfig_KubectlExec_DefaultTargetType(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"ssh_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"service_type": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"kubeconfig_path": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}, map[string]any{
		"ssh_configuration": []any{
			map[string]any{
				"service_type":    service.SshServiceTypeKubectlExec,
				"kubeconfig_path": "/root/.kube/config",
			},
		},
	})

	config := &service.SshServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.SshServiceTypeKubectlExec, config.SshServiceType)
	assert.NotNil(t, config.KubectlExecSshServiceConfiguration)
	// Should default to standard when not specified
	assert.Equal(t, service.KubectlExecTargetTypeStandard, config.KubectlExecSshServiceConfiguration.KubectlExecTargetType)
}
