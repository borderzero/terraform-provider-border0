package exit_node

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "exit_node_configuration"
// attribute into a socket's Exit Node service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.ExitNodeServiceConfiguration) diag.Diagnostics {
	// NOTE(@adrianosela): currently exit_node_configuration and the config object itself have
	// no attributes, so there is nothing to do here. If that ever changes, follow the pattern
	// in subnet_routes package in the parent directory.
	if config == nil {
		config = new(service.ExitNodeServiceConfiguration)
	}
	return nil
}
