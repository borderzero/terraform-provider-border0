package aws_s3

import (
	"testing"

	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToUpstreamConfig_WithAwsCredentials(t *testing.T) {
	awsCredsList := []any{
		map[string]any{
			"access_key_id":     "AKIAIOSFODNN7EXAMPLE",
			"secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		},
	}

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"aws_s3_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_credentials": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"access_key_id": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"secret_access_key": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}, map[string]any{
		"aws_s3_configuration": []any{
			map[string]any{
				"aws_credentials": awsCredsList,
			},
		},
	})

	config := &service.AwsS3ServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	require.NotNil(t, config.AwsCredentials)
	require.NotNil(t, config.AwsCredentials.AwsAccessKeyId)
	assert.Equal(t, "AKIAIOSFODNN7EXAMPLE", *config.AwsCredentials.AwsAccessKeyId)
	require.NotNil(t, config.AwsCredentials.AwsSecretAccessKey)
	assert.Equal(t, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", *config.AwsCredentials.AwsSecretAccessKey)
}

func TestToUpstreamConfig_WithoutAwsCredentials(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"aws_s3_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_credentials": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Resource{Schema: map[string]*schema.Schema{}},
					},
				},
			},
		},
	}, map[string]any{
		"aws_s3_configuration": []any{
			map[string]any{},
		},
	})

	config := &service.AwsS3ServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Nil(t, config.AwsCredentials)
}

func TestToUpstreamConfig_EmptyConfiguration(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"aws_s3_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Resource{Schema: map[string]*schema.Schema{}},
		},
	}, map[string]any{})

	config := &service.AwsS3ServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	assert.Nil(t, config.AwsCredentials)
}

func TestToUpstreamConfig_WithAwsProfile(t *testing.T) {
	awsCredsList := []any{
		map[string]any{
			"profile": "my-profile",
		},
	}

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"aws_s3_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_credentials": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"profile": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}, map[string]any{
		"aws_s3_configuration": []any{
			map[string]any{
				"aws_credentials": awsCredsList,
			},
		},
	})

	config := &service.AwsS3ServiceConfiguration{}
	diags := ToUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")
	require.NotNil(t, config.AwsCredentials)
	require.NotNil(t, config.AwsCredentials.AwsProfile)
	assert.Equal(t, "my-profile", *config.AwsCredentials.AwsProfile)
}
