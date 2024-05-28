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

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "The user resource allows you to create and manage Border0 users.",
		ReadContext:   resourceUserRead,
		CreateContext: resourceUserCreate,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name for the user. A friendly name to help distinguish it among other users.",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email address of the user. User email must be unique per organization.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The role for the user. Currently valid values include 'admin', 'member', 'read only', and 'client'.",
			},
			"notify_by_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to notify the user that they have been added via email. Defaults to true",
			},
		},
	}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	user, err := client.User(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the user was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] User (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch user")
	}

	return schemautil.SetValues(d, map[string]any{
		"display_name":    user.DisplayName,
		"email":           user.Email,
		"role":            user.Role,
		"notify_by_email": d.Get("notify_by_email").(bool),
	})
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	user := &border0client.User{
		DisplayName: d.Get("display_name").(string),
		Email:       d.Get("email").(string),
		Role:        d.Get("role").(string),
	}

	notifyByEmail := true
	if v := d.Get("notify_by_email"); v != nil {
		notifyByEmail = v.(bool)
	}

	opts := []border0client.UserOption{}
	if !notifyByEmail {
		opts = append(opts, border0client.WithSkipNotification(!notifyByEmail))
	}

	created, err := client.CreateUser(ctx, user, opts...)
	if err != nil {
		return diagnostics.Error(err, "Failed to create user")
	}

	d.SetId(created.ID)

	if diags := resourceUserRead(ctx, d, m); diags.HasError() {
		return diags
	}

	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	fieldsToCheckForChanges := []string{
		"display_name",
		"email",
		"role",
	}

	if d.HasChanges(fieldsToCheckForChanges...) {
		userUpdate := &border0client.User{
			ID:          d.Id(),
			Email:       d.Get("email").(string),
			DisplayName: d.Get("display_name").(string),
			Role:        d.Get("role").(string),
		}

		_, err := client.UpdateUser(ctx, userUpdate)
		if err != nil {
			return diagnostics.Error(err, "Failed to update user")
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	if err := client.DeleteUser(ctx, d.Id()); err != nil {
		return diagnostics.Error(err, "Failed to delete user")
	}
	d.SetId("")
	return nil
}
