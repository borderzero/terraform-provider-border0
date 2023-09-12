package schemautil

import (
	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/database"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/http"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/ssh"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/tls"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromUpstreamConfig translates a socket's upstream service config to terraform resource data.
// In short: *border0client.SocketUpstreamConfigs -> *schema.ResourceData
func FromUpstreamConfig(
	d *schema.ResourceData,
	socket *border0client.Socket,
	configs *border0client.SocketUpstreamConfigs,
) diag.Diagnostics {
	// noop if upstream config is not set
	if len(configs.List) == 0 {
		return nil
	}

	firstConfig := configs.List[0]
	config := firstConfig.Config

	switch config.ServiceType {
	case service.ServiceTypeSsh:
		return ssh.FromUpstreamConfig(d, config.SshServiceConfiguration)
	case service.ServiceTypeDatabase:
		return database.FromUpstreamConfig(d, config.DatabaseServiceConfiguration)
	case service.ServiceTypeHttp:
		return http.FromUpstreamConfig(d, socket, config.HttpServiceConfiguration)
	case service.ServiceTypeTls:
		return tls.FromUpstreamConfig(d, config.TlsServiceConfiguration)
	default:
		return diag.Errorf(`sockets with service type "%s" not yet supported`, config.ServiceType)
	}
}

// ToUpstreamConfig translates terraform resource data to a socket's upstream service config.
// In short: *schema.ResourceData -> *border0client.SocketUpstreamConfigs
func ToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	// noop if connector id is not set or empty
	if v, ok := d.GetOk("connector_id"); !ok || v.(string) == "" {
		return nil
	}

	// if connector id is given, upstream config is required
	if socket.UpstreamConfig == nil {
		socket.UpstreamConfig = new(service.Configuration)
	}

	socket.UpstreamConfig.ServiceType = socket.SocketType // default to socket type

	var diags diag.Diagnostics

	switch socket.UpstreamConfig.ServiceType {
	case service.ServiceTypeSsh:
		if socket.UpstreamConfig.SshServiceConfiguration == nil {
			socket.UpstreamConfig.SshServiceConfiguration = new(service.SshServiceConfiguration)
		}
		diags = ssh.ToUpstreamConfig(d, socket.UpstreamConfig.SshServiceConfiguration)

	case service.ServiceTypeDatabase:
		if socket.UpstreamConfig.DatabaseServiceConfiguration == nil {
			socket.UpstreamConfig.DatabaseServiceConfiguration = new(service.DatabaseServiceConfiguration)
		}
		diags = database.ToUpstreamConfig(d, socket.UpstreamConfig.DatabaseServiceConfiguration)

	case service.ServiceTypeHttp:
		if socket.UpstreamConfig.HttpServiceConfiguration == nil {
			socket.UpstreamConfig.HttpServiceConfiguration = new(service.HttpServiceConfiguration)
		}
		diags = http.ToUpstreamConfig(d, socket, socket.UpstreamConfig.HttpServiceConfiguration)

	case service.ServiceTypeTls:
		if socket.UpstreamConfig.TlsServiceConfiguration == nil {
			socket.UpstreamConfig.TlsServiceConfiguration = new(service.TlsServiceConfiguration)
		}
		diags = tls.ToUpstreamConfig(d, socket.UpstreamConfig.TlsServiceConfiguration)

	default:
		return diag.Errorf(`sockets with service type "%s" not yet supported`, socket.UpstreamConfig.ServiceType)
	}

	if err := socket.UpstreamConfig.Validate(); err != nil {
		return diagnostics.Error(err, "Upstream configuration is invalid")
	}

	return diags
}
