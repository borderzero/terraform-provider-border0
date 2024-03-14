package tls

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig converts a socket's TLS service configuration into terraform resource data for
// the "tls_configuration" attribute on the "border0_socket" resource.
func FromUpstreamConfig(d *schema.ResourceData, config *service.TlsServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with service type "tls" but TLS service configuration was not present`)
	}

	data := map[string]any{
		"service_type": config.TlsServiceType,
	}

	var diags diag.Diagnostics

	switch config.TlsServiceType {
	case service.TlsServiceTypeStandard:
		diags = standardFromUpstreamConfig(&data, config.StandardTlsServiceConfiguration)
	default:
		return diag.Errorf(`sockets with TLS service type "%s" not yet supported`, config.TlsServiceType)
	}

	if diags.HasError() {
		return diags
	}

	if err := d.Set("tls_configuration", []map[string]any{data}); err != nil {
		return diagnostics.Error(err, `Failed to set "tls_configuration"`)
	}

	return nil
}

func standardFromUpstreamConfig(data *map[string]any, config *service.StandardTlsServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with TLS service type "standard" but standard TLS service configuration was not present`)
	}

	(*data)["hostname"] = config.Hostname
	(*data)["port"] = config.Port

	return nil
}
