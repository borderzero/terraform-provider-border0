package border0

import (
	"context"
	"log"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/schemaconvert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApprovalFlow() *schema.Resource {
	return &schema.Resource{
		Description:   "The approval flow resource allows you to create and manage a Border0 approval flow. An approval flow defines who can request just-in-time access to a set of sockets, and who can approve those requests.",
		ReadContext:   resourceApprovalFlowRead,
		CreateContext: resourceApprovalFlowCreate,
		UpdateContext: resourceApprovalFlowUpdate,
		DeleteContext: resourceApprovalFlowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the approval flow.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the approval flow.",
			},
			"socket_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Description:  "Set of socket IDs this approval flow applies to. At least one of `socket_ids` or `socket_tags` is required.",
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"socket_ids", "socket_tags"},
			},
			"socket_tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				Description:  "Map of socket tags this approval flow applies to. At least one of `socket_ids` or `socket_tags` is required.",
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"socket_ids", "socket_tags"},
			},
			"requester_user_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Description:  "Set of user IDs allowed to request access. At least one of `requester_user_ids` or `requester_group_ids` is required.",
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"requester_user_ids", "requester_group_ids"},
			},
			"requester_group_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Description:  "Set of group IDs allowed to request access. At least one of `requester_user_ids` or `requester_group_ids` is required.",
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"requester_user_ids", "requester_group_ids"},
			},
			"approver_user_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Description:  "Set of user IDs allowed to approve requests. At least one of `approver_user_ids` or `approver_group_ids` is required.",
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"approver_user_ids", "approver_group_ids"},
			},
			"approver_group_ids": {
				Type:         schema.TypeSet,
				Optional:     true,
				Description:  "Set of group IDs allowed to approve requests. At least one of `approver_user_ids` or `approver_group_ids` is required.",
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"approver_user_ids", "approver_group_ids"},
			},
			"allow_self_approval": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether a requester who is also an approver may approve their own request. Defaults to `false`.",
			},
		},
	}
}

func resourceApprovalFlowRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)

	flow, err := client.ApprovalWorkflow(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the approval flow was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] Approval flow (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch approval flow")
	}

	// Set socket_tags when present, and also clear it in state if it was previously set but is now empty.
	if _, ok := d.GetOk("socket_tags"); ok || len(flow.SocketTags) > 0 {
		if err := d.Set("socket_tags", flow.SocketTags); err != nil {
			return diagnostics.Error(err, "Failed to set socket_tags")
		}
	}

	return schemautil.SetValues(d, map[string]any{
		"name":                flow.Name,
		"description":         flow.Description,
		"socket_ids":          flow.SocketIDs,
		"requester_user_ids":  flow.RequesterUserIDs,
		"requester_group_ids": flow.RequesterGroupIDs,
		"approver_user_ids":   flow.ApproverUserIDs,
		"approver_group_ids":  flow.ApproverGroupIDs,
		"allow_self_approval": flow.AllowSelfApproval,
	})
}

func resourceApprovalFlowCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	helper := m.(*ProviderHelper)
	client := helper.Requester

	created, err := client.CreateApprovalWorkflow(ctx, approvalFlowFromResourceData(d))
	if err != nil {
		return diagnostics.Error(err, "Failed to create approval flow")
	}
	d.SetId(created.ID)

	helper.ReadAfterWriteDelay()
	return resourceApprovalFlowRead(ctx, d, m)
}

func resourceApprovalFlowUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	helper := m.(*ProviderHelper)
	client := helper.Requester

	flow := approvalFlowFromResourceData(d)
	flow.ID = d.Id()

	if _, err := client.UpdateApprovalWorkflow(ctx, flow); err != nil {
		return diagnostics.Error(err, "Failed to update approval flow")
	}

	helper.ReadAfterWriteDelay()
	return resourceApprovalFlowRead(ctx, d, m)
}

func resourceApprovalFlowDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)
	if err := client.DeleteApprovalWorkflow(ctx, d.Id()); err != nil {
		return diagnostics.Error(err, "Failed to delete approval flow")
	}
	d.SetId("")
	return nil
}

func approvalFlowFromResourceData(d *schema.ResourceData) *border0client.ApprovalWorkflow {
	flow := &border0client.ApprovalWorkflow{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		AllowSelfApproval: d.Get("allow_self_approval").(bool),
	}

	if v, ok := d.GetOk("socket_ids"); ok {
		flow.SocketIDs = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}
	if v, ok := d.GetOk("socket_tags"); ok {
		tags := make(map[string]string)
		for key, value := range v.(map[string]any) {
			tags[key] = value.(string)
		}
		flow.SocketTags = tags
	}
	if v, ok := d.GetOk("requester_user_ids"); ok {
		flow.RequesterUserIDs = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}
	if v, ok := d.GetOk("requester_group_ids"); ok {
		flow.RequesterGroupIDs = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}
	if v, ok := d.GetOk("approver_user_ids"); ok {
		flow.ApproverUserIDs = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}
	if v, ok := d.GetOk("approver_group_ids"); ok {
		flow.ApproverGroupIDs = schemaconvert.SetToSlice[string](v.(*schema.Set))
	}

	return flow
}
