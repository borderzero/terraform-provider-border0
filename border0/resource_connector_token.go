package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConnectorToken() *schema.Resource {
	return &schema.Resource{
		Description:   "The connector resource allows you to create and delete a token for a Border0 connector.",
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
				Description: "The name of the connector token.",
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

func resourceConnectorTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	connectorID := d.Get("connector_id").(string)
	connectorTokens, err := client.ConnectorTokens(ctx, connectorID)
	if !d.IsNewResource() && border0client.NotFound(err) {
		log.Printf("[WARN] No tokens found for connector (%s), removing from state", connectorID)
		d.SetId("")
		return diag.Diagnostics{}
	}
	if err != nil {
		return diagnosticsError(err, "Failed to fetch connector tokens")
	}

	var connectorToken *border0client.ConnectorToken
	for _, token := range connectorTokens.List {
		if token.ID == d.Id() {
			connectorToken = &token
			break
		}
	}

	if connectorToken == nil {
		log.Printf("[WARN] Token (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diag.Diagnostics{}
	}

	return setValues(d, map[string]any{
		"name":       connectorToken.Name,
		"expires_at": connectorToken.ExpiresAt.String(),
	})
}

func resourceConnectorTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	connectorID := d.Get("connector_id").(string)
	connectorToken := &border0client.ConnectorToken{
		ConnectorID: connectorID,
		Name:        d.Get("name").(string),
	}

	if v, ok := d.GetOk("expires_at"); ok {
		expiresAt, err := border0client.FlexibleTimeFrom(v.(string))
		if err != nil {
			return diagnosticsError(err, "Failed to parse expires_at")
		}
		connectorToken.ExpiresAt = expiresAt
	}

	created, err := client.CreateConnectorToken(ctx, connectorToken)
	if err != nil {
		return diagnosticsError(err, "Failed to create connector token")
	}

	d.SetId(created.ID)

	diagnotics := setValues(d, map[string]any{
		"connector_id": connectorID,
		"token":        created.Token,
	})
	if diagnotics != nil && diagnotics.HasError() {
		return diagnotics
	}

	return resourceConnectorTokenRead(ctx, d, m)
}

func resourceConnectorTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*border0client.APIClient)
	connectorID := d.Get("connector_id").(string)
	if err := client.DeleteConnectorToken(ctx, connectorID, d.Id()); err != nil {
		return diagnosticsError(err, "Failed to delete connector token")
	}
	d.SetId("")
	return nil
}
