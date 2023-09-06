package http

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ToUpstreamConfig(d *schema.ResourceData, config *service.HttpServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("database_configuration"); ok {
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
		return standardToUpstreamConfig(data, config.StandardHttpServiceConfiguration)

	case service.HttpServiceTypeConnectorFileServer:
		if config.FileServerHttpServiceConfiguration == nil {
			config.FileServerHttpServiceConfiguration = new(service.FileServerHttpServiceConfiguration)
		}
		return connectorFileServerToUpstreamConfig(data, config.FileServerHttpServiceConfiguration)

	default:
		return diag.Errorf(`sockets with http service type "%s" not yet supported`, httpServiceType)
	}
}

func standardToUpstreamConfig(data map[string]any, config *service.StandardHttpServiceConfiguration) diag.Diagnostics {
	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}
	if v, ok := data["host_header"]; ok {
		config.HostHeader = v.(string)
	}

	return nil
}

func connectorFileServerToUpstreamConfig(data map[string]any, config *service.FileServerHttpServiceConfiguration) diag.Diagnostics {
	if v, ok := data["top_level_directory"]; ok {
		config.TopLevelDirectory = v.(string)
	}

	return nil
}
