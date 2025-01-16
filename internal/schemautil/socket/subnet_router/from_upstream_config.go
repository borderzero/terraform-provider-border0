package subnet_router

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's subnet router service configuration into terraform resource
// data for the "subnet_router_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.SubnetRouterServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "subnet_router" but Subnet Router service configuration was not present`)
	}
	data := map[string]any{
		"ipv4_cidr_ranges": config.IPv4CIDRRanges,
		"ipv6_cidr_ranges": config.IPv6CIDRRanges,
	}
	if err := d.Set("subnet_router_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "subnet_router_configuration"`)
	}
	return nil
}
