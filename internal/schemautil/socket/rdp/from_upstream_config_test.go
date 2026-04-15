package rdp

import (
	"testing"

	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromUpstreamConfig_AllFields(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"rdp_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hostname": {Type: schema.TypeString, Optional: true},
					"port":     {Type: schema.TypeInt, Optional: true},
					"username": {Type: schema.TypeString, Optional: true},
					"password": {Type: schema.TypeString, Optional: true},
					"domain":   {Type: schema.TypeString, Optional: true},
				},
			},
		},
	}, map[string]any{})

	config := &service.RdpServiceConfiguration{
		HostnameAndPort: service.HostnameAndPort{
			Hostname: "10.0.0.1",
			Port:     3389,
		},
		UsernameAndPassword: service.UsernameAndPassword{
			Username: "admin",
			Password: "secret",
		},
		Domain: "corp.example.com",
	}

	diags := FromUpstreamConfig(d, config)
	require.False(t, diags.HasError(), "Expected no errors")

	rdpConfig := d.Get("rdp_configuration").([]any)
	require.Len(t, rdpConfig, 1)

	configMap := rdpConfig[0].(map[string]any)
	assert.Equal(t, "10.0.0.1", configMap["hostname"])
	assert.Equal(t, 3389, configMap["port"])
	assert.Equal(t, "admin", configMap["username"])
	assert.Equal(t, "secret", configMap["password"])
	assert.Equal(t, "corp.example.com", configMap["domain"])
}

func TestFromUpstreamConfig_HostnameAndPortOnly(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"rdp_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hostname": {Type: schema.TypeString, Optional: true},
					"port":     {Type: schema.TypeInt, Optional: true},
					"username": {Type: schema.TypeString, Optional: true},
					"password": {Type: schema.TypeString, Optional: true},
					"domain":   {Type: schema.TypeString, Optional: true},
				},
			},
		},
	}, map[string]any{})

	config := &service.RdpServiceConfiguration{
		HostnameAndPort: service.HostnameAndPort{
			Hostname: "10.0.0.1",
			Port:     3389,
		},
	}

	diags := FromUpstreamConfig(d, config)
	require.False(t, diags.HasError(), "Expected no errors")

	rdpConfig := d.Get("rdp_configuration").([]any)
	require.Len(t, rdpConfig, 1)

	configMap := rdpConfig[0].(map[string]any)
	assert.Equal(t, "10.0.0.1", configMap["hostname"])
	assert.Equal(t, 3389, configMap["port"])
	assert.Equal(t, "", configMap["username"])
	assert.Equal(t, "", configMap["password"])
	assert.Equal(t, "", configMap["domain"])
}

func TestFromUpstreamConfig_NilConfig(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"rdp_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hostname": {Type: schema.TypeString, Optional: true},
					"port":     {Type: schema.TypeInt, Optional: true},
					"username": {Type: schema.TypeString, Optional: true},
					"password": {Type: schema.TypeString, Optional: true},
					"domain":   {Type: schema.TypeString, Optional: true},
				},
			},
		},
	}, map[string]any{})

	diags := FromUpstreamConfig(d, nil)
	require.True(t, diags.HasError(), "Expected error for nil config")
	assert.Contains(t, diags[0].Summary, "not present")
}
