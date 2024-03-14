package vnc

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "vnc_configuration" attribute into
// a socket's VNC service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.VncServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("vnc_configuration"); ok {
		if vncConfigsList := v.([]any); len(vncConfigsList) > 0 {
			data = vncConfigsList[0].(map[string]any)
		}
	}

	if config == nil {
		config = new(service.VncServiceConfiguration)
	}

	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}

	return nil
}
