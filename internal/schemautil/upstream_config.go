package schemautil

import (
	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/aws_s3"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/database"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/elasticsearch"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/exit_node"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/http"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/kubernetes"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/rdp"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/snowflake"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/ssh"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/subnet_router"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/tls"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/vnc"
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
	case service.ServiceTypeVnc:
		return vnc.FromUpstreamConfig(d, config.VncServiceConfiguration)
	case service.ServiceTypeRdp:
		return rdp.FromUpstreamConfig(d, config.RdpServiceConfiguration)
	case service.ServiceTypeSubnetRouter:
		return subnet_router.FromUpstreamConfig(d, config.SubnetRouterServiceConfiguration)
	case service.ServiceTypeExitNode:
		return exit_node.FromUpstreamConfig(d, config.ExitNodeServiceConfiguration)
	case service.ServiceTypeSnowflake:
		return snowflake.FromUpstreamConfig(d, config.SnowflakeServiceConfiguration)
	case service.ServiceTypeKubernetes:
		return kubernetes.FromUpstreamConfig(d, config.KubernetesServiceConfiguration)
	case service.ServiceTypeElasticsearch:
		return elasticsearch.FromUpstreamConfig(d, config.ElasticsearchServiceConfiguration)
	case service.ServiceTypeAwsS3:
		return aws_s3.FromUpstreamConfig(d, config.AwsS3ServiceConfiguration)
	default:
		return diag.Errorf(`sockets with service type "%s" not yet supported`, config.ServiceType)
	}
}

// ToUpstreamConfig translates terraform resource data to a socket's upstream service config.
// In short: *schema.ResourceData -> *border0client.SocketUpstreamConfigs
func ToUpstreamConfig(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	// noop if connector ids is not set or empty
	// NOTE: the connector_id check can be removed when we remove the field from the schema (next major release).
	connectorIDAny, connectorIDFound := d.GetOk("connector_id")
	connectorIDsAny, connectorIDsFound := d.GetOk("connector_ids")
	if (!connectorIDFound || connectorIDAny.(string) == "") && (!connectorIDsFound || connectorIDsAny.(*schema.Set).Len() == 0) {
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

	case service.ServiceTypeVnc:
		if socket.UpstreamConfig.VncServiceConfiguration == nil {
			socket.UpstreamConfig.VncServiceConfiguration = new(service.VncServiceConfiguration)
		}
		diags = vnc.ToUpstreamConfig(d, socket.UpstreamConfig.VncServiceConfiguration)

	case service.ServiceTypeRdp:
		if socket.UpstreamConfig.RdpServiceConfiguration == nil {
			socket.UpstreamConfig.RdpServiceConfiguration = new(service.RdpServiceConfiguration)
		}
		diags = rdp.ToUpstreamConfig(d, socket.UpstreamConfig.RdpServiceConfiguration)

	case service.ServiceTypeSubnetRouter:
		if socket.UpstreamConfig.SubnetRouterServiceConfiguration == nil {
			socket.UpstreamConfig.SubnetRouterServiceConfiguration = new(service.SubnetRouterServiceConfiguration)
		}
		diags = subnet_router.ToUpstreamConfig(d, socket.UpstreamConfig.SubnetRouterServiceConfiguration)

	case service.ServiceTypeExitNode:
		if socket.UpstreamConfig.ExitNodeServiceConfiguration == nil {
			socket.UpstreamConfig.ExitNodeServiceConfiguration = new(service.ExitNodeServiceConfiguration)
		}
		diags = exit_node.ToUpstreamConfig(d, socket.UpstreamConfig.ExitNodeServiceConfiguration)

	case service.ServiceTypeSnowflake:
		if socket.UpstreamConfig.SnowflakeServiceConfiguration == nil {
			socket.UpstreamConfig.SnowflakeServiceConfiguration = new(service.SnowflakeServiceConfiguration)
		}
		diags = snowflake.ToUpstreamConfig(d, socket.UpstreamConfig.SnowflakeServiceConfiguration)

	case service.ServiceTypeKubernetes:
		if socket.UpstreamConfig.KubernetesServiceConfiguration == nil {
			socket.UpstreamConfig.KubernetesServiceConfiguration = new(service.KubernetesServiceConfiguration)
		}
		diags = kubernetes.ToUpstreamConfig(d, socket.UpstreamConfig.KubernetesServiceConfiguration)
	case service.ServiceTypeElasticsearch:
		if socket.UpstreamConfig.ElasticsearchServiceConfiguration == nil {
			socket.UpstreamConfig.ElasticsearchServiceConfiguration = new(service.ElasticsearchServiceConfiguration)
		}
		diags = elasticsearch.ToUpstreamConfig(d, socket.UpstreamConfig.ElasticsearchServiceConfiguration)

	case service.ServiceTypeAwsS3:
		if socket.UpstreamConfig.AwsS3ServiceConfiguration == nil {
			socket.UpstreamConfig.AwsS3ServiceConfiguration = new(service.AwsS3ServiceConfiguration)
		}
		diags = aws_s3.ToUpstreamConfig(d, socket.UpstreamConfig.AwsS3ServiceConfiguration)
	default:
		return diag.Errorf(`sockets with service type "%s" not yet supported`, socket.UpstreamConfig.ServiceType)
	}

	if err := socket.UpstreamConfig.Validate(); err != nil {
		return diagnostics.Error(err, "Upstream configuration is invalid")
	}

	return diags
}
