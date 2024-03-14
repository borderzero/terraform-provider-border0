package vnc

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's VNC service configuration into terraform resource data for
// the "vnc_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.VncServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "vnc" but VNC service configuration was not present`)
	}
	data := map[string]any{
		"hostname": config.Hostname,
		"port":     config.Port,
	}
	if err := d.Set("vnc_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "vnc_configuration"`)
	}
	return nil
}
