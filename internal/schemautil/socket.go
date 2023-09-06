package schemautil

import (
	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/client/enum"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func FromSocket(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if socket.Tags != nil && len(socket.Tags) > 0 {
		// only set tags if there are any, this prevents a drift in the state
		// if no tags are set in the terraform resource border0_socket
		if err := d.Set("tags", socket.Tags); err != nil {
			return diagnostics.Error(err, "Failed to set tags")
		}
	}

	return SetValues(d, map[string]any{
		"name":                             socket.Name,
		"socket_type":                      socket.SocketType,
		"description":                      socket.Description,
		"upstream_type":                    socket.UpstreamType,
		"upstream_http_hostname":           socket.UpstreamHTTPHostname,
		"recording_enabled":                socket.RecordingEnabled,
		"connector_authentication_enabled": socket.ConnectorAuthenticationEnabled,
	})
}

func FromConnector(d *schema.ResourceData, connectors *border0client.SocketConnectors) diag.Diagnostics {
	var connectorID string

	if len(connectors.List) > 0 {
		connectorID = connectors.List[0].ConnectorID
	}

	return SetValues(d, map[string]any{
		"connector_id": connectorID,
	})
}

func ToSocket(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if v, ok := d.GetOk("description"); ok {
		socket.Description = v.(string)
	}

	if v, ok := d.GetOk("upstream_type"); ok {
		socket.UpstreamType = v.(string)
	} else {
		// if upstream type is not set, use socket type as default upstream type
		// except for database sockets, which use mysql as default upstream type
		switch socket.SocketType {
		case enum.SocketTypeHTTP, enum.SocketTypeSSH, enum.SocketTypeTLS:
			socket.UpstreamType = socket.SocketType
		case enum.SocketTypeDatabase:
			socket.UpstreamType = "mysql"
		}
	}

	if v, ok := d.GetOk("upstream_http_hostname"); ok {
		socket.UpstreamHTTPHostname = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		socket.Tags = make(map[string]string)
		tags := v.(map[string]any)
		for key, value := range tags {
			socket.Tags[key] = value.(string)
		}
	}

	if v, ok := d.GetOk("recording_enabled"); ok {
		socket.RecordingEnabled = v.(bool)
	}

	if v, ok := d.GetOk("connector_authentication_enabled"); ok {
		socket.ConnectorAuthenticationEnabled = v.(bool)
	}

	if v, ok := d.GetOk("connector_id"); ok {
		socket.ConnectorID = v.(string)
	}

	return nil
}
