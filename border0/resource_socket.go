package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/client/enum"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSocket() *schema.Resource {
	return &schema.Resource{
		Description:   "The socket resource allows you to create and manage a Border0 socket.",
		ReadContext:   resourceSocketRead,
		CreateContext: resourceSocketCreate,
		UpdateContext: resourceSocketUpdate,
		DeleteContext: resourceSocketDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the socket. Must be unique within your Border0 organization. Socket name can have alphanumerics and hyphens, but it must start or end with alphanumeric.",
			},
			"socket_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the socket. Valid values: `ssh`, `http`, `database`, `tls`.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the socket.",
			},
			"recording_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if session recording is enabled for the socket.",
			},
			"connector_authentication_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if connector authentication is enabled for the socket.",
			},
			"tags": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The tags of the socket.",
			},
			"connector_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The connector id that the socket is associated with.",
			},
			"upstream_type": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" && new == "" {
						return true
					}
					return old == new
				},
				Description: "The upstream type of the socket.",
			},
			"upstream_http_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The upstream http hostname of the socket.",
			},
			"upstream_service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The upstream service type. Valid values depend on the socket type, for ssh: `standard`, `aws_ec2_instance_connect`, `aws_ssm`. Defaults to `standard`.",
			},
			"upstream_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The upstream hostname.",
			},
			"upstream_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The upstream port number.",
			},
			"upstream_authentication_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The upstream authentication type. Valid values: `username_password`, `border0_certificate`, `ssh_private_key`. Defaults to `border0_certificate`.",
			},
			"upstream_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The upstream username.",
			},
			"upstream_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The upstream password.",
			},
			"upstream_private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The upstream private key.",
			},
		},
	}
}

func resourceSocketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	// get socket details
	socket, err := client.Socket(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		log.Printf("[WARN] Socket (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diag.Diagnostics{}
	}
	if err != nil {
		return schemautil.DiagnosticsError(err, "Failed to fetch socket")
	}
	// get socket linked connectors and their ids
	connectors, err := client.SocketConnectors(ctx, d.Id())
	if err != nil {
		return schemautil.DiagnosticsError(err, "Failed to fetch socket connectors")
	}
	// get socket upstream configs for linked connectors
	upstreamConfigs, err := client.SocketUpstreamConfigs(ctx, d.Id())
	if err != nil {
		return schemautil.DiagnosticsError(err, "Failed to fetch socket upstream configs")
	}

	// now inject the socket details, connector id and upstream config into the resource data
	if diagnostics := injectSocketFieldsTo(d, socket); diagnostics.HasError() {
		return diagnostics
	}
	if diagnostics := injectSocketConnectorIDTo(d, connectors); diagnostics.HasError() {
		return diagnostics
	}
	return schemautil.FromUpstreamConfig(d, upstreamConfigs)
}

func resourceSocketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	socket := &border0client.Socket{
		Name:       d.Get("name").(string),
		SocketType: d.Get("socket_type").(string),
	}

	populateSocketOptionalFieldsFrom(d, socket)

	if diagnostics := schemautil.ToUpstreamConfig(d, socket); diagnostics.HasError() {
		return diagnostics
	}

	created, err := client.CreateSocket(ctx, socket)
	if err != nil {
		return schemautil.DiagnosticsError(err, "Failed to create socket")
	}

	d.SetId(created.SocketID)

	return resourceSocketRead(ctx, d, m)
}

func resourceSocketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	if d.HasChangesExcept("socket_type") {
		existingSocket, err := client.Socket(ctx, d.Id())
		if err != nil {
			return schemautil.DiagnosticsError(err, "Failed to fetch socket")
		}
		socketUpdate := &border0client.Socket{
			Name:         d.Get("name").(string),
			SocketType:   d.Get("socket_type").(string),
			UpstreamType: existingSocket.UpstreamType,
		}

		populateSocketOptionalFieldsFrom(d, socketUpdate)

		if diagnostics := schemautil.ToUpstreamConfig(d, socketUpdate); diagnostics.HasError() {
			return diagnostics
		}

		_, err = client.UpdateSocket(ctx, d.Id(), socketUpdate)
		if err != nil {
			return schemautil.DiagnosticsError(err, "Failed to update socket")
		}
	}

	return resourceSocketRead(ctx, d, m)
}

func resourceSocketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	if err := client.DeleteSocket(ctx, d.Id()); err != nil {
		return schemautil.DiagnosticsError(err, "Failed to delete socket")
	}
	d.SetId("")
	return nil
}

func injectSocketFieldsTo(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if socket.Tags != nil && len(socket.Tags) > 0 {
		// only set tags if there are any, this prevents a drift in the state
		// if no tags are set in the terraform resource border0_socket
		if err := d.Set("tags", socket.Tags); err != nil {
			return schemautil.DiagnosticsError(err, "Failed to set tags")
		}
	}

	return schemautil.SetValues(d, map[string]any{
		"name":                             socket.Name,
		"socket_type":                      socket.SocketType,
		"description":                      socket.Description,
		"upstream_type":                    socket.UpstreamType,
		"upstream_http_hostname":           socket.UpstreamHTTPHostname,
		"recording_enabled":                socket.RecordingEnabled,
		"connector_authentication_enabled": socket.ConnectorAuthenticationEnabled,
	})
}

func injectSocketConnectorIDTo(d *schema.ResourceData, connectors *border0client.SocketConnectors) diag.Diagnostics {
	var connectorID string

	if len(connectors.List) > 0 {
		connectorID = connectors.List[0].ConnectorID
	}

	return schemautil.SetValues(d, map[string]any{
		"connector_id": connectorID,
	})
}

func populateSocketOptionalFieldsFrom(d *schema.ResourceData, socket *border0client.Socket) {
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
}
