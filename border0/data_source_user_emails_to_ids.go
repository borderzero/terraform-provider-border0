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

func dataSourceUserEmailsToIDs() *schema.Resource {
	return &schema.Resource{
		Description: "`border0_user_emails_to_ids` data source can be used to get ids for non-terraform-managed users (by email) for use with `border0_group` resource.",
		ReadContext: dataSourceUserEmailsToIDsRead,
		Schema: map[string]*schema.Schema{
			"emails": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of user emails to get IDs for",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ids": {
				Computed:    true,
				Type:        schema.TypeSet,
				Description: "Set of user IDs for the given user emails",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceUserEmailsToIDsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(border0client.Requester)

	emails := set.New[string]()
	if v := d.Get("emails"); v != nil {
		emails = set.New(schemaconvert.SetToSlice[string](v.(*schema.Set))...)
	}

	// NOTE: internally uses the border0 go sdk's users paginator
	// to retrieve all pages of users in the organization, using
	// the default page size defined there.
	users, err := client.Users(ctx)
	if err != nil {
		return diagnostics.Error(err, "Failed to fetch users list")
	}

	ids := []string{}
	for _, user := range users.List {
		if emails.Has(user.Email) {
			ids = append(ids, user.ID)
		}
	}
	d.Set("ids", ids)

	d.SetId(strconv.Itoa(stringHashcode(strings.Join(ids, ","))))
	return nil
}
