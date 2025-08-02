package schemautil

import (
	"sort"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/client/enum"
	"github.com/borderzero/border0-go/lib/types/slice"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FromSocket converts a border0client.Socket's top level fields into terraform resource data.
// Those fields are:
// - name
// - display_name
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
		"display_name":      socket.DisplayName,
		"socket_type":       socket.SocketType,
		"description":       socket.Description,
		"upstream_type":     socket.UpstreamType,
		"recording_enabled": socket.RecordingEnabled,
	})
}

// FromConnector reads the `connector_id` from the first connector in the list of connectors from a Border0
// socket, and sets the `connector_id` in the terraform resource data.
func FromConnector(d *schema.ResourceData, connectors *border0client.SocketConnectors) diag.Diagnostics {
	connectorIDs := slice.Transform(connectors.List, func(c border0client.SocketConnector) string { return c.ConnectorID })
	sort.Strings(connectorIDs)
	connectorIDsAny := slice.Transform(connectorIDs, func(stringID string) any { return stringID })

	// only set connector_ids if the original state had connector_ids
	if _, ok := d.GetOk("connector_ids"); ok {
		return SetValues(d, map[string]any{"connector_ids": schema.NewSet(schema.HashString, connectorIDsAny)})
	}
	// backwards compatibility, can remove next major version
	if _, ok := d.GetOk("connector_id"); ok {
		if len(connectorIDs) > 0 {
			return SetValues(d, map[string]any{"connector_id": connectorIDs[0]})
		}
	}
	return nil
}

// ToSocket read top level socket fields from terraform resource data and sets them in a border0client.Socket.
// Those fields are:
// - display_name
// - description
// - upstream_type
// - upstream_http_hostname
// - tags
// - recording_enabled
// - connector_id
func ToSocket(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if v, ok := d.GetOk("display_name"); ok {
		socket.DisplayName = v.(string)
	}

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

	if v, ok := d.GetOk("connector_ids"); ok {
		rawSet := v.(*schema.Set)
		socket.ConnectorIDs = make([]string, rawSet.Len())
		for i, id := range rawSet.List() {
			socket.ConnectorIDs[i] = id.(string)
		}
	}

	// backwards compatibility, will be deprecated in next major release
	if len(socket.ConnectorIDs) == 0 {
		if v, ok := d.GetOk("connector_id"); ok {
			socket.ConnectorIDs = []string{v.(string)}
		}
	}
	return nil
}
