package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
	border0types "github.com/borderzero/border0-go/service/connector/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSocket() *schema.Resource {
	return &schema.Resource{
		Description:   "The socket resource allows you to create and manage a Border socket.",
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
				Description: "The name of the socket.",
			},
			"socket_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the socket.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the socket.",
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
			"tags": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The tags of the socket.",
			},
			"recording_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if recording is enabled for the socket.",
			},
			"connector_authentication_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if connector authentication is enabled for the socket.",
			},
			"connector_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The connector id that the socket is associated with.",
			},
			"upstream_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The upstream SSH username.",
			},
			"upstream_hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The upstream SSH hostname.",
			},
			"upstream_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The upstream SSH port.",
			},
			"upstream_authentication_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The upstream SSH authentication type.",
			},
		},
	}
}

func resourceSocketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	socket, err := client.Socket(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		log.Printf("[WARN] Socket (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diag.Diagnostics{}
	}
	if err != nil {
		return diagnosticsError(err, "Failed to fetch socket")
	}

	if socket.Tags != nil && len(socket.Tags) > 0 {
		// only set tags if there are any, this prevents a drift in the state
		// if no tags are set in the terraform resource border0_socket
		if err := d.Set("tags", socket.Tags); err != nil {
			return diagnosticsError(err, "Failed to set tags")
		}
	}

	diagnotics := setValues(d, map[string]any{
		"name":                             socket.Name,
		"socket_type":                      socket.SocketType,
		"description":                      socket.Description,
		"upstream_type":                    socket.UpstreamType,
		"upstream_http_hostname":           socket.UpstreamHTTPHostname,
		"recording_enabled":                socket.RecordingEnabled,
		"connector_authentication_enabled": socket.ConnectorAuthenticationEnabled,
	})
	if diagnotics != nil && diagnotics.HasError() {
		return diagnotics
	}

	return injectSocketConnectorDataTo(d, socket)
}

func resourceSocketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	socket := &border0client.Socket{
		Name:       d.Get("name").(string),
		SocketType: d.Get("socket_type").(string),
	}

	populateSocketOptionalFieldsFrom(d, socket)
	populateSocketConnectorDataFrom(d, socket)

	created, err := client.CreateSocket(ctx, socket)
	if err != nil {
		return diagnosticsError(err, "Failed to create key")
	}

	d.SetId(created.SocketID)

	return resourceSocketRead(ctx, d, m)
}

func resourceSocketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)

	if d.HasChangesExcept("socket_type") {
		existingSocket, err := client.Socket(ctx, d.Id())
		if err != nil {
			return diagnosticsError(err, "Failed to fetch socket")
		}
		socketUpdate := &border0client.Socket{
			Name:         d.Get("name").(string),
			SocketType:   d.Get("socket_type").(string),
			UpstreamType: existingSocket.UpstreamType,
		}

		populateSocketOptionalFieldsFrom(d, socketUpdate)
		populateSocketConnectorDataFrom(d, socketUpdate)

		_, err = client.UpdateSocket(ctx, d.Id(), socketUpdate)
		if err != nil {
			return diagnosticsError(err, "Failed to update socket")
		}
	}

	return resourceSocketRead(ctx, d, m)
}

func resourceSocketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	if err := client.DeleteSocket(ctx, d.Id()); err != nil {
		return diagnosticsError(err, "Failed to delete socket")
	}
	d.SetId("")
	return nil
}

func injectSocketConnectorDataTo(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	var (
		connectorID                string
		upstreamUsername           string
		upstreamHostname           string
		upstreamPort               int
		upstreamAuthenticationType string
	)

	if socket.ConnectorData != nil {
		connectorID = socket.ConnectorData.ConnectorID
		switch socket.SocketType {
		case "ssh":
			if socket.ConnectorData.Config != nil {
				upstreamHostname = socket.ConnectorData.Config.Hostname
				upstreamPort = socket.ConnectorData.Config.Port
				if socket.ConnectorData.Config.SSHConfiguration != nil {
					upstreamAuthenticationType = socket.ConnectorData.Config.SSHConfiguration.UpstreamAuthenticationType
					if socket.ConnectorData.Config.SSHConfiguration.BasicCredentials != nil {
						upstreamUsername = socket.ConnectorData.Config.SSHConfiguration.BasicCredentials.Username
					}
				}
			}
		}
	}

	return setValues(d, map[string]any{
		"connector_id":                 connectorID,
		"upstream_username":            upstreamUsername,
		"upstream_hostname":            upstreamHostname,
		"upstream_port":                upstreamPort,
		"upstream_authentication_type": upstreamAuthenticationType,
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
		case "http", "ssh", "tls":
			socket.UpstreamType = socket.SocketType
		case "database":
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
}

func populateSocketConnectorDataFrom(d *schema.ResourceData, socket *border0client.Socket) {
	if v, ok := d.GetOk("connector_id"); ok {
		if socket.ConnectorData == nil {
			socket.ConnectorData = new(border0client.SocketConnectorData)
		}
		socket.ConnectorData.ConnectorID = v.(string)

		// TODO: need to validate input data
		if socket.ConnectorData.Config == nil {
			socket.ConnectorData.Config = new(border0types.ConnectorServiceUpstreamConfig)
		}
		switch socket.SocketType {
		case "ssh":
			socket.ConnectorData.Config.UpstreamConnectionType = border0types.UpstreamConnectionTypeSSH
			if v, ok := d.GetOk("upstream_hostname"); ok {
				socket.ConnectorData.Config.Hostname = v.(string)
			}

			if v, ok := d.GetOk("upstream_port"); ok {
				socket.ConnectorData.Config.Port = v.(int)
			}

			if socket.ConnectorData.Config.SSHConfiguration == nil {
				socket.ConnectorData.Config.SSHConfiguration = new(border0types.SSHConfiguration)
			}
			if v, ok := d.GetOk("upstream_authentication_type"); ok {
				socket.ConnectorData.Config.SSHConfiguration.UpstreamAuthenticationType = v.(string)
			}

			if socket.ConnectorData.Config.SSHConfiguration.BasicCredentials == nil {
				socket.ConnectorData.Config.SSHConfiguration.BasicCredentials = new(border0types.BasicCredentials)
			}
			if v, ok := d.GetOk("upstream_username"); ok {
				socket.ConnectorData.Config.SSHConfiguration.BasicCredentials.Username = v.(string)
			}
		}
	}
}
