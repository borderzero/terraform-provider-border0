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

func resourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "The service account resource allows you to create and manage Border0 service accounts.",
		ReadContext:   resourceServiceAccountRead,
		CreateContext: resourceServiceAccountCreate,
		UpdateContext: resourceServiceAccountUpdate,
		DeleteContext: resourceServiceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name for the user in slug format. Must be unique per organization.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The role for the service account. Currently valid values include 'admin', 'member', 'read only', and 'client'.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the connector.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the service account should be active or not. Defaults to true",
			},
		},
	}
}

func resourceServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	serviceAccount, err := client.ServiceAccount(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the service account was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] Service Account (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch service account")
	}

	return schemautil.SetValues(d, map[string]any{
		"name":        serviceAccount.Name, // name is the primary identifier for service accounts
		"description": serviceAccount.Description,
		"role":        serviceAccount.Role,
		"active":      serviceAccount.Active,
	})
}

func resourceServiceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	serviceAccount := &border0client.ServiceAccount{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Role:        d.Get("role").(string),
		Active:      true, // default to true
	}
	if v := d.Get("active"); v != nil {
		serviceAccount.Active = v.(bool)
	}

	created, err := client.CreateServiceAccount(ctx, serviceAccount)
	if err != nil {
		return diagnostics.Error(err, "Failed to create service account")
	}
	d.SetId(created.Name)

	if diags := resourceServiceAccountRead(ctx, d, m); diags.HasError() {
		return diags
	}

	return nil
}

func resourceServiceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	fieldsToCheckForChanges := []string{
		"description",
		"role",
		"active",
	}

	if d.HasChanges(fieldsToCheckForChanges...) {
		serviceAccountUpdate := &border0client.ServiceAccount{
			Name:        d.Id(),
			Description: d.Get("description").(string),
			Role:        d.Get("role").(string),
			Active:      d.Get("active").(bool),
		}

		_, err := client.UpdateServiceAccount(ctx, serviceAccountUpdate)
		if err != nil {
			return diagnostics.Error(err, "Failed to update service account")
		}
	}

	return resourceServiceAccountRead(ctx, d, m)
}

func resourceServiceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	if err := client.DeleteServiceAccount(ctx, d.Get("name").(string)); err != nil {
		return diagnostics.Error(err, "Failed to delete service account")
	}
	d.SetId("")
	return nil
}
