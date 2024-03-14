package vpn

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's VPN service configuration into terraform resource data for
// the "vpn_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.VpnServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "vpn" but VPN service configuration was not present`)
	}
	data := map[string]any{
		"dhcp_pool_subnet":  config.DHCPPoolSubnet,
		"advertised_routes": config.AdvertisedRoutes,
	}
	if err := d.Set("vpn_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "vpn_configuration"`)
	}
	return nil
}
