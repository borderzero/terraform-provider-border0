package aws_s3

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's AWS S3 service configuration into terraform resource data for
// the "aws_s3_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.AwsS3ServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "aws_s3" but AWS S3 service configuration was not present`)
	}

	data := make(map[string]any)

	if config.AwsCredentials != nil {
		data["aws_credentials"] = shared.FromAwsCredentials(config.AwsCredentials)
	}

	if err := d.Set("aws_s3_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "aws_s3_configuration"`)
	}

	return nil
}
