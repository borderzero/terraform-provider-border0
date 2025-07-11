package schemautil

import (
	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/client/enum"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromSocket converts a border0client.Socket's top level fields into terraform resource data.
// Those fields are:
// - name
// - socket_type
// - description
// - upstream_type
// - recording_enabled
func FromSocket(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if len(socket.Tags) > 0 {
		// only set tags if there are any, this prevents a drift in the state
		// if no tags are set in the terraform resource border0_socket
		if err := d.Set("tags", socket.Tags); err != nil {
			return diagnostics.Error(err, "Failed to set tags")
		}
	}

	return SetValues(d, map[string]any{
		"name":              socket.Name,
		"socket_type":       socket.SocketType,
		"description":       socket.Description,
		"upstream_type":     socket.UpstreamType,
		"recording_enabled": socket.RecordingEnabled,
	})
}

// FromConnector reads the `connector_id` from the first connector in the list of connectors from a Border0
// socket, and sets the `connector_id` in the terraform resource data.
func FromConnector(d *schema.ResourceData, connectors *border0client.SocketConnectors) diag.Diagnostics {
	var connectorID string

	if len(connectors.List) > 0 {
		connectorID = connectors.List[0].ConnectorID
	}

	return SetValues(d, map[string]any{
		"connector_id": connectorID,
	})
}

// ToSocket read top level socket fields from terraform resource data and sets them in a border0client.Socket.
// Those fields are:
// - description
// - upstream_type
// - upstream_http_hostname
// - tags
// - recording_enabled
// - connector_id
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
		case
			enum.SocketTypeHTTP,
			enum.SocketTypeSSH,
			enum.SocketTypeTLS,
			enum.SocketTypeVNC,
			enum.SocketTypeRDP,
			enum.SocketTypeVPN:

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

	if v, ok := d.GetOk("connector_id"); ok {
		socket.ConnectorID = v.(string)
	}

	return nil
}
