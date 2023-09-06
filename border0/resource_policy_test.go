package border0_test

import (
	"encoding/json"
	"testing"

	"github.com/autarch/testify/require"
	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialPolicyConfig = `
resource "border0_policy" "unit_test" {
  name = "unit-test-policy-1"
  description = "policy created from terraform unit test"
  policy_data = jsonencode({
    "version": "v1",
    "action": [ "database", "ssh", "http", "tls" ],
    "condition": {
      "who": {
        "email": [ "johndoe@example.com" ],
        "domain": [ "example.com" ]
      },
      "where": {
        "allowed_ip": [ "0.0.0.0/0", "::/0" ],
        "country": [ "NL", "CA", "US", "BR", "FR" ],
        "country_not": [ "BE" ]
      },
      "when": {
        "after": "2022-10-13T05:12:26Z",
        "before": null,
        "time_of_day_after": "00:00 UTC",
        "time_of_day_before": "23:59 UTC"
      }
    }
  })
}
`

var updatePolicyConfig = `
resource "border0_policy" "unit_test" {
  name = "unit-test-policy-1"
  description = "update policy description"
  policy_data = jsonencode({
    "version": "v1",
    "action": [ "database", "ssh", "http", "tls" ],
    "condition": {
      "who": {
        "email": [ "johndoe@example.com", "another@example.com" ],
        "domain": [ "example.com" ]
      },
      "where": {
        "allowed_ip": [ "0.0.0.0/0", "::/0" ],
        "country": [ "NL", "CA", "US", "BR", "FR" ],
        "country_not": [ "BE" ]
      },
      "when": {
        "after": "2022-10-13T05:12:26Z",
        "before": null,
        "time_of_day_after": "00:00 UTC",
        "time_of_day_before": "23:59 UTC"
      }
    }
  })
}
`

func Test_Resource_Border0Policy(t *testing.T) {
	initialPolicyData := border0client.PolicyData{
		Version: "v1",
		Action:  []string{"database", "ssh", "http", "tls"},
		Condition: border0client.PolicyCondition{
			Who: border0client.PolicyWho{
				Email:  []string{"johndoe@example.com"},
				Domain: []string{"example.com"},
			},
			Where: border0client.PolicyWhere{
				AllowedIP:  []string{"0.0.0.0/0", "::/0"},
				Country:    []string{"NL", "CA", "US", "BR", "FR"},
				CountryNot: []string{"BE"},
			},
			When: border0client.PolicyWhen{
				After:           "2022-10-13T05:12:26Z",
				Before:          "",
				TimeOfDayAfter:  "00:00 UTC",
				TimeOfDayBefore: "23:59 UTC",
			},
		},
	}
	initialInput := border0client.Policy{
		Name:        "unit-test-policy-1",
		Description: "policy created from terraform unit test",
		PolicyData:  initialPolicyData,
	}
	initialOutput := border0client.Policy{
		ID:          "unit-test-id-1",
		Name:        "unit-test-policy-1",
		Description: "policy created from terraform unit test",
		PolicyData:  initialPolicyData,
	}

	updatePolicyData := border0client.PolicyData{
		Version: "v1",
		Action:  []string{"database", "ssh", "http", "tls"},
		Condition: border0client.PolicyCondition{
			Who: border0client.PolicyWho{
				Email: []string{
					"johndoe@example.com",
					"another@example.com",
				},
				Domain: []string{"example.com"},
			},
			Where: border0client.PolicyWhere{
				AllowedIP:  []string{"0.0.0.0/0", "::/0"},
				Country:    []string{"NL", "CA", "US", "BR", "FR"},
				CountryNot: []string{"BE"},
			},
			When: border0client.PolicyWhen{
				After:           "2022-10-13T05:12:26Z",
				Before:          "",
				TimeOfDayAfter:  "00:00 UTC",
				TimeOfDayBefore: "23:59 UTC",
			},
		},
	}
	updateInput := border0client.Policy{
		Name:        "unit-test-policy-1",
		Description: "update policy description",
		PolicyData:  updatePolicyData,
	}
	updateOutput := border0client.Policy{
		ID:          "unit-test-id-1",
		Name:        "unit-test-policy-1",
		Description: "update policy description",
		PolicyData:  updatePolicyData,
	}

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// read = client.Policy()
		// create = client.CreatePolicy()
		// update = client.UpdatePolicy()
		// delete = client.DeletePolicy()

		// terraform apply (create + read + read)
		clientMock.EXPECT().CreatePolicy(matchContext, &initialInput).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,

		// terraform aplly (update + read + read)
		clientMock.EXPECT().UpdatePolicy(matchContext, "unit-test-id-1", &updateInput).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&updateOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&updateOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeletePolicy(matchContext, "unit-test-id-1").Return(nil).Call,
	)

	initialPolicyDataJSON, err := json.Marshal(initialPolicyData)
	require.NoError(t, err)

	updatePolicyDataJSON, err := json.Marshal(updatePolicyData)
	require.NoError(t, err)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: initialPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_policy.unit_test", "name", "unit-test-policy-1"),
					resource.TestCheckResourceAttr("border0_policy.unit_test", "description", "policy created from terraform unit test"),
					resource.TestCheckResourceAttrSet("border0_policy.unit_test", "id"),
					testMatchResourceAttrJSON("border0_policy.unit_test", "policy_data", string(initialPolicyDataJSON)),
				),
			},
			{
				Config: updatePolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_policy.unit_test", "name", "unit-test-policy-1"),
					resource.TestCheckResourceAttr("border0_policy.unit_test", "description", "update policy description"),
					testMatchResourceAttrJSON("border0_policy.unit_test", "policy_data", string(updatePolicyDataJSON)),
				),
			},
			{
				ResourceName:      "border0_policy.unit_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
