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

func resourceConnectorToken() *schema.Resource {
	return &schema.Resource{
		Description:   "The connector token resource allows you to create and delete a token for a Border0 connector.",
		ReadContext:   resourceConnectorTokenRead,
		CreateContext: resourceConnectorTokenCreate,
		DeleteContext: resourceConnectorTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"connector_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the connector.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the connector token. Connector token name must contain only lowercase letters, numbers and dashes.",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The expiration date and time of the token. Leave empty for no expiration.",
			},
			"token": {
				Type:        schema.TypeString,
				Description: "The generated connector token.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceConnectorTokenRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)

	var connectorID, connectorTokenID string
	diags := schemautil.LoadMultipartID(d, &connectorID, &connectorTokenID)
	if diags.HasError() {
		return diags
	}
	if connectorID == "" {
		connectorID = d.Get("connector_id").(string)
	}

	connectorToken, err := client.ConnectorToken(ctx, connectorID, connectorTokenID)
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the connector token was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] Token (%s) for connector (%s) not found, removing from state", connectorTokenID, connectorID)
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch connector tokens")
	}

	return schemautil.SetValues(d, map[string]any{
		"name":       connectorToken.Name,
		"expires_at": connectorToken.ExpiresAt.String(),
	})
}

func resourceConnectorTokenCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	helper := m.(*ProviderHelper)
	client := helper.Requester

	connectorID := d.Get("connector_id").(string)
	connectorToken := &border0client.ConnectorToken{
		ConnectorID: connectorID,
		Name:        d.Get("name").(string),
	}

	if v, ok := d.GetOk("expires_at"); ok {
		expiresAt, err := border0client.FlexibleTimeFrom(v.(string))
		if err != nil {
			return diagnostics.Error(err, "Failed to parse expires_at")
		}
		connectorToken.ExpiresAt = expiresAt
	}

	created, err := client.CreateConnectorToken(ctx, connectorToken)
	if err != nil {
		return diagnostics.Error(err, "Failed to create connector token")
	}
	schemautil.SetMultipartID(d, connectorID, created.ID)

	if diags := schemautil.SetValues(d, map[string]any{
		"connector_id": connectorID,
		"token":        created.Token,
	}); diags.HasError() {
		return diags
	}

	helper.ReadAfterWriteDelay()
	return resourceConnectorTokenRead(ctx, d, m)
}

func resourceConnectorTokenDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)

	var connectorID, connectorTokenID string
	diags := schemautil.LoadMultipartID(d, &connectorID, &connectorTokenID)
	if diags.HasError() {
		return diags
	}

	if connectorID == "" {
		connectorID = d.Get("connector_id").(string)
	}

	if err := client.DeleteConnectorToken(ctx, connectorID, connectorTokenID); err != nil {
		return diagnostics.Error(err, "Failed to delete connector token")
	}
	d.SetId("")
	return nil
}
