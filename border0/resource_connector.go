package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
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
				Description: "The name of the connector.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the connector.",
			},
		},
	}
}

func resourceConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	connector, err := client.Connector(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		log.Printf("[WARN] Connector (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diag.Diagnostics{}
	}
	if err != nil {
		return diagnosticsError(err, "Failed to fetch connector")
	}

	return setValues(d, map[string]any{
		"name":        connector.Name,
		"description": connector.Description,
	})
}

func resourceConnectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	connector := &border0client.Connector{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		connector.Description = v.(string)
	}

	created, err := client.CreateConnector(ctx, connector)
	if err != nil {
		return diagnosticsError(err, "Failed to create connector")
	}

	d.SetId(created.ConnectorID)

	return resourceConnectorRead(ctx, d, m)
}

func resourceConnectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)

	if d.HasChanges("name", "description") {
		connectorUpdate := &border0client.Connector{
			ConnectorID: d.Id(),
			Name:        d.Get("name").(string),
		}

		if v, ok := d.GetOk("description"); ok {
			connectorUpdate.Description = v.(string)
		}

		_, err := client.UpdateConnector(ctx, connectorUpdate)
		if err != nil {
			return diagnosticsError(err, "Failed to update connector")
		}
	}

	return resourceConnectorRead(ctx, d, m)
}

func resourceConnectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	if err := client.DeleteConnector(ctx, d.Id()); err != nil {
		return diagnosticsError(err, "Failed to delete connector")
	}
	d.SetId("")
	return nil
}
