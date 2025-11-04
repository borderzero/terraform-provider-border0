package ssh

import (
	"testing"

	"github.com/borderzero/border0-go/types/common"
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKubectlExecFromUpstreamConfig_Standard(t *testing.T) {
	upstream := &service.KubectlExecSshServiceConfiguration{
		KubectlExecTargetType: service.KubectlExecTargetTypeStandard,
		BaseKubectlExecTargetConfiguration: service.BaseKubectlExecTargetConfiguration{
			NamespaceAllowlist: []string{"default", "kube-system"},
			NamespaceSelectorsAllowlist: map[string]map[string][]string{
				"default": {
					"app": []string{"nginx", "redis"},
				},
			},
		},
		StandardKubectlExecTargetConfiguration: &service.StandardKubectlExecTargetConfiguration{
			KubeconfigPath: "/root/.kube/config",
			MasterUrl:      "https://k8s.example.com",
		},
	}

	data := make(map[string]any)
	diags := kubectlExecFromUpstreamConfig(&data, upstream)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.KubectlExecTargetTypeStandard, data["kubectl_exec_target_type"])
	assert.Equal(t, "/root/.kube/config", data["kubeconfig_path"])
	assert.Equal(t, "https://k8s.example.com", data["master_url"])
	assert.Equal(t, []string{"default", "kube-system"}, data["namespace_allowlist"])

	// Check namespace selectors were serialized to JSON
	selectorsJSON, ok := data["namespace_selectors_allowlist"].(string)
	require.True(t, ok, "namespace_selectors_allowlist should be a string")
	assert.Contains(t, selectorsJSON, "default")
	assert.Contains(t, selectorsJSON, "nginx")
}

func TestKubectlExecFromUpstreamConfig_StandardMinimal(t *testing.T) {
	upstream := &service.KubectlExecSshServiceConfiguration{
		KubectlExecTargetType: service.KubectlExecTargetTypeStandard,
		StandardKubectlExecTargetConfiguration: &service.StandardKubectlExecTargetConfiguration{
			KubeconfigPath: "/root/.kube/config",
		},
	}

	data := make(map[string]any)
	diags := kubectlExecFromUpstreamConfig(&data, upstream)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.KubectlExecTargetTypeStandard, data["kubectl_exec_target_type"])
	assert.Equal(t, "/root/.kube/config", data["kubeconfig_path"])
	assert.NotContains(t, data, "master_url")
	assert.NotContains(t, data, "namespace_allowlist")
	assert.NotContains(t, data, "namespace_selectors_allowlist")
}

func TestKubectlExecFromUpstreamConfig_AwsEks(t *testing.T) {
	accessKeyId := "AKIAIOSFODNN7EXAMPLE"
	secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	upstream := &service.KubectlExecSshServiceConfiguration{
		KubectlExecTargetType: service.KubectlExecTargetTypeAwsEks,
		AwsEksKubectlExecTargetConfiguration: &service.AwsEksKubectlExecTargetConfiguration{
			EksClusterName:   "my-cluster",
			EksClusterRegion: "us-west-2",
			AwsCredentials: &common.AwsCredentials{
				AwsAccessKeyId:     &accessKeyId,
				AwsSecretAccessKey: &secretAccessKey,
			},
		},
	}

	data := make(map[string]any)
	diags := kubectlExecFromUpstreamConfig(&data, upstream)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, service.KubectlExecTargetTypeAwsEks, data["kubectl_exec_target_type"])
	assert.Equal(t, "my-cluster", data["eks_cluster_name"])
	assert.Equal(t, "us-west-2", data["eks_cluster_region"])
	assert.NotNil(t, data["aws_credentials"])
}

func TestKubectlExecFromUpstreamConfig_EmptyTargetType(t *testing.T) {
	upstream := &service.KubectlExecSshServiceConfiguration{
		KubectlExecTargetType: "",
		StandardKubectlExecTargetConfiguration: &service.StandardKubectlExecTargetConfiguration{
			KubeconfigPath: "/root/.kube/config",
		},
	}

	data := make(map[string]any)
	diags := kubectlExecFromUpstreamConfig(&data, upstream)

	require.False(t, diags.HasError(), "Expected no errors for empty target type")
	assert.Equal(t, "", data["kubectl_exec_target_type"])
	assert.Equal(t, "/root/.kube/config", data["kubeconfig_path"])
}

func TestKubectlExecFromUpstreamConfig_NilConfig(t *testing.T) {
	data := make(map[string]any)
	diags := kubectlExecFromUpstreamConfig(&data, nil)

	require.True(t, diags.HasError(), "Expected error for nil config")
	assert.Contains(t, diags[0].Summary, "not present")
}

func TestKubectlExecFromUpstreamConfig_InvalidTargetType(t *testing.T) {
	upstream := &service.KubectlExecSshServiceConfiguration{
		KubectlExecTargetType: "invalid_type",
	}

	data := make(map[string]any)
	diags := kubectlExecFromUpstreamConfig(&data, upstream)

	require.True(t, diags.HasError(), "Expected error for invalid target type")
	assert.Contains(t, diags[0].Summary, "invalid")
}

func TestDockerExecFromUpstreamConfig_WithAllowlist(t *testing.T) {
	upstream := &service.DockerExecSshServiceConfiguration{
		ContainerNameAllowlist: []string{"nginx*", "postgres", "redis*"},
	}

	data := make(map[string]any)
	diags := dockerExecFromUpstreamConfig(&data, upstream)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, []string{"nginx*", "postgres", "redis*"}, data["container_name_allowlist"])
}

func TestDockerExecFromUpstreamConfig_EmptyAllowlist(t *testing.T) {
	upstream := &service.DockerExecSshServiceConfiguration{
		ContainerNameAllowlist: []string{},
	}

	data := make(map[string]any)
	diags := dockerExecFromUpstreamConfig(&data, upstream)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.NotContains(t, data, "container_name_allowlist")
}

func TestDockerExecFromUpstreamConfig_NilConfig(t *testing.T) {
	data := make(map[string]any)
	diags := dockerExecFromUpstreamConfig(&data, nil)

	require.True(t, diags.HasError(), "Expected error for nil config")
	assert.Contains(t, diags[0].Summary, "not present")
}

func TestFromUpstreamConfig_KubectlExec(t *testing.T) {
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
	}, map[string]any{})

	config := &service.SshServiceConfiguration{
		SshServiceType: service.SshServiceTypeKubectlExec,
		KubectlExecSshServiceConfiguration: &service.KubectlExecSshServiceConfiguration{
			KubectlExecTargetType: service.KubectlExecTargetTypeStandard,
			StandardKubectlExecTargetConfiguration: &service.StandardKubectlExecTargetConfiguration{
				KubeconfigPath: "/root/.kube/config",
			},
		},
	}

	diags := FromUpstreamConfig(d, config)
	require.False(t, diags.HasError(), "Expected no errors")

	sshConfig := d.Get("ssh_configuration").([]any)
	require.Len(t, sshConfig, 1)

	configMap := sshConfig[0].(map[string]any)
	assert.Equal(t, service.SshServiceTypeKubectlExec, configMap["service_type"])
	assert.Equal(t, service.KubectlExecTargetTypeStandard, configMap["kubectl_exec_target_type"])
	assert.Equal(t, "/root/.kube/config", configMap["kubeconfig_path"])
}

func TestFromUpstreamConfig_DockerExec(t *testing.T) {
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
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	}, map[string]any{})

	config := &service.SshServiceConfiguration{
		SshServiceType: service.SshServiceTypeDockerExec,
		DockerExecSshServiceConfiguration: &service.DockerExecSshServiceConfiguration{
			ContainerNameAllowlist: []string{"nginx*", "postgres"},
		},
	}

	diags := FromUpstreamConfig(d, config)
	require.False(t, diags.HasError(), "Expected no errors")

	sshConfig := d.Get("ssh_configuration").([]any)
	require.Len(t, sshConfig, 1)

	configMap := sshConfig[0].(map[string]any)
	assert.Equal(t, service.SshServiceTypeDockerExec, configMap["service_type"])

	// Terraform returns []any not []string
	containerList := configMap["container_name_allowlist"].([]any)
	assert.Len(t, containerList, 2)
	assert.Contains(t, containerList, "nginx*")
	assert.Contains(t, containerList, "postgres")
}
