package elasticsearch

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's elasticsearch service configuration into terraform resource data for
// the "elasticsearch_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.ElasticsearchServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "elasticsearch" but elasticsearch service configuration was not present`)
	}

	data := map[string]any{
		"service_type": config.ElasticsearchServiceType,
	}

	var diags diag.Diagnostics

	switch config.ElasticsearchServiceType {
	case service.ElasticsearchServiceTypeStandard:
		diags = standardFromUpstreamConfig(&data, config.StandardElasticsearchServiceConfiguration)
	default:
		return diag.Errorf(`sockets with elasticsearch service type "%s" not yet supported`, config.ElasticsearchServiceType)
	}

	if diags.HasError() {
		return diags
	}

	if err := d.Set("elasticsearch_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "elasticsearch_configuration"`)
	}

	return nil
}

func standardFromUpstreamConfig(data *map[string]any, config *service.StandardElasticsearchServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with elasticsearch service type "standard" but standard elasticsearch service configuration was not present`)
	}

	(*data)["authentication_type"] = config.AuthenticationType
	(*data)["protocol"] = config.Protocol
	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port

	switch config.AuthenticationType {
	case service.ElasticsearchAuthenticationTypeBasic:
		if config.BasicAuthentication == nil {
			return diag.Errorf(`got a socket with elasticsearch authentication type "basic" but username and password auth configuration was not present`)
		}
		(*data)["username"] = config.BasicAuthentication.Username
		(*data)["password"] = config.BasicAuthentication.Password
	default:
		return diag.Errorf(`elasticsearch authentication type "%s" is invalid`, config.AuthenticationType)
	}

	return nil
}
