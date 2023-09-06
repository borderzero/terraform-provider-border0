package tls

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	case service.TlsServiceTypeVpn:
		diags = vpnFromUpstreamConfig(&data, config.VpnTlsServiceConfiguration)
	case service.TlsServiceTypeHttpProxy:
		diags = httpProxyFromUpstreamConfig(&data, config.HttpProxyTlsServiceConfiguration)
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

func vpnFromUpstreamConfig(data *map[string]any, config *service.VpnTlsServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with TLS service type "vpn" but VPN TLS service configuration was not present`)
	}

	(*data)["vpn_subnet"] = config.VpnSubnet
	(*data)["vpn_routes"] = config.Routes

	return nil
}

func httpProxyFromUpstreamConfig(data *map[string]any, config *service.HttpProxyTlsServiceConfiguration) diag.Diagnostics {
	if config == nil {
		return diag.Errorf(`got a socket with TLS service type "http_proxy" but HTTP proxy TLS service configuration was not present`)
	}

	(*data)["http_proxy_host_allowlist"] = config.HostAllowlist

	return nil
}
