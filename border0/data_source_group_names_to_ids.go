package border0

import (
	"context"
	"strconv"
	"strings"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/lib/types/set"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/borderzero/terraform-provider-border0/internal/schemautil/schemaconvert"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroupNamesToIDs() *schema.Resource {
	return &schema.Resource{
		Description: "`border0_group_names_to_ids` data source can be used to get ids for non-terraform-managed groups.",
		ReadContext: dataSourceGroupNamesToIDsRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of group names to get IDs for",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ids": {
				Computed:    true,
				Type:        schema.TypeSet,
				Description: "Set of group IDs for the given group names",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGroupNamesToIDsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)

	groups := set.New[string]()
	if v := d.Get("names"); v != nil {
		groups = set.New(schemaconvert.SetToSlice[string](v.(*schema.Set))...)
	}

	groupsList, err := client.Groups(ctx)
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch groups list")
	}

	ids := []string{}
	for _, group := range groupsList.List {
		if groups.Has(group.DisplayName) {
			ids = append(ids, group.ID)
		}
	}
	d.Set("ids", ids)

	d.SetId(strconv.Itoa(stringHashcode(strings.Join(ids, ","))))
	return nil
}
