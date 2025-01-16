package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/socket/shared"
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
				Description: "The type of the socket. Valid values: `ssh`, `http`, `database`, `tls`, `vnc`, `rdp`, `vpn`, `subnet_router`, `exit_node`.",
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

			"http_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     service.HttpServiceTypeStandard,
							Description: "The upstream service type. Valid values: `standard`, `connector_file_server`. Defaults to `standard`.",
						},
						"upstream_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream HTTP URL. Format: `http(s)://<hostname>:<port>`. Example: `https://example.com` or `http://another.example.com:8080`. Only used when service type is `standard`.",
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream host header. Only used when service type is `standard`, and it's different from the hostname in `upstream_url`.",
						},
						"file_server_directory": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream file server directory. Only used when service type is `connector_file_server`.",
						},
					},
				},
			},

			"ssh_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     service.SshServiceTypeStandard,
							Description: "The upstream service type. Valid values: `standard`, `aws_ec2_instance_connect`, `aws_ssm`. Defaults to `standard`.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream SSH hostname.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The upstream SSH port number.",
						},
						"authentication_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The upstream authentication type for standard SSH service. Valid values: `username_and_password`, `border0_certificate`, `ssh_private_key`. Defaults to `border0_certificate`.",
						},
						"username_provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The upstream username provider. Valid values: `defined`, `prompt_client`, `use_connector_user`. Defaults to `prompt_client`.",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream username.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream password. Only used when authentication type is `username_and_password`.",
						},
						"private_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream private key. Only used when authentication type is `private_key`.",
						},
						"aws_credentials": shared.AwsCredentialsSchema,
						"ec2_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream EC2 instance id. Used when service type is either `aws_ec2_instance_connect` or `aws_ssm`.",
						},
						"ec2_instance_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream EC2 instance region. Used when service type is either `aws_ec2_instance_connect` or `aws_ssm` (SSM target type is `ec2`).",
						},
						"ssm_target_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The upstream SSM target type. Valid values: `ec2`, `ecs`. Defaults to `ec2`. Only used when service type is `aws_ssm`.",
						},
						"ecs_cluster_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream ECS cluster region. Only used when service type is `aws_ssm`, and SSM target type is `ecs`.",
						},
						"ecs_cluster_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream ECS cluster name. Only used when service type is `aws_ssm`, and SSM target type is `ecs`.",
						},
						"ecs_service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream ECS service name. Only used when service type is `aws_ssm`, and SSM target type is `ecs`.",
						},
					},
				},
			},

			"database_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     service.DatabaseServiceTypeStandard,
							Description: "The upstream service type. Valid values: `standard`, `aws_rds`, `gcp_cloud_sql`, `azure_sql`. Defaults to `standard`.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     service.DatabaseProtocolMySql,
							Description: "The upstream database protocol. Valid values: `mysql`, `postgres`. Defaults to `mysql`.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream database hostname.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The upstream database port number.",
						},
						"authentication_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     service.DatabaseAuthenticationTypeUsernameAndPassword,
							Description: "The upstream authentication type. Valid values: `username_and_password`, `tls`, `iam`. Defaults to `username_and_password`.",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream username. Used when authentication type is either `username_and_password` or `tls`.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream password. Used when authentication type is either `username_and_password` or `tls`.",
						},
						"certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream certificate. Only used when authentication type is `tls`.",
						},
						"private_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream private key. Only used when authentication type is `tls`.",
						},
						"ca_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream CA certificate.",
						},
						"aws_credentials": shared.AwsCredentialsSchema,
						"rds_instance_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream RDS database region. Only used when service type is `aws_rds`, and authentication type is `iam`.",
						},
						// gcp_cloud_sql only
						"cloudsql_connector_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if CloudSQL connector is enabled. Only used when service type is `gcp_cloud_sql`.",
						},
						"tls_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if TLS authentication is enabled. Only used when service type is `gcp_cloud_sql`.",
						},
						"cloudsql_iam_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if GCP IAM authentication is enabled. Only used when service type is `gcp_cloud_sql`.",
						},
						"gcp_credentials": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream GCP credentials.",
						},
						"cloudsql_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream CloudSQL instance id. Only used when service type is `gcp_cloud_sql`.",
						},
						// azure_sql only
						"azure_ad_integrated": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if Azure integrated authentication is enabled. Only used when service type is `azure_sql`.",
						},
						"azure_ad_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if Azure AD authentication is enabled. Only used when service type is `azure_sql`.",
						},
						"kerberos_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if Kerberos authentication is enabled. Only used when service type is `azure_sql`.",
						},
						"sql_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if standard SQL authentication (username and password) is enabled. Only used when service type is `azure_sql`.",
						},
					},
				},
			},

			"tls_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     service.TlsServiceTypeStandard,
							Description: "The upstream service type. Valid values: `standard`. Defaults to `standard`.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream TLS hostname.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The upstream TLS port number.",
						},
					},
				},
			},

			"vnc_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream VNC hostname.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The upstream VNC port number.",
						},
					},
				},
			},

			"rdp_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream RDP hostname.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The upstream RDP port number.",
						},
					},
				},
			},

			"vpn_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dhcp_pool_subnet": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The subnet in IPv4 CIDR notation to use for VPN client IP allocations (DHCP pool)",
						},
						"advertised_routes": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Set of routes to advertise to VPN clients in IPv4 CIDR notation",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			// DEPRECATED: will be removed after February 2025
			"subnet_routes_configuration": {
				Type:        schema.TypeList,
				Description: "THIS FIELD IS DEPRECATED, USE subnet_router_configuration",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_cidr_ranges": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Set of IPv4 routes to advertise to VPN clients in CIDR notation",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ipv6_cidr_ranges": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Set of IPv6 routes to advertise to VPN clients in CIDR notation",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"subnet_router_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4_cidr_ranges": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Set of IPv4 routes to advertise to VPN clients in CIDR notation",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ipv6_cidr_ranges": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Set of IPv6 routes to advertise to VPN clients in CIDR notation",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"exit_node_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// NOTE(@adrianosela): currently exit_node_configuration and the config object itself have
						// no attributes, so there is nothing to do here. If that ever changes, follow the pattern
						// used for subnet_router above.
					},
				},
			},
		},
	}
}

func resourceSocketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	// first let's try to get the socket by id
	// this api call may return socket that's soft deleted
	// we will try to fetch the same socket again by its name to make sure it's not deleted
	socket, diags := fetchSocket(ctx, d, m, d.Id())
	if diags.HasError() {
		return diags
	}
	if socket == nil {
		return nil
	}

	// and then get the socket by its name, in case if the socket is deleted
	// api only returns 404 error when the socket is fetched by name
	socket, diags = fetchSocket(ctx, d, m, socket.Name)
	if diags.HasError() {
		return diags
	}
	if socket == nil {
		return nil
	}

	// get socket linked connectors and their ids
	connectors, err := client.SocketConnectors(ctx, d.Id())
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch socket connectors")
	}
	// get socket upstream configs for linked connectors
	upstreamConfigs, err := client.SocketUpstreamConfigs(ctx, d.Id())
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch socket upstream configs")
	}

	// now inject the socket details, connector id and upstream config into the resource data
	if diags := schemautil.FromSocket(d, socket); diags.HasError() {
		return diags
	}
	if diags := schemautil.FromConnector(d, connectors); diags.HasError() {
		return diags
	}
	return schemautil.FromUpstreamConfig(d, socket, upstreamConfigs)
}

func fetchSocket(ctx context.Context, d *schema.ResourceData, m interface{}, idOrName string) (*border0client.Socket, diag.Diagnostics) {
	client := m.(border0client.Requester)

	// get socket by id or by name
	// when getting socket by id, api returns the socket even if it's soft deleted
	// when getting socket by name, api returns 404 error when the socket is deleted
	socket, err := client.Socket(ctx, idOrName)
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the socket was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] Socket (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil, nil
	}
	if err != nil {
		return nil, diagnostics.Error(err, "Failed to fetch socket")
	}

	return socket, nil
}

func resourceSocketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	socket := &border0client.Socket{
		Name:       d.Get("name").(string),
		SocketType: d.Get("socket_type").(string),
	}

	if diags := schemautil.ToSocket(d, socket); diags.HasError() {
		return diags
	}
	if diags := schemautil.ToUpstreamConfig(d, socket); diags.HasError() {
		return diags
	}

	created, err := client.CreateSocket(ctx, socket)
	if err != nil {
		return diagnostics.Error(err, "Failed to create socket")
	}

	d.SetId(created.SocketID)

	return resourceSocketRead(ctx, d, m)
}

func resourceSocketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	if d.HasChangesExcept("socket_type") {
		existingSocket, err := client.Socket(ctx, d.Id())
		if err != nil {
			return diagnostics.Error(err, "Failed to fetch socket")
		}
		socketUpdate := &border0client.Socket{
			Name:         d.Get("name").(string),
			SocketType:   d.Get("socket_type").(string),
			UpstreamType: existingSocket.UpstreamType,
		}

		if diags := schemautil.ToSocket(d, socketUpdate); diags.HasError() {
			return diags
		}
		if diags := schemautil.ToUpstreamConfig(d, socketUpdate); diags.HasError() {
			return diags
		}

		_, err = client.UpdateSocket(ctx, d.Id(), socketUpdate)
		if err != nil {
			return diagnostics.Error(err, "Failed to update socket")
		}
	}

	return resourceSocketRead(ctx, d, m)
}

func resourceSocketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	if err := client.DeleteSocket(ctx, d.Id()); err != nil {
		return diagnostics.Error(err, "Failed to delete socket")
	}
	d.SetId("")
	return nil
}
