package border0

import (
	"context"
	"fmt"
	"log"
	"strings"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
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

func resourceConnectorTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	connectorID, connectorTokenID, diags := determineConnectorIDAndConnectorTokenID(d)
	if diags.HasError() {
		return diags
	}

	connectorTokens, err := client.ConnectorTokens(ctx, connectorID)
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the connector was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] No tokens found for connector (%s), removing from state", connectorID)
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch connector tokens")
	}

	var connectorToken *border0client.ConnectorToken
	for _, token := range connectorTokens.List {
		if token.ID == connectorTokenID {
			connectorToken = &token
			break
		}
	}

	if connectorToken == nil {
		if !d.IsNewResource() {
			// in case if the connector token was deleted without Terraform knowing about it, we need to remove it from the state
			log.Printf("[WARN] Token (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		// otherwise, we need to fail the read
		return diag.Errorf("Connector token not found (connector id: %s, connector token id: %s)", connectorID, connectorTokenID)
	}

	return schemautil.SetValues(d, map[string]any{
		"name":       connectorToken.Name,
		"expires_at": connectorToken.ExpiresAt.String(),
	})
}

func resourceConnectorTokenCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
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

	d.SetId(fmt.Sprintf("%s:%s", connectorID, created.ID))

	if diags := schemautil.SetValues(d, map[string]any{
		"connector_id": connectorID,
		"token":        created.Token,
	}); diags.HasError() {
		return diags
	}

	return resourceConnectorTokenRead(ctx, d, m)
}

func resourceConnectorTokenDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	connectorID, connectorTokenID, diags := determineConnectorIDAndConnectorTokenID(d)
	if diags.HasError() {
		return diags
	}
	if err := client.DeleteConnectorToken(ctx, connectorID, connectorTokenID); err != nil {
		return diagnostics.Error(err, "Failed to delete connector token")
	}
	d.SetId("")
	return nil
}

func determineConnectorIDAndConnectorTokenID(d *schema.ResourceData) (connectorID string, connectorTokenID string, diags diag.Diagnostics) {
	ids := strings.Split(d.Id(), ":")
	if len(ids) == 1 {
		connectorID = d.Get("connector_id").(string)
		connectorTokenID = ids[0]
	} else if len(ids) == 2 {
		connectorID = ids[0]
		connectorTokenID = ids[1]
	} else {
		diags = diag.Errorf("Invalid ID format: %s", d.Id())
	}
	return connectorID, connectorTokenID, diags
}
