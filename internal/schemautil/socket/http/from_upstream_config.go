package http

import (
	"fmt"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's HTTP service configuration into terraform resource data for
// the "http_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket, config *service.HttpServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "http" but HTTP service configuration was not present`)
	}

	data := map[string]any{
		"service_type": config.HttpServiceType,
	}

	var diags diag.Diagnostics

	switch config.HttpServiceType {
	case service.HttpServiceTypeStandard:
		diags = standardFromUpstreamConfig(&data, socket, config.StandardHttpServiceConfiguration)
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

func standardFromUpstreamConfig(data *map[string]any, socket *border0client.Socket, config *service.StandardHttpServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with HTTP service type "standard" but standard HTTP service configuration was not present`)
	}

	port := fmt.Sprintf(":%d", config.Port)
	if (socket.UpstreamType == "https" && config.Port == 443) ||
		(socket.UpstreamType == "http" && config.Port == 80) {
		port = "" // omit port if it's the default for the upstream type
	}
	(*data)["upstream_url"] = fmt.Sprintf("%s://%s%s", socket.UpstreamType, config.Hostname, port)

	var hostHeader string
	if config.Hostname != config.HostHeader {
		hostHeader = config.HostHeader
	}
	(*data)["host_header"] = hostHeader

	if len(config.Headers) > 0 {
		headerMap := make(map[string][]string)
		for _, header := range config.Headers {
			if header.Key == "" || header.Value == "" {
				continue
			}
			headerMap[header.Key] = append(headerMap[header.Key], header.Value)
		}
		var headersList []map[string]any
		for key, values := range headerMap {
			headersList = append(headersList, map[string]any{"key": key, "values": values})
		}
		if len(headersList) > 0 {
			(*data)["header"] = headersList
		}
	}

	return nil
}

func connectorFileServerFromUpstreamConfig(data *map[string]any, config *service.FileServerHttpServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with http service type "connector_file_server" but file server HTTP service configuration was not present`)
	}

	(*data)["top_level_directory"] = config.TopLevelDirectory

	return nil
}
