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
							Description: "The upstream password.",
						},
						"private_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream private key.",
						},
						"aws_credentials": shared.AwsCredentialsSchema,
						"ec2_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream EC2 instance id.",
						},
						"ec2_instance_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream EC2 instance region.",
						},
						"ssm_target_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The upstream SSM target type. Valid values: `ec2`, `ecs`. Defaults to `ec2`.",
						},
						"ecs_cluster_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream ECS cluster region.",
						},
						"ecs_cluster_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream ECS cluster name.",
						},
						"ecs_service_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream ECS service name.",
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
							Description: "The upstream service type. Valid values: `standard`, `aws_rds`, `gcp_cloud_sql`. Defaults to `standard`.",
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
							Description: "The upstream username.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream password.",
						},
						"certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream certificate.",
						},
						"private_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The upstream private key.",
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
							Description: "The upstream RDS database region.",
						},
						"cloudsql_connector_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if CloudSQL connector is enabled. Only used when service type is `gcp_cloud_sql`.",
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
							Description: "The upstream CloudSQL instance id.",
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
							Description: "The upstream service type. Valid values: `standard`, `vpn`, http_proxy`. Defaults to `standard`.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The upstream TLS hostname. Only used when service type is `standard`.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The upstream TLS port number. Only used when service type is `standard`.",
						},
						"vpn_subnet": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The VPN subnet. Only used when service type is `vpn`.",
						},
						"vpn_routes": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The VPN routes. Only used when service type is `vpn`.",
						},
						"http_proxy_host_allowlist": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The HTTP proxy host allowlist. Only used when service type is `http_proxy`.",
						},
					},
				},
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
		return diagnostics.Error(err, "Failed to fetch socket")
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
