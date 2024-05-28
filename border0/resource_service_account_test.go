package border0_test

import (
	"fmt"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Resource_Border0ServiceAccount(t *testing.T) {
	initialInput := border0client.ServiceAccount{
		Name:        "unit-test-service-account",
		Description: "Description Before",
		Role:        "admin",
		Active:      true,
	}
	initialOutput := border0client.ServiceAccount{
		Name:        initialInput.Name,
		Description: "Description Before",
		Role:        initialInput.Role,
		Active:      initialInput.Active,
	}
	updateInputOutput := border0client.ServiceAccount{
		Name:        initialInput.Name,
		Description: "Description After",
		Role:        initialInput.Role,
		Active:      !initialOutput.Active,
	}

	serviceAccountResourceFmt := `
	resource "border0_service_account" "unit_test" {
		name        = "%s"
		description = "%s"
		role        = "%s"
		active      = %t
	}`

	initialServiceAccountConfig := fmt.Sprintf(
		serviceAccountResourceFmt,
		initialInput.Name,
		initialInput.Description,
		initialInput.Role,
		initialInput.Active,
	)

	updateServiceAccountConfig := fmt.Sprintf(
		serviceAccountResourceFmt,
		updateInputOutput.Name,
		updateInputOutput.Description,
		updateInputOutput.Role,
		updateInputOutput.Active,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateServiceAccount(matchContext, &initialInput).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().ServiceAccount(matchContext, initialOutput.Name).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().ServiceAccount(matchContext, initialOutput.Name).Return(&initialOutput, nil).Call,

		// terraform import (read) - required for update
		clientMock.EXPECT().ServiceAccount(matchContext, initialOutput.Name).Return(&initialOutput, nil).Call,

		// terraform apply (update + read + read)
		clientMock.EXPECT().UpdateServiceAccount(matchContext, &updateInputOutput).Return(&updateInputOutput, nil).Call,
		clientMock.EXPECT().ServiceAccount(matchContext, updateInputOutput.Name).Return(&updateInputOutput, nil).Call,
		clientMock.EXPECT().ServiceAccount(matchContext, updateInputOutput.Name).Return(&updateInputOutput, nil).Call,

		// terraform import (read) - required for delete
		clientMock.EXPECT().ServiceAccount(matchContext, updateInputOutput.Name).Return(&updateInputOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteServiceAccount(matchContext, updateInputOutput.Name).Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: initialServiceAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_service_account.unit_test", "name", initialOutput.Name),
					resource.TestCheckResourceAttr("border0_service_account.unit_test", "description", initialOutput.Description),
					resource.TestCheckResourceAttr("border0_service_account.unit_test", "role", initialOutput.Role),
					resource.TestCheckResourceAttr("border0_service_account.unit_test", "active", fmt.Sprintf("%t", initialOutput.Active)),
					resource.TestCheckResourceAttrSet("border0_service_account.unit_test", "id"),
				),
			},
			{
				Config: updateServiceAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_service_account.unit_test", "name", updateInputOutput.Name),
					resource.TestCheckResourceAttr("border0_service_account.unit_test", "description", updateInputOutput.Description),
					resource.TestCheckResourceAttr("border0_service_account.unit_test", "role", updateInputOutput.Role),
					resource.TestCheckResourceAttr("border0_service_account.unit_test", "active", fmt.Sprintf("%t", updateInputOutput.Active)),
					resource.TestCheckResourceAttrSet("border0_service_account.unit_test", "id"),
				),
			},
			{
				ResourceName:      "border0_service_account.unit_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
