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

func resourceServiceAccountToken() *schema.Resource {
	return &schema.Resource{
		Description:   "The service account token resource allows you to create and delete a token for a Border0 service account.",
		ReadContext:   resourceServiceAccountTokenRead,
		CreateContext: resourceServiceAccountTokenCreate,
		DeleteContext: resourceServiceAccountTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"service_account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service account.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the service account token. Service account token name must contain only lowercase letters, numbers and dashes.",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The expiration date and time of the token. Leave empty for no expiration.",
			},
			"token": {
				Type:        schema.TypeString,
				Description: "The generated service account token.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceServiceAccountTokenRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)

	var serviceAccountTokenID, serviceAccountName string
	diags := schemautil.LoadMultipartID(d, &serviceAccountTokenID, &serviceAccountName)
	if diags.HasError() {
		return diags
	}

	serviceAccountTokens, err := client.ServiceAccountTokens(ctx, serviceAccountName)
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the service account was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] Service account (%s) not found, removing service account token from state", serviceAccountName)
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch service account tokens")
	}

	for _, tk := range serviceAccountTokens.List {
		if tk.ID == serviceAccountTokenID {
			return schemautil.SetValues(d, map[string]any{
				"name":       tk.Name,
				"expires_at": tk.ExpiresAt.String(),
			})
		}
	}

	// in case if the connector token was deleted without Terraform knowing about it, we need to remove it from the state
	log.Printf("[WARN] Token (%s) for service account (%s) not found, removing from state", serviceAccountTokenID, serviceAccountName)
	d.SetId("")
	return nil
}

func resourceServiceAccountTokenCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	helper := m.(*ProviderHelper)
	client := helper.Requester

	serviceAccountName := d.Get("service_account_name").(string)
	tokenName := d.Get("name").(string)

	serviceAccountToken := &border0client.ServiceAccountToken{
		Name: tokenName,
	}

	if v, ok := d.GetOk("expires_at"); ok {
		expiresAt, err := border0client.FlexibleTimeFrom(v.(string))
		if err != nil {
			return diagnostics.Error(err, "Failed to parse expires_at")
		}
		serviceAccountToken.ExpiresAt = expiresAt
	}

	created, err := client.CreateServiceAccountToken(ctx, serviceAccountName, serviceAccountToken)
	if err != nil {
		return diagnostics.Error(err, "Failed to create service account token")
	}

	schemautil.SetMultipartID(d, created.ID, serviceAccountName)
	if diags := schemautil.SetValues(d, map[string]any{
		"service_account_name": serviceAccountName,
		"token":                created.Token,
	}); diags.HasError() {
		return diags
	}

	helper.ReadAfterWriteDelay()
	return resourceServiceAccountTokenRead(ctx, d, m)
}

func resourceServiceAccountTokenDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)

	var serviceAccountTokenID, serviceAccountName string
	diags := schemautil.LoadMultipartID(d, &serviceAccountTokenID, &serviceAccountName)
	if diags.HasError() {
		return diags
	}

	if err := client.DeleteServiceAccountToken(ctx, serviceAccountName, serviceAccountTokenID); err != nil {
		return diagnostics.Error(err, "Failed to delete service account token")
	}
	d.SetId("")
	return nil
}
