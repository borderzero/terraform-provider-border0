package vpn

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/schemaconvert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "vpn_configuration" attribute into
// a socket's VPN service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.VpnServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("vpn_configuration"); ok {
		if vncConfigsList := v.([]any); len(vncConfigsList) > 0 {
			data = vncConfigsList[0].(map[string]any)
		}
	}

	if config == nil {
		config = new(service.VpnServiceConfiguration)
	}

	if v, ok := data["dhcp_pool_subnet"]; ok {
		config.DHCPPoolSubnet = v.(string)
	}
	if v, ok := data["advertised_routes"]; ok {
		config.AdvertisedRoutes = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}

	return nil
}
