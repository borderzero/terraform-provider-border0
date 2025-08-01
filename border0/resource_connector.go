package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConnector() *schema.Resource {
	return &schema.Resource{
		Description:   "The connector resource allows you to create and manage a Border0 connector.",
		ReadContext:   resourceConnectorRead,
		CreateContext: resourceConnectorCreate,
		UpdateContext: resourceConnectorUpdate,
		DeleteContext: resourceConnectorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the connector. Connector name must contain only lowercase letters, numbers and dashes.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the connector.",
			},
			"built_in_ssh_service_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to expose the connector as an ssh service.",
			},
			"built_in_ssh_service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The socket id of the built-in ssh service.",
			},
		},
	}
}

func resourceConnectorRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)
	connector, err := client.Connector(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the connector was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] Connector (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch connector")
	}

	var builtInSshServiceID string
	if connector.BuiltInSshServiceEnabled && connector.BuiltInSshService != nil {
		builtInSshServiceID = connector.BuiltInSshService.SocketID
	}

	return schemautil.SetValues(d, map[string]any{
		"name":                         connector.Name,
		"description":                  connector.Description,
		"built_in_ssh_service_enabled": connector.BuiltInSshServiceEnabled,
		"built_in_ssh_service_id":      builtInSshServiceID,
	})
}

func resourceConnectorCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)
	connector := &border0client.Connector{
		Name: d.Get("name").(string),
		// we need to create the connector first and then enable the built-in ssh service
		// otherwise we may end up creating an orphaned connector if the built-in ssh service creation fails
		BuiltInSshServiceEnabled: false,
	}

	if v, ok := d.GetOk("description"); ok {
		connector.Description = v.(string)
	}

	var enableBuiltInSSHService bool
	if v, ok := d.GetOk("built_in_ssh_service_enabled"); ok {
		enableBuiltInSSHService = v.(bool)
	}

	// first create the connector
	// the built-in ssh service will not be created at this step to avoid orphaned connectors
	created, err := client.CreateConnector(ctx, connector)
	if err != nil {
		return diagnostics.Error(err, "Failed to create connector")
	}
	d.SetId(created.ConnectorID)
	if diags := resourceConnectorRead(ctx, d, m); diags.HasError() {
		return diags
	}

	// and then update the connector to create the built-in ssh service, if enabled
	if enableBuiltInSSHService {
		created.BuiltInSshServiceEnabled = true
		_, err := client.UpdateConnector(ctx, created)
		if err != nil {
			return diagnostics.Error(err, "Failed to enable built-in ssh service")
		}
		return resourceConnectorRead(ctx, d, m)
	}

	return nil
}

func resourceConnectorUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)

	fieldsToCheckForChanges := []string{
		"name",
		"description",
		"built_in_ssh_service_enabled",
	}

	if d.HasChanges(fieldsToCheckForChanges...) {
		connectorUpdate := &border0client.Connector{
			ConnectorID: d.Id(),
			Name:        d.Get("name").(string),
		}

		if v, ok := d.GetOk("description"); ok {
			connectorUpdate.Description = v.(string)
		}

		if v, ok := d.GetOk("built_in_ssh_service_enabled"); ok {
			connectorUpdate.BuiltInSshServiceEnabled = v.(bool)
		}

		_, err := client.UpdateConnector(ctx, connectorUpdate)
		if err != nil {
			return diagnostics.Error(err, "Failed to update connector")
		}
	}

	return resourceConnectorRead(ctx, d, m)
}

func resourceConnectorDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)
	if err := client.DeleteConnector(ctx, d.Id()); err != nil {
		return diagnostics.Error(err, "Failed to delete connector")
	}
	d.SetId("")
	return nil
}
