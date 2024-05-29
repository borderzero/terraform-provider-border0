package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/lib/types/slice"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/schemaconvert"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "The group resource allows you to create and manage Border0 groups.",
		ReadContext:   resourceGroupRead,
		CreateContext: resourceGroupCreate,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name for the user. A friendly name to help distinguish it among other users.",
			},
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of user ids (members of the group)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	group, err := client.Group(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the group was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch group")
	}

	return schemautil.SetValues(d, map[string]any{
		"display_name": group.DisplayName,
		"members":      slice.Transform(group.Members, func(u border0client.User) string { return u.ID }),
	})
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	group := &border0client.Group{
		DisplayName: d.Get("display_name").(string),
	}

	members := []string{}
	if v := d.Get("members"); v != nil {
		members = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}
	// client side validation for member IDs
	for _, member := range members {
		if _, err := uuid.ParseUUID(member); err != nil {
			return diagnostics.Error(err, "member id %s not a valid user uuid", member)
		}
	}

	created, err := client.CreateGroup(ctx, group)
	if err != nil {
		return diagnostics.Error(err, "Failed to create group")
	}
	d.SetId(created.ID)

	if _, err = client.UpdateGroupMemberships(ctx, created, members); err != nil {
		if delErr := client.DeleteGroup(ctx, created.ID); delErr != nil {
			return diagnostics.Error(err, "failed to create group memberships and failed to cleanup group afterwards: %v", delErr)
		}
		d.SetId("")
		return diagnostics.Error(err, "failed to create group memberships")
	}

	if diags := resourceGroupRead(ctx, d, m); diags.HasError() {
		return diags
	}

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	group := &border0client.Group{
		ID: d.Id(),
	}

	if d.HasChanges("members") {
		members := []string{}
		if v := d.Get("members"); v != nil {
			members = schemaconvert.SetToSlice[string](v.(*schema.Set))
		}
		for _, member := range members {
			if _, err := uuid.ParseUUID(member); err != nil {
				return diagnostics.Error(err, "member id %s not a valid user uuid", member)
			}
		}
		if _, err := client.UpdateGroupMemberships(ctx, group, members); err != nil {
			return diagnostics.Error(err, "Failed to update group memberships")
		}
	}

	if d.HasChanges("display_name") {
		group.DisplayName = d.Get("display_name").(string)
		if _, err := client.UpdateGroup(ctx, group); err != nil {
			return diagnostics.Error(err, "Failed to update group")
		}
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	if err := client.DeleteGroup(ctx, d.Id()); err != nil {
		return diagnostics.Error(err, "Failed to delete group")
	}
	d.SetId("")
	return nil
}
