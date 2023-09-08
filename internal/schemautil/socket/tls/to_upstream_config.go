package tls

import (
	"github.com/borderzero/border0-go/types/service"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ToUpstreamConfig converts a "border0_socket" resource data from "tls_configuration" attribute into
// a socket's TLS service upstream configuration.
func ToUpstreamConfig(d *schema.ResourceData, config *service.TlsServiceConfiguration) diag.Diagnostics {
	data := make(map[string]any)

	if v, ok := d.GetOk("tls_configuration"); ok {
		if tlsConfigsList := v.([]any); len(tlsConfigsList) > 0 {
			data = tlsConfigsList[0].(map[string]any)
		}
	}

	tlsServiceType := service.TlsServiceTypeStandard // default to "standard"

	if v, ok := data["service_type"]; ok {
		tlsServiceType = v.(string)
	}
	config.TlsServiceType = tlsServiceType

	switch tlsServiceType {
	case service.TlsServiceTypeStandard:
		if config.StandardTlsServiceConfiguration == nil {
			config.StandardTlsServiceConfiguration = new(service.StandardTlsServiceConfiguration)
		}
		return standardToUpstreamConfig(data, config.StandardTlsServiceConfiguration)

	case service.TlsServiceTypeVpn:
		if config.VpnTlsServiceConfiguration == nil {
			config.VpnTlsServiceConfiguration = new(service.VpnTlsServiceConfiguration)
		}
		return vpnToUpstreamConfig(data, config.VpnTlsServiceConfiguration)

	case service.TlsServiceTypeHttpProxy:
		if config.HttpProxyTlsServiceConfiguration == nil {
			config.HttpProxyTlsServiceConfiguration = new(service.HttpProxyTlsServiceConfiguration)
		}
		return httpProxyToUpstreamConfig(data, config.HttpProxyTlsServiceConfiguration)

	default:
		return diag.Errorf(`sockets with tls service type "%s" not yet supported`, tlsServiceType)
	}
}

func standardToUpstreamConfig(data map[string]any, config *service.StandardTlsServiceConfiguration) diag.Diagnostics {
	if v, ok := data["hostname"]; ok {
		config.Hostname = v.(string)
	}
	if v, ok := data["port"]; ok {
		config.Port = uint16(v.(int))
	}

	return nil
}

func vpnToUpstreamConfig(data map[string]any, config *service.VpnTlsServiceConfiguration) diag.Diagnostics {
	if v, ok := data["vpn_subnet"]; ok {
		config.VpnSubnet = v.(string)
	}
	if v, ok := data["vpn_routes"]; ok {
		config.Routes = v.([]string)
	}

	return nil
}

func httpProxyToUpstreamConfig(data map[string]any, config *service.HttpProxyTlsServiceConfiguration) diag.Diagnostics {
	if v, ok := data["http_proxy_host_allowlist"]; ok {
		config.HostAllowlist = v.([]string)
	}

	return nil
}
