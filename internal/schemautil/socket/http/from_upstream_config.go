package http

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's HTTP service configuration into terraform resource data for
// the "http_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.HttpServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "http" but HTTP service configuration was not present`)
	}

	data := map[string]any{
		"service_type": config.HttpServiceType,
	}

	var diags diag.Diagnostics

	switch config.HttpServiceType {
	case service.HttpServiceTypeStandard:
		diags = standardFromUpstreamConfig(&data, config.StandardHttpServiceConfiguration)
	case service.HttpServiceTypeConnectorFileServer:
		diags = connectorFileServerFromUpstreamConfig(&data, config.FileServerHttpServiceConfiguration)
	default:
		return diag.Errorf(`sockets with HTTP service type "%s" not yet supported`, config.HttpServiceType)
	}

	if diags.HasError() {
		return diags
	}

	if err := d.Set("http_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "http_configuration"`)
	}

	return nil
}

func standardFromUpstreamConfig(data *map[string]any, config *service.StandardHttpServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with HTTP service type "standard" but standard HTTP service configuration was not present`)
	}

	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port
	(*data)["host_header"] = config.HostHeader

	return nil
}

func connectorFileServerFromUpstreamConfig(data *map[string]any, config *service.FileServerHttpServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with http service type "connector_file_server" but file server HTTP service configuration was not present`)
	}

	(*data)["top_level_directory"] = config.TopLevelDirectory

	return nil
}
