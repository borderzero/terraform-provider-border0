package border0

import (
	"context"
	"encoding/json"
	"hash/crc32"
	"sort"
	"strconv"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicyDocument() *schema.Resource {
	return &schema.Resource{
		Description: "`border0_policy_document` data source can be used to generate a policy document in JSON format for use with `border0_policy` resource.",
		ReadContext: dataSourcePolicyDocumentRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy language version.",
			},
			"action": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The actions that you want to allow.",
			},
			"condition": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"who": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The email address of the user who is allowed to perform the actions.",
									},
									"domain": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The domain of the user who is allowed to perform the actions.",
									},
								},
							},
							Description: "Who is allowed to perform the actions.",
						},
						"where": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_ip": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The IP address that the request must originate from to be allowed to perform the actions.",
									},
									"country": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The country that the request must originate from to be allowed to perform the actions.",
									},
									"country_not": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The country that the request must _NOT_ originate from to be allowed to perform the actions.",
									},
								},
							},
							Description: "Where the request must originate from to be allowed to perform the actions.",
						},
						"when": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"after": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the request must be made after to be allowed to perform the actions.",
									},
									"before": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the request must be made before to be allowed to perform the actions.",
									},
									"time_of_day_after": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the request must be made after to be allowed to perform the actions.",
									},
									"time_of_day_before": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the request must be made before to be allowed to perform the actions.",
									},
								},
							},
							Description: "When the request must be made to be allowed to perform the actions.",
						},
					},
				},
				Description: "The conditions under which you want to allow the actions.",
			},
		},
	}
}

func dataSourcePolicyDocumentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var policyData border0client.PolicyData

	if v, ok := d.GetOk("version"); ok {
		policyData.Version = v.(string)
	}
	if v, ok := d.GetOk("action"); ok {
		policyData.Action = policyDecodeStringList(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("condition"); ok {
		if conditions := v.(*schema.Set).List(); len(conditions) > 0 {
			condition := conditions[0].(map[string]interface{})
			if v, ok := condition["who"]; ok {
				if whos := v.(*schema.Set).List(); len(whos) > 0 {
					who := whos[0].(map[string]interface{})
					if v, ok := who["email"]; ok {
						policyData.Condition.Who.Email = policyDecodeStringList(v.(*schema.Set).List())
					}
					if v, ok := who["domain"]; ok {
						policyData.Condition.Who.Domain = policyDecodeStringList(v.(*schema.Set).List())
					}
				}
			}
			if v, ok := condition["where"]; ok {
				if wheres := v.(*schema.Set).List(); len(wheres) > 0 {
					where := wheres[0].(map[string]interface{})
					if v, ok := where["allowed_ip"]; ok {
						policyData.Condition.Where.AllowedIP = policyDecodeStringList(v.(*schema.Set).List())
					}
					if v, ok := where["country"]; ok {
						policyData.Condition.Where.Country = policyDecodeStringList(v.(*schema.Set).List())
					}
					if v, ok := where["country_not"]; ok {
						policyData.Condition.Where.CountryNot = policyDecodeStringList(v.(*schema.Set).List())
					}
				}
			}
			if v, ok := condition["when"]; ok {
				if whens := v.(*schema.Set).List(); len(whens) > 0 {
					when := whens[0].(map[string]interface{})
					if v, ok := when["after"]; ok {
						policyData.Condition.When.After = v.(string)
					}
					if v, ok := when["before"]; ok {
						policyData.Condition.When.Before = v.(string)
					}
					if v, ok := when["time_of_day_after"]; ok {
						policyData.Condition.When.TimeOfDayAfter = v.(string)
					}
					if v, ok := when["time_of_day_before"]; ok {
						policyData.Condition.When.TimeOfDayBefore = v.(string)
					}
				}
			}
		}
	}

	jsonPolicyData, err := json.MarshalIndent(policyData, "", "  ")
	if err != nil {
		return diagnostics.Error(err, "Failed to marshal policy data")
	}
	jsonString := string(jsonPolicyData)

	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(stringHashcode(jsonString)))
	return nil
}

func policyDecodeStringList(list []interface{}) []string {
	ret := make([]string, len(list))
	for i, value := range list {
		ret[i] = value.(string)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ret)))
	return ret
}

func stringHashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
