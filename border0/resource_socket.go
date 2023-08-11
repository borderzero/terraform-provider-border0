package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
	border0enum "github.com/borderzero/border0-go/client/enum"
	border0types "github.com/borderzero/border0-go/service/connector/types"
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
			"upstream_connection_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The upstream connection type. Valid values: `ssh`, `aws_ec2_connect`, `aws_ssm`, `database`. Defaults to `ssh`.",
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
				Description: "The upstream authentication type. Valid values: `username_password`, `border0_cert`, `ssh_private_key`. Defaults to `border0_cert`.",
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
	client := m.(*border0client.APIClient)

	// get socket details
	socket, err := client.Socket(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		log.Printf("[WARN] Socket (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diag.Diagnostics{}
	}
	if err != nil {
		return diagnosticsError(err, "Failed to fetch socket")
	}
	// get socket linked connectors and their ids
	connectors, err := client.SocketConnectors(ctx, d.Id())
	if err != nil {
		return diagnosticsError(err, "Failed to fetch socket connectors")
	}
	// get socket upstream configs for linked connectors
	upstreamConfigs, err := client.SocketUpstreamConfigs(ctx, d.Id())
	if err != nil {
		return diagnosticsError(err, "Failed to fetch socket upstream configs")
	}

	// now inject the socket details, connector id and upstream config into the resource data
	if diagnotics := injectSocketFieldsTo(d, socket); diagnotics != nil && diagnotics.HasError() {
		return diagnotics
	}
	if diagnotics := injectSocketConnectorIDTo(d, connectors); diagnotics != nil && diagnotics.HasError() {
		return diagnotics
	}
	return injectSocketUpstreamConfigTo(d, socket, upstreamConfigs)
}

func resourceSocketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	socket := &border0client.Socket{
		Name:       d.Get("name").(string),
		SocketType: d.Get("socket_type").(string),
	}

	populateSocketOptionalFieldsFrom(d, socket)
	populateSocketUpstreamConfigFrom(d, socket)

	created, err := client.CreateSocket(ctx, socket)
	if err != nil {
		return diagnosticsError(err, "Failed to create socket")
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
		populateSocketUpstreamConfigFrom(d, socketUpdate)

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

func injectSocketFieldsTo(d *schema.ResourceData, socket *border0client.Socket) diag.Diagnostics {
	if socket.Tags != nil && len(socket.Tags) > 0 {
		// only set tags if there are any, this prevents a drift in the state
		// if no tags are set in the terraform resource border0_socket
		if err := d.Set("tags", socket.Tags); err != nil {
			return diagnosticsError(err, "Failed to set tags")
		}
	}

	return setValues(d, map[string]any{
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

	return setValues(d, map[string]any{
		"connector_id": connectorID,
	})
}

// TODO: refactor this function
func injectSocketUpstreamConfigTo(d *schema.ResourceData, socket *border0client.Socket, configs *border0client.SocketUpstreamConfigs) diag.Diagnostics {
	var (
		upstreamHostname           string
		upstreamPort               int
		upstreamConnectionType     string
		upstreamAuthenticationType string
		upstreamUsername           string
		upstreamPassword           string
		upstreamPrivateKey         string
	)

	if len(configs.List) > 0 {
		config := configs.List[0].Config

		upstreamConnectionType = config.UpstreamConnectionType
		upstreamHostname = config.Hostname
		upstreamPort = config.Port

		switch socket.SocketType {
		case border0enum.SocketTypeSSH:
			if config.SSHConfiguration != nil {
				switch config.UpstreamConnectionType {
				case border0types.UpstreamConnectionTypeSSH:
					upstreamAuthenticationType = config.SSHConfiguration.UpstreamAuthenticationType

					switch config.SSHConfiguration.UpstreamAuthenticationType {
					case border0types.UpstreamAuthenticationTypeUsernamePassword:
						if config.SSHConfiguration.BasicCredentials != nil {
							upstreamUsername = config.SSHConfiguration.BasicCredentials.Username
							upstreamPassword = config.SSHConfiguration.BasicCredentials.Password
						}
					case border0types.UpstreamAuthenticationTypeBorder0Cert:
						if config.SSHConfiguration.Border0CertificateDetails != nil {
							upstreamUsername = config.SSHConfiguration.Border0CertificateDetails.Username
						}
					case border0types.UpstreamAuthenticationTypeSSHPrivateKey:
						if config.SSHConfiguration.SSHPrivateKeyDetails != nil {
							upstreamUsername = config.SSHConfiguration.SSHPrivateKeyDetails.Username
							upstreamPrivateKey = config.SSHConfiguration.SSHPrivateKeyDetails.Key
						}
					}
				}
			}
		}
	}

	return setValues(d, map[string]any{
		"upstream_connection_type":     upstreamConnectionType,
		"upstream_hostname":            upstreamHostname,
		"upstream_port":                upstreamPort,
		"upstream_authentication_type": upstreamAuthenticationType,
		"upstream_username":            upstreamUsername,
		"upstream_password":            upstreamPassword,
		"upstream_private_key":         upstreamPrivateKey,
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
		case border0enum.SocketTypeHTTP, border0enum.SocketTypeSSH, border0enum.SocketTypeTLS:
			socket.UpstreamType = socket.SocketType
		case border0enum.SocketTypeDatabase:
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

// TODO: need to validate input data
// TODO: refactor this function
func populateSocketUpstreamConfigFrom(d *schema.ResourceData, socket *border0client.Socket) {
	// noop if connector id is not set or empty
	if v, ok := d.GetOk("connector_id"); !ok || v.(string) == "" {
		return
	}

	// if connector id is given, upstream config is required
	if socket.UpstreamConfig == nil {
		socket.UpstreamConfig = new(border0types.ConnectorServiceUpstreamConfig)
	}
	if v, ok := d.GetOk("upstream_connection_type"); ok {
		socket.UpstreamConfig.UpstreamConnectionType = v.(string)
	}
	if v, ok := d.GetOk("upstream_hostname"); ok {
		socket.UpstreamConfig.Hostname = v.(string)
	}
	if v, ok := d.GetOk("upstream_port"); ok {
		socket.UpstreamConfig.Port = v.(int)
	}

	switch socket.SocketType {
	case border0enum.SocketTypeSSH:
		if socket.UpstreamConfig.UpstreamConnectionType == "" {
			socket.UpstreamConfig.UpstreamConnectionType = border0types.UpstreamConnectionTypeSSH
		}
		if socket.UpstreamConfig.SSHConfiguration == nil {
			socket.UpstreamConfig.SSHConfiguration = new(border0types.SSHConfiguration)
		}

		switch socket.UpstreamConfig.UpstreamConnectionType {
		case border0types.UpstreamConnectionTypeSSH:
			if v, ok := d.GetOk("upstream_authentication_type"); ok {
				socket.UpstreamConfig.SSHConfiguration.UpstreamAuthenticationType = v.(string)
			}
			if socket.UpstreamConfig.SSHConfiguration.UpstreamAuthenticationType == "" {
				socket.UpstreamConfig.SSHConfiguration.UpstreamAuthenticationType = border0types.UpstreamAuthenticationTypeBorder0Cert
			}

			switch socket.UpstreamConfig.SSHConfiguration.UpstreamAuthenticationType {
			case border0types.UpstreamAuthenticationTypeUsernamePassword:
				if socket.UpstreamConfig.SSHConfiguration.BasicCredentials == nil {
					socket.UpstreamConfig.SSHConfiguration.BasicCredentials = new(border0types.BasicCredentials)
				}
				if v, ok := d.GetOk("upstream_username"); ok {
					socket.UpstreamConfig.SSHConfiguration.BasicCredentials.Username = v.(string)
				}
				if v, ok := d.GetOk("upstream_password"); ok {
					socket.UpstreamConfig.SSHConfiguration.BasicCredentials.Password = v.(string)
				}
			case border0types.UpstreamAuthenticationTypeBorder0Cert:
				if socket.UpstreamConfig.SSHConfiguration.Border0CertificateDetails == nil {
					socket.UpstreamConfig.SSHConfiguration.Border0CertificateDetails = new(border0types.Border0CertificateDetails)
				}
				if v, ok := d.GetOk("upstream_username"); ok {
					socket.UpstreamConfig.SSHConfiguration.Border0CertificateDetails.Username = v.(string)
				}
			case border0types.UpstreamAuthenticationTypeSSHPrivateKey:
				if socket.UpstreamConfig.SSHConfiguration.SSHPrivateKeyDetails == nil {
					socket.UpstreamConfig.SSHConfiguration.SSHPrivateKeyDetails = new(border0types.SSHPrivateKeyDetails)
				}
				if v, ok := d.GetOk("upstream_username"); ok {
					socket.UpstreamConfig.SSHConfiguration.SSHPrivateKeyDetails.Username = v.(string)
				}
				if v, ok := d.GetOk("upstream_private_key"); ok {
					socket.UpstreamConfig.SSHConfiguration.SSHPrivateKeyDetails.Key = v.(string)
				}
			}
		}
	}
}
