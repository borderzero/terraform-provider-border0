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
	// NOTE(@adrianosela): currently exit_node_configuration and the config object itself have
	// no attributes, so there is nothing to do here. If that ever changes, follow the pattern
	// in subnet_routes package in the parent directory.
	if err := d.Set("exit_node_configuration", []map[string]any{}); err != nil {
		return diagnostics.Error(err, `Failed to set "exit_node_configuration"`)
	}
	return nil
}
