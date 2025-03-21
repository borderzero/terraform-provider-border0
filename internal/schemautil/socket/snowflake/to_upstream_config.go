package snowflake

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "snowflake_configuration" attribute into
// a socket's Snowflake service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.SnowflakeServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("snowflake_configuration"); ok {
		snowflakeConfigsList := v.([]any)
		if len(snowflakeConfigsList) < 1 {
			return diag.Errorf(`no snowflake_configuration found`)
		}
		if len(snowflakeConfigsList) > 1 {
			return diag.Errorf(`multiple instances of snowflake_configuration found, must have exactly 1`)
		}
		data = snowflakeConfigsList[0].(map[string]any)
	}

	if config == nil {
		config = new(service.SnowflakeServiceConfiguration)
	}
	if v, ok := data["account"]; ok {
		config.Account = v.(string)
	}
	if v, ok := data["username"]; ok {
		config.Username = v.(string)
	}
	if v, ok := data["password"]; ok {
		config.Password = v.(string)
	}

	return nil
}
