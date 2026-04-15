package rdp

import (
	"testing"

	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToUpstreamConfig_AllFields(t *testing.T) {
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
	}, map[string]any{
		"rdp_configuration": []any{
			map[string]any{
				"hostname": "10.0.0.1",
				"port":     3389,
				"username": "admin",
				"password": "secret",
				"domain":   "corp.example.com",
			},
		},
	})

	config := &service.RdpServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, "10.0.0.1", config.Hostname)
	assert.Equal(t, uint16(3389), config.Port)
	assert.Equal(t, "admin", config.Username)
	assert.Equal(t, "secret", config.Password)
	assert.Equal(t, "corp.example.com", config.Domain)
}

func TestToUpstreamConfig_HostnameAndPortOnly(t *testing.T) {
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
	}, map[string]any{
		"rdp_configuration": []any{
			map[string]any{
				"hostname": "10.0.0.1",
				"port":     3389,
			},
		},
	})

	config := &service.RdpServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, "10.0.0.1", config.Hostname)
	assert.Equal(t, uint16(3389), config.Port)
	assert.Empty(t, config.Username)
	assert.Empty(t, config.Password)
	assert.Empty(t, config.Domain)
}

func TestToUpstreamConfig_UsernameAndPasswordWithoutDomain(t *testing.T) {
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
	}, map[string]any{
		"rdp_configuration": []any{
			map[string]any{
				"hostname": "10.0.0.1",
				"port":     3389,
				"username": "admin",
				"password": "secret",
			},
		},
	})

	config := &service.RdpServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Equal(t, "admin", config.Username)
	assert.Equal(t, "secret", config.Password)
	assert.Empty(t, config.Domain)
}

func TestToUpstreamConfig_EmptyConfig(t *testing.T) {
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

	config := &service.RdpServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Empty(t, config.Hostname)
	assert.Equal(t, uint16(0), config.Port)
	assert.Empty(t, config.Username)
	assert.Empty(t, config.Password)
	assert.Empty(t, config.Domain)
}

func TestToUpstreamConfig_NilConfig(t *testing.T) {
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
	}, map[string]any{
		"rdp_configuration": []any{
			map[string]any{
				"hostname": "10.0.0.1",
				"port":     3389,
				"username": "admin",
				"password": "secret",
			},
		},
	})

	// passing nil config should not panic
	diags := ToUpstreamConfig(d, nil)
	require.False(t, diags.HasError(), "Expected no errors")
}
