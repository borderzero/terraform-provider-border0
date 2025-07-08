package elasticsearch

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "elasticsearch_configuration" attribute into
// a socket's elasticsearch service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.ElasticsearchServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("elasticsearch_configuration"); ok {
		if configsList := v.([]any); len(configsList) > 0 {
			data = configsList[0].(map[string]any)
		}
	}

	serviceType := service.ElasticsearchServiceTypeStandard // default to "standard"

	if v, ok := data["service_type"]; ok {
		serviceType = v.(string)
	}
	config.ElasticsearchServiceType = serviceType

	switch serviceType {
	case service.ElasticsearchServiceTypeStandard:
		if config.StandardElasticsearchServiceConfiguration == nil {
			config.StandardElasticsearchServiceConfiguration = new(service.StandardElasticsearchServiceConfiguration)
		}
		return standardToUpstreamConfig(data, config.StandardElasticsearchServiceConfiguration)

	default:
		return diag.Errorf(`sockets with elasticsearch service type "%s" not yet supported`, serviceType)
	}
}

func standardToUpstreamConfig(data map[string]any, config *service.StandardElasticsearchServiceConfiguration) diag.Diagnostics {
	authType := service.ElasticsearchAuthenticationTypeBasic // default to "basic"

	if v, ok := data["authentication_type"]; ok {
		authType = v.(string)
	}
	config.AuthenticationType = authType

	if v, ok := data["protocol"]; ok {
		config.Protocol = v.(string)
	}
	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}

	switch authType {
	case service.ElasticsearchAuthenticationTypeBasic:
		if config.BasicAuthentication == nil {
			config.BasicAuthentication = new(service.ElasticsearchServiceTypeBasicAuth)
		}

		if v, ok := data["username"]; ok {
			config.BasicAuthentication.Username = v.(string)
		}
		if v, ok := data["password"]; ok {
			config.BasicAuthentication.Password = v.(string)
		}
	default:
		return diag.Errorf(`elasticsearch authentication type "%s" is invalid`, authType)
	}

	return nil
}
