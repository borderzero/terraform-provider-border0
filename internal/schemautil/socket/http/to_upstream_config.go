package http

import (
	"regexp"
	"strconv"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "http_configuration" attribute into
// a socket's HTTP service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket, config *service.HttpServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("http_configuration"); ok {
		if httpConfigsList := v.([]any); len(httpConfigsList) > 0 {
			data = httpConfigsList[0].(map[string]any)
		}
	}

	httpServiceType := service.HttpServiceTypeStandard // default to "standard"

	if v, ok := data["service_type"]; ok {
		httpServiceType = v.(string)
	}
	config.HttpServiceType = httpServiceType

	switch httpServiceType {
	case service.HttpServiceTypeStandard:
		if config.StandardHttpServiceConfiguration == nil {
			config.StandardHttpServiceConfiguration = new(service.StandardHttpServiceConfiguration)
		}
		return standardToUpstreamConfig(data, socket, config.StandardHttpServiceConfiguration)

	case service.HttpServiceTypeConnectorFileServer:
		if config.FileServerHttpServiceConfiguration == nil {
			config.FileServerHttpServiceConfiguration = new(service.FileServerHttpServiceConfiguration)
		}
		return connectorFileServerToUpstreamConfig(data, config.FileServerHttpServiceConfiguration)

	default:
		return diag.Errorf(`sockets with http service type "%s" not yet supported`, httpServiceType)
	}
}

func standardToUpstreamConfig(data map[string]any, socket *border0client.Socket, config *service.StandardHttpServiceConfiguration) diag.Diagnostics {
	if v, ok := data["upstream_url"]; ok {
		upstreamURL := v.(string)

		re := regexp.MustCompile(`^(https?):\/\/([^:\/\s]+)(?::(\d+))?`)
		matches := re.FindStringSubmatch(upstreamURL)

		if len(matches) == 0 {
			return diag.Errorf("invalid upstream URL: %s", upstreamURL)
		}

		// matches[0] is the whole string

		upstreamType := matches[1]
		hostname := matches[2]
		port := uint16(0) // Default value

		switch upstreamType {
		case "http":
			port = 80
		case "https":
			port = 443
		default:
			// should never happen
		}

		if len(matches) == 4 && matches[3] != "" {
			portUint64, err := strconv.ParseUint(matches[3], 10, 16)
			if err != nil {
				return diag.Errorf("cannot parse port from the upstream URL: %s", upstreamURL)
			}
			port = uint16(portUint64)
		}

		config.Hostname = hostname
		config.Port = port
		config.HostHeader = hostname

		// backward compatibility with the old "upstream_type" and "upstream_http_hostname" attributes
		socket.UpstreamType = upstreamType
		socket.UpstreamHTTPHostname = hostname
	}

	if v, ok := data["host_header"]; ok {
		hostHeader := v.(string)
		if hostHeader != "" {
			config.HostHeader = hostHeader
		}
	}

	return nil
}

func connectorFileServerToUpstreamConfig(data map[string]any, config *service.FileServerHttpServiceConfiguration) diag.Diagnostics {
	if v, ok := data["top_level_directory"]; ok {
		config.TopLevelDirectory = v.(string)
	}

	return nil
}
