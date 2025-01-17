package subnet_router

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/schemaconvert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "subnet_router_configuration"
// attribute into a socket's Subnet Router service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.SubnetRouterServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("subnet_router_configuration"); ok {
		if srConfigsList := v.([]any); len(srConfigsList) > 0 {
			data = srConfigsList[0].(map[string]any)
		}
	}

	if config == nil {
		config = new(service.SubnetRouterServiceConfiguration)
	}

	if v, ok := data["ipv4_cidr_ranges"]; ok {
		config.IPv4CIDRRanges = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}
	if v, ok := data["ipv6_cidr_ranges"]; ok {
		config.IPv6CIDRRanges = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}

	return nil
}
