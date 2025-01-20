package exit_node

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's exit node service configuration into terraform resource
// data for the "exit_node_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.ExitNodeServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "exit_node" but Exit Node service configuration was not present`)
	}
	data := make(map[string]any)
	if config.Mode != "" {
		data["mode"] = config.Mode
	}
	if err := d.Set("exit_node_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "subnet_router_configuration"`)
	}
	return nil
}
