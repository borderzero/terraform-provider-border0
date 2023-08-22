package border0

import (
	"context"
	"fmt"
	"strings"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Attaches a managed policy to a socket.",
		ReadContext:   resourcePolicyAttachmentRead,
		CreateContext: resourcePolicyAttachmentCreate,
		DeleteContext: resourcePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the policy to attach.",
			},
			"socket_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the socket to attach the policy to.",
			},
		},
	}
}

func resourcePolicyAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	ids := strings.Split(d.Id(), ":")
	if len(ids) != 2 {
		return diag.Errorf("Invalid ID format: %s. Correct format is policyID:socketID", d.Id())
	}

	policyID, socketID := ids[0], ids[1]

	policy, err := client.Policy(ctx, policyID)
	if err != nil {
		return schemautil.DiagnosticsError(err, "Failed to fetch policy")
	}

	for _, eachSocketID := range policy.SocketIDs {
		if eachSocketID == socketID {
			return schemautil.SetValues(d, map[string]any{
				"policy_id": policyID,
				"socket_id": socketID,
			})
		}
	}
	return nil
}

func resourcePolicyAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	policyID := d.Get("policy_id").(string)
	socketID := d.Get("socket_id").(string)
	err := client.AttachPolicyToSocket(ctx, policyID, socketID)
	if err != nil {
		return schemautil.DiagnosticsError(err, "Failed to attach policy to socket")
	}
	d.SetId(fmt.Sprintf("%s:%s", policyID, socketID))
	return resourcePolicyAttachmentRead(ctx, d, m)
}

func resourcePolicyAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	ids := strings.Split(d.Id(), ":")
	if len(ids) != 2 {
		return diag.Errorf("Invalid ID format: %s. Correct format is policyID:socketID", d.Id())
	}
	policyID, socketID := ids[0], ids[1]
	if err := client.RemovePolicyFromSocket(ctx, policyID, socketID); err != nil {
		return schemautil.DiagnosticsError(err, "Failed to remove policy from socket")
	}
	d.SetId("")
	return nil
}
