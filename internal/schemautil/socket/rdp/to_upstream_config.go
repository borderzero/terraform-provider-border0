package rdp

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "rdp_configuration" attribute into
// a socket's RDP service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.RdpServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("rdp_configuration"); ok {
		if rdpConfigsList := v.([]any); len(rdpConfigsList) > 0 {
			data = rdpConfigsList[0].(map[string]any)
		}
	}

	if config == nil {
		config = new(service.RdpServiceConfiguration)
	}

	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}

	return nil
}
