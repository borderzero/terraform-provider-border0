package snowflake

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's Snowflake service configuration into terraform resource data for
// the "snowflake_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.SnowflakeServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "snowflake" but Snowflake service configuration was not present`)
	}
	data := map[string]any{
		"account":  config.Account,
		"username": config.Username,
		"password": config.Password,
	}
	if err := d.Set("snowflake_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "snowflake_configuration"`)
	}
	return nil
}
