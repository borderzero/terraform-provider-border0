package border0

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strings"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "The policy resource allows you to create and manage a Border0 policy.",
		ReadContext:   resourcePolicyRead,
		CreateContext: resourcePolicyCreate,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy. Policy name must contain only lowercase letters, numbers and dashes.",
			},
			"policy_data": {
				Type:                  schema.TypeString,
				Required:              true,
				DiffSuppressFunc:      suppressEquivalentPolicyDiffs,
				DiffSuppressOnRefresh: true,
				Description:           "The policy data. This is a JSON string.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the policy.",
			},
			"org_wide": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the policy should be applied to all sockets in the organization.",
			},
		},
	}
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	policy, err := client.Policy(ctx, d.Id())
	if !d.IsNewResource() && border0client.NotFound(err) {
		// in case if the policy was deleted without Terraform knowing about it, we need to remove it from the state
		log.Printf("[WARN] Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch policy")
	}

	policyData, err := json.Marshal(&policy.PolicyData)
	if err != nil {
		return diagnostics.Error(err, "Failed to marshal policy data")
	}

	return schemautil.SetValues(d, map[string]any{
		"name":        policy.Name,
		"policy_data": string(policyData),
		"description": policy.Description,
		"org_wide":    policy.OrgWide,
	})
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	policy := &border0client.Policy{
		Name: d.Get("name").(string),
	}

	var policyData border0client.PolicyData
	if err := json.Unmarshal([]byte(d.Get("policy_data").(string)), &policyData); err != nil {
		return diagnostics.Error(err, "Failed to unmarshal policy data")
	}
	policy.PolicyData = policyData

	if v, ok := d.GetOk("description"); ok {
		policy.Description = v.(string)
	}
	if v, ok := d.GetOk("org_wide"); ok {
		policy.OrgWide = v.(bool)
	}

	created, err := client.CreatePolicy(ctx, policy)
	if err != nil {
		return diagnostics.Error(err, "Failed to create policy")
	}

	d.SetId(created.ID)

	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)

	if d.HasChangesExcept("org_wide") {
		policyUpdate := &border0client.Policy{
			Name: d.Get("name").(string),
		}

		var policyData border0client.PolicyData
		if err := json.Unmarshal([]byte(d.Get("policy_data").(string)), &policyData); err != nil {
			return diagnostics.Error(err, "Failed to unmarshal policy data")
		}
		policyUpdate.PolicyData = policyData

		if v, ok := d.GetOk("description"); ok {
			policyUpdate.Description = v.(string)
		}

		_, err := client.UpdatePolicy(ctx, d.Id(), policyUpdate)
		if err != nil {
			return diagnostics.Error(err, "Failed to update policy")
		}
	}

	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	if err := client.DeletePolicy(ctx, d.Id()); err != nil {
		return diagnostics.Error(err, "Failed to delete policy")
	}
	d.SetId("")
	return nil
}

func suppressEquivalentPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	if strings.TrimSpace(old) == "" && strings.TrimSpace(new) == "" {
		return true
	}

	if strings.TrimSpace(old) == "{}" && strings.TrimSpace(new) == "" {
		return true
	}

	if strings.TrimSpace(old) == "" && strings.TrimSpace(new) == "{}" {
		return true
	}

	if strings.TrimSpace(old) == "{}" && strings.TrimSpace(new) == "{}" {
		return true
	}

	var oldJSONAsPolicyData, newJSONAsPolicyData border0client.PolicyData

	if err := json.Unmarshal([]byte(old), &oldJSONAsPolicyData); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(new), &newJSONAsPolicyData); err != nil {
		return false
	}

	return reflect.DeepEqual(oldJSONAsPolicyData, newJSONAsPolicyData)
}
