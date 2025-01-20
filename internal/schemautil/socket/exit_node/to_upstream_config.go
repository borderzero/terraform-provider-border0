package exit_node

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "exit_node_configuration"
// attribute into a socket's Exit Node service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.ExitNodeServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)
	if v, ok := d.GetOk("exit_node_configuration"); ok {
		if srConfigsList := v.([]any); len(srConfigsList) > 0 {
			data = srConfigsList[0].(map[string]any)
		}
	}
	if config == nil {
		config = new(service.ExitNodeServiceConfiguration)
	}
	if v, ok := data["mode"]; ok {
		if modeStr, ok := v.(string); ok && modeStr != "" {
			config.Mode = modeStr
		}
	}
	return nil

}
