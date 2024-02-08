package border0_test

import (
	"testing"

	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var policyDocumentConfig = `
data "border0_policy_document" "unit_test" {
  version = "v1"
  action = [ "database", "ssh", "http", "tls" ]
  condition {
    who {
      email = [ "johndoe@example.com" ]
      group = [ "db5c2352-b689-4135-babc-e97a8893128b" ]
      domain = [ "example.com" ]
    }
    where {
      allowed_ip = [ "0.0.0.0/0", "::/0" ]
      country = [ "NL", "CA", "US", "BR", "FR" ]
      country_not = [ "BE" ]
    }
    when {
      after = "2022-10-13T05:12:27Z"
      time_of_day_after = "00:00 UTC"
      time_of_day_before = "23:59 UTC"
    }
  }
}
`

func Test_DataSource_PolicyDocument(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, new(mocks.APIClientRequester)),
		Steps: []resource.TestStep{
			{
				Config: policyDocumentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "version", "v1"),
					// action list item gets sorted
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "action.0", "database"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "action.1", "http"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "action.2", "ssh"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "action.3", "tls"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.who.0.email.0", "johndoe@example.com"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.who.0.group.0", "db5c2352-b689-4135-babc-e97a8893128b"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.who.0.domain.0", "example.com"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.where.0.allowed_ip.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.where.0.allowed_ip.1", "::/0"),
					// country list item gets sorted
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.where.0.country.0", "BR"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.where.0.country.1", "CA"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.where.0.country.2", "FR"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.where.0.country.3", "NL"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.where.0.country.4", "US"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.where.0.country_not.0", "BE"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.when.0.after", "2022-10-13T05:12:27Z"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.when.0.time_of_day_after", "00:00 UTC"),
					resource.TestCheckResourceAttr("data.border0_policy_document.unit_test", "condition.0.when.0.time_of_day_before", "23:59 UTC"),
					resource.TestCheckResourceAttrSet("data.border0_policy_document.unit_test", "json"),
				),
			},
		},
	})
}
