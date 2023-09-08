package shared

import (
	"github.com/borderzero/border0-go/lib/types/pointer"
	"github.com/borderzero/border0-go/types/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AwsCredentialsSchema is the schema for the `aws_credentials` block. This block is used to configure
// the upstream service's AWS credentials. Configure `access_key_id`, `secret_access_key` and `session_token`
// to use static credentials, or configure `profile` to use credentials from the shared credentials file.
var AwsCredentialsSchema = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The upstream AWS access key id.",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The upstream AWS secret access key.",
			},
			"session_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The upstream AWS session token.",
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The upstream AWS profile.",
			},
		},
	},
	Description: "The upstream service's AWS credentials.",
}

// FromAwsCredentials converts the `common.AwsCredentials` struct to a `[]map[string]any` for use in the
// `aws_credentials` block.
func FromAwsCredentials(creds *common.AwsCredentials) []map[string]any {
	if creds == nil {
		return nil
	}
	return []map[string]any{
		{
			"access_key_id":     creds.AwsAccessKeyId,
			"secret_access_key": creds.AwsSecretAccessKey,
			"session_token":     creds.AwsSessionToken,
			"profile":           creds.AwsProfile,
		},
	}
}

// ToAwsCredentials converts the `aws_credentials` block to a `common.AwsCredentials` struct.
func ToAwsCredentials(v any, creds *common.AwsCredentials) {
	if awsCredentialsList := v.([]any); len(awsCredentialsList) > 0 {
		awsCredentials := awsCredentialsList[0].(map[string]any)
		if creds == nil {
			creds = new(common.AwsCredentials)
		}
		if v, ok := awsCredentials["access_key_id"]; ok {
			creds.AwsAccessKeyId = pointer.To(v.(string))
		}
		if v, ok := awsCredentials["secret_access_key"]; ok {
			creds.AwsSecretAccessKey = pointer.To(v.(string))
		}
		if v, ok := awsCredentials["session_token"]; ok {
			creds.AwsSessionToken = pointer.To(v.(string))
		}
		if v, ok := awsCredentials["profile"]; ok {
			creds.AwsProfile = pointer.To(v.(string))
		}
	}
}
