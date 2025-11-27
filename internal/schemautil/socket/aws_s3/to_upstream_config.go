package aws_s3

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "aws_s3_configuration" attribute into
// a socket's AWS S3 service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.AwsS3ServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("aws_s3_configuration"); ok {
		if awsS3ConfigsList := v.([]any); len(awsS3ConfigsList) > 0 && awsS3ConfigsList[0] != nil {
			data = awsS3ConfigsList[0].(map[string]any)
		}
	}

	if v, ok := data["aws_credentials"]; ok {
		config.AwsCredentials = shared.ToAwsCredentials(v)
	}

	return nil
}
