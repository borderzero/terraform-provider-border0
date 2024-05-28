package border0_test

import (
	"fmt"
	"testing"

	"github.com/autarch/testify/mock"
	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Resource_Border0User(t *testing.T) {
	notifyByEmail := false

	initialInput := border0client.User{
		DisplayName: "Unit Test Before",
		Email:       "user@unit-test.com",
		Role:        "admin",
	}
	initialOutput := border0client.User{
		ID:          "uuiduuid-uuid-uuid-uuid-uuiduuiduuid",
		DisplayName: initialInput.DisplayName,
		Email:       initialInput.Email,
		Role:        initialInput.Role,
	}
	updateInputOutput := border0client.User{
		ID:          initialOutput.ID,
		DisplayName: "Unit Test After",
		Email:       initialInput.Email,
		Role:        initialInput.Role,
	}

	userResourceFmt := `
	resource "border0_user" "unit_test" {
		display_name    = "%s"
		email           = "%s"
		role            = "%s"
		notify_by_email = %t
	}`

	initialUserConfig := fmt.Sprintf(
		userResourceFmt,
		initialInput.DisplayName,
		initialInput.Email,
		initialInput.Role,
		notifyByEmail,
	)

	updateUserConfig := fmt.Sprintf(
		userResourceFmt,
		updateInputOutput.DisplayName,
		initialInput.Email,
		initialInput.Role,
		notifyByEmail,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateUser(matchContext, &initialInput, mock.Anything).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().User(matchContext, initialOutput.ID).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().User(matchContext, initialOutput.ID).Return(&initialOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().User(matchContext, initialOutput.ID).Return(&initialOutput, nil).Call,

		// terraform apply (update + read + read)
		clientMock.EXPECT().UpdateUser(matchContext, &updateInputOutput).Return(&updateInputOutput, nil).Call,
		clientMock.EXPECT().User(matchContext, updateInputOutput.ID).Return(&updateInputOutput, nil).Call,
		clientMock.EXPECT().User(matchContext, updateInputOutput.ID).Return(&updateInputOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().User(matchContext, updateInputOutput.ID).Return(&updateInputOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteUser(matchContext, updateInputOutput.ID).Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: initialUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_user.unit_test", "display_name", initialOutput.DisplayName),
					resource.TestCheckResourceAttr("border0_user.unit_test", "email", initialOutput.Email),
					resource.TestCheckResourceAttr("border0_user.unit_test", "role", initialOutput.Role),
					resource.TestCheckResourceAttr("border0_user.unit_test", "notify_by_email", fmt.Sprintf("%t", notifyByEmail)),
					resource.TestCheckResourceAttrSet("border0_user.unit_test", "id"),
				),
			},
			{
				Config: updateUserConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_user.unit_test", "display_name", updateInputOutput.DisplayName),
					resource.TestCheckResourceAttr("border0_user.unit_test", "email", updateInputOutput.Email),
					resource.TestCheckResourceAttr("border0_user.unit_test", "role", updateInputOutput.Role),
					resource.TestCheckResourceAttr("border0_user.unit_test", "notify_by_email", fmt.Sprintf("%t", notifyByEmail)),
					resource.TestCheckResourceAttrSet("border0_user.unit_test", "id"),
				),
			},
			{
				ResourceName:      "border0_user.unit_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
