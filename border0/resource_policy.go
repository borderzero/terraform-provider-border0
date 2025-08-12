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
	"github.com/borderzero/terraform-provider-border0/internal/lib/sem"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePolicy(semaphore sem.Semaphore) *schema.Resource {
	return &schema.Resource{
		Description:   "The policy resource allows you to create and manage a Border0 policy.",
		ReadContext:   resourcePolicyRead,
		CreateContext: getResourcePolicyCreate(semaphore),
		UpdateContext: getResourcePolicyUpdate(semaphore),
		DeleteContext: getResourcePolicyDelete(semaphore),
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
			"tag_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of tag rules to apply to the sockets that this policy is applied to.",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
	pruneNullValues(pdIface)
	filteredPolicyData, err := json.Marshal(pdIface)
	if err != nil {
		return diagnostics.Error(err, "Failed to marshal filtered policy data")
	}

	tagRulesSlice := make([]map[string]any, 0, len(policy.TagRules))
	for _, rule := range policy.TagRules {
		ruleMap := make(map[string]any)
		for key, val := range rule {
			ruleMap[key] = val
		}
		tagRulesSlice = append(tagRulesSlice, ruleMap)
	}

	return schemautil.SetValues(d, map[string]any{
		"name":        policy.Name,
		"policy_data": string(filteredPolicyData),
		"description": policy.Description,
		"org_wide":    policy.OrgWide,
		"version":     policy.Version,
		"tag_rules":   tagRulesSlice,
	})
}

func getResourcePolicyCreate(sem sem.Semaphore) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		sem.Acquire()
		defer sem.Release()

		helper := m.(*ProviderHelper)
		client := helper.Requester

		policy := &border0client.Policy{
			Name:     d.Get("name").(string),
			Version:  d.Get("version").(string),
			TagRules: mapTagRules(d.Get("tag_rules")),
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

		helper.ReadAfterWriteDelay()
		return resourcePolicyRead(ctx, d, m)
	}
}

func getResourcePolicyUpdate(sem sem.Semaphore) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		sem.Acquire()
		defer sem.Release()

		helper := m.(*ProviderHelper)
		client := helper.Requester

		if d.HasChangesExcept("org_wide") {
			policyUpdate := &border0client.Policy{
				Name:     d.Get("name").(string),
				TagRules: mapTagRules(d.Get("tag_rules")),
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

		helper.ReadAfterWriteDelay()
		return resourcePolicyRead(ctx, d, m)
	}
}

func getResourcePolicyDelete(sem sem.Semaphore) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
		sem.Acquire()
		defer sem.Release()

		client := m.(border0client.Requester)
		if err := client.DeletePolicy(ctx, d.Id()); err != nil {
			return diagnostics.Error(err, "Failed to delete policy")
		}
		d.SetId("")
		return nil
	}
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
	old = pruneNullJSON(old)
	new = pruneNullJSON(new)
	return jsoneq.AreEqual(old, new)
}

func mapTagRules(tagRules any) []map[string]string {
	if raw, ok := tagRules.([]any); ok {
		if len(raw) == 0 {
			return []map[string]string{}
		}
		rules := make([]map[string]string, 0, len(raw))
		for _, rule := range raw {
			if ruleMap, ok := rule.(map[string]any); ok {
				stringMap := make(map[string]string)
				for key, val := range ruleMap {
					if strVal, ok := val.(string); ok {
						stringMap[key] = strVal
					} else {
						log.Printf("[WARN] Expected string value for key '%s', got %T", key, val)
					}
				}
				rules = append(rules, stringMap)
			}
		}
		return rules
	}
	return []map[string]string{}
}

func pruneNullValues(v any) {
	switch x := v.(type) {
	case map[string]any:
		for k, e := range x {
			switch val := e.(type) {
			case nil:
				delete(x, k)
				continue
			case string:
				if val == "" {
					delete(x, k)
					continue
				}
			case []any:
				if len(val) == 0 {
					delete(x, k)
					continue
				}
				pruneNullValues(val)
				continue
			case map[string]any:
				pruneNullValues(val)
				continue
			default:
				pruneNullValues(val)
			}
		}
	case []any:
		for _, e := range x {
			pruneNullValues(e)
		}
	}
}

func pruneNullJSON(s string) string {
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	pruneNullValues(v)
	b, err := json.Marshal(v)
	if err != nil {
		return s
	}
	return string(b)
}
