package border0

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strings"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/lib/types/jsoneq"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "v2",
				Description:  "The version of the policy. The default value is 'v2', the other valid value is 'v1'.",
				ValidateFunc: validation.StringInSlice([]string{"v1", "v2"}, false),
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
			"socket_tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A set of tags to apply to the sockets that this policy is applied to.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

	rawPolicyData, err := json.Marshal(&policy.PolicyData)
	if err != nil {
		return diagnostics.Error(err, "Failed to marshal policy data")
	}
	var pdIface any
	if err := json.Unmarshal(rawPolicyData, &pdIface); err != nil {
		return diagnostics.Error(err, "Failed to process policy data")
	}
	jsoneq.Prune(pdIface, jsoneq.PruneEmptySlices(), jsoneq.PruneEmptyStrings())
	filteredPolicyData, err := json.Marshal(pdIface)
	if err != nil {
		return diagnostics.Error(err, "Failed to marshal filtered policy data")
	}

	return schemautil.SetValues(d, map[string]any{
		"name":        policy.Name,
		"policy_data": string(filteredPolicyData),
		"description": policy.Description,
		"org_wide":    policy.OrgWide,
		"version":     policy.Version,
		"socket_tags": func() map[string]any {
			m := make(map[string]interface{})
			for key, val := range policy.SocketTags {
				m[key] = val
			}
			return m
		}(),
	})
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(border0client.Requester)
	policy := &border0client.Policy{
		Name:       d.Get("name").(string),
		Version:    d.Get("version").(string),
		SocketTags: mapSocketTags(d.Get("socket_tags")),
	}

	switch policy.Version {
	case "v1":
		var policyData border0client.PolicyData
		if err := json.Unmarshal([]byte(d.Get("policy_data").(string)), &policyData); err != nil {
			return diagnostics.Error(err, "Failed to unmarshal policy data")
		}
		policy.PolicyData = policyData
	case "v2":
		var policyData border0client.PolicyDataV2
		if err := json.Unmarshal([]byte(d.Get("policy_data").(string)), &policyData); err != nil {
			return diagnostics.Error(err, "Failed to unmarshal policy data")
		}
		policy.PolicyData = policyData
	default:
		return diag.Errorf("Invalid policy version: %s", policy.Version)
	}

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
			Name:       d.Get("name").(string),
			SocketTags: mapSocketTags(d.Get("socket_tags")),
		}

		switch d.Get("version").(string) {
		case "v1":
			var policyData border0client.PolicyData
			if err := json.Unmarshal([]byte(d.Get("policy_data").(string)), &policyData); err != nil {
				return diagnostics.Error(err, "Failed to unmarshal policy data")
			}
			policyUpdate.Version = "v1"
			policyUpdate.PolicyData = policyData
		case "v2":
			var policyData border0client.PolicyDataV2
			if err := json.Unmarshal([]byte(d.Get("policy_data").(string)), &policyData); err != nil {
				return diagnostics.Error(err, "Failed to unmarshal policy data")
			}
			policyUpdate.Version = "v2"
			policyUpdate.PolicyData = policyData
		default:
			return diag.Errorf("Invalid policy version: %s", policyUpdate.Version)
		}

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

// suppressEquivalentPolicyDiffs suppresses spurious diffs in policy_data by checking
// for semantic equivalence, ignoring default values and unordered arrays.
func suppressEquivalentPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldTrim := strings.TrimSpace(old)
	newTrim := strings.TrimSpace(new)
	if oldTrim == "" && newTrim == "" {
		return true
	}
	if (oldTrim == "{}" && newTrim == "") || (oldTrim == "" && newTrim == "{}") || (oldTrim == "{}" && newTrim == "{}") {
		return true
	}

	// first, unmarshal into typed structs to ignore default values
	switch d.Get("version").(string) {
	case "v1":
		var oldPD, newPD border0client.PolicyData
		if err := json.Unmarshal([]byte(old), &oldPD); err == nil {
			if err := json.Unmarshal([]byte(new), &newPD); err == nil {
				if reflect.DeepEqual(oldPD, newPD) {
					return true
				}
			}
		}
	case "v2":
		var oldPD, newPD border0client.PolicyDataV2
		if err := json.Unmarshal([]byte(old), &oldPD); err == nil {
			if err := json.Unmarshal([]byte(new), &newPD); err == nil {
				if reflect.DeepEqual(oldPD, newPD) {
					return true
				}
			}
		}
	default:
		return false
	}

	// fallback to generic unordered JSON comparison for arrays
	// NOTE: we MUST keep empty objects.
	return jsoneq.AreEqual(old, new, jsoneq.PruneEmptySlices(), jsoneq.PruneEmptyStrings())
}

func mapSocketTags(socketTags any) map[string]string {
	if raw, ok := socketTags.(map[string]interface{}); ok {
		if len(raw) == 0 {
			return nil
		}
		m := make(map[string]string)
		for key, val := range raw {
			m[key] = val.(string)
		}
		return m
	}
	return nil
}
