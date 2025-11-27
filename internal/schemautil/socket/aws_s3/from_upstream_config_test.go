package aws_s3

import (
	"testing"

	"github.com/borderzero/border0-go/types/common"
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromUpstreamConfig_WithAwsCredentials(t *testing.T) {
	accessKeyId := "AKIAIOSFODNN7EXAMPLE"
	secretAccessKey := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	config := &service.AwsS3ServiceConfiguration{
		AwsCredentials: &common.AwsCredentials{
			AwsAccessKeyId:     &accessKeyId,
			AwsSecretAccessKey: &secretAccessKey,
		},
	}

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"aws_s3_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_credentials": shared.AwsCredentialsSchema,
				},
			},
		},
	}, map[string]any{})

	diags := FromUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")

	awsS3ConfigRaw, ok := d.GetOk("aws_s3_configuration")
	require.True(t, ok, "Expected aws_s3_configuration to be set")

	awsS3ConfigList := awsS3ConfigRaw.([]any)
	require.Len(t, awsS3ConfigList, 1)

	awsS3Config := awsS3ConfigList[0].(map[string]any)
	awsCredsRaw, ok := awsS3Config["aws_credentials"]
	require.True(t, ok, "Expected aws_credentials to be set")

	awsCredsList := awsCredsRaw.([]any)
	require.Len(t, awsCredsList, 1)

	awsCreds := awsCredsList[0].(map[string]any)
	assert.Equal(t, "AKIAIOSFODNN7EXAMPLE", awsCreds["access_key_id"])
	assert.Equal(t, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", awsCreds["secret_access_key"])
}

func TestFromUpstreamConfig_WithAwsProfile(t *testing.T) {
	profile := "my-profile"

	config := &service.AwsS3ServiceConfiguration{
		AwsCredentials: &common.AwsCredentials{
			AwsProfile: &profile,
		},
	}

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"aws_s3_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_credentials": shared.AwsCredentialsSchema,
				},
			},
		},
	}, map[string]any{})

	diags := FromUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")

	awsS3ConfigRaw, ok := d.GetOk("aws_s3_configuration")
	require.True(t, ok, "Expected aws_s3_configuration to be set")

	awsS3ConfigList := awsS3ConfigRaw.([]any)
	require.Len(t, awsS3ConfigList, 1)

	awsS3Config := awsS3ConfigList[0].(map[string]any)
	awsCredsRaw, ok := awsS3Config["aws_credentials"]
	require.True(t, ok, "Expected aws_credentials to be set")

	awsCredsList := awsCredsRaw.([]any)
	require.Len(t, awsCredsList, 1)

	awsCreds := awsCredsList[0].(map[string]any)
	assert.Equal(t, "my-profile", awsCreds["profile"])
}

func TestFromUpstreamConfig_WithoutAwsCredentials(t *testing.T) {
	config := &service.AwsS3ServiceConfiguration{}

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"aws_s3_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_credentials": shared.AwsCredentialsSchema,
				},
			},
		},
	}, map[string]any{})

	diags := FromUpstreamConfig(d, config)

	require.False(t, diags.HasError(), "Expected no errors")

	awsS3ConfigRaw, ok := d.GetOk("aws_s3_configuration")
	require.True(t, ok, "Expected aws_s3_configuration to be set")

	awsS3ConfigList := awsS3ConfigRaw.([]any)
	require.Len(t, awsS3ConfigList, 1)

	// When aws_credentials is nil, the list item might be nil
	if awsS3ConfigList[0] == nil {
		// This is expected when there are no credentials
		return
	}

	awsS3Config := awsS3ConfigList[0].(map[string]any)
	_, hasAwsCreds := awsS3Config["aws_credentials"]
	assert.False(t, hasAwsCreds, "Expected aws_credentials not to be set when nil")
}

func TestFromUpstreamConfig_NilConfiguration(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"aws_s3_configuration": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_credentials": shared.AwsCredentialsSchema,
				},
			},
		},
	}, map[string]any{})

	diags := FromUpstreamConfig(d, nil)

	require.True(t, diags.HasError(), "Expected error for nil configuration")
	assert.Contains(t, diags[0].Summary, "AWS S3 service configuration was not present")
}
