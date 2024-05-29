package border0_test

import (
	"fmt"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_DataSource_UserEmailsToIDs(t *testing.T) {
	mockUsersResponse := border0client.Users{
		List: []border0client.User{
			{
				ID:    "mock1-id",
				Email: "mock1@gmail.com",
			},
			{
				ID:    "mock2-id",
				Email: "mock2@gmail.com",
			},
			{
				ID:    "mock3-id",
				Email: "mock3@gmail.com",
			},
			{
				ID:    "mock4-id",
				Email: "mock4@gmail.com",
			},
			{
				ID:    "mock5-id",
				Email: "mock5@gmail.com",
			},
		},
	}

	config := fmt.Sprintf(`
		data "border0_user_emails_to_ids" "unit_test" {
			emails = [ "%s", "%s", "%s" ]
		}`,
		mockUsersResponse.List[0].Email,
		mockUsersResponse.List[2].Email,
		mockUsersResponse.List[4].Email,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// refresh for startup
		clientMock.EXPECT().Users(matchContext).Return(&mockUsersResponse, nil).Call,

		// refresh for apply, apply, and post-apply
		clientMock.EXPECT().Users(matchContext).Return(&mockUsersResponse, nil).Call,
		clientMock.EXPECT().Users(matchContext).Return(&mockUsersResponse, nil).Call,
		clientMock.EXPECT().Users(matchContext).Return(&mockUsersResponse, nil).Call,

		// refresh for cleanup
		clientMock.EXPECT().Users(matchContext).Return(&mockUsersResponse, nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.border0_user_emails_to_ids.unit_test", "emails.#", "3"),
					resource.TestCheckResourceAttr("data.border0_user_emails_to_ids.unit_test", "emails.0", mockUsersResponse.List[0].Email),
					resource.TestCheckResourceAttr("data.border0_user_emails_to_ids.unit_test", "emails.1", mockUsersResponse.List[2].Email),
					resource.TestCheckResourceAttr("data.border0_user_emails_to_ids.unit_test", "emails.2", mockUsersResponse.List[4].Email),
					resource.TestCheckResourceAttr("data.border0_user_emails_to_ids.unit_test", "ids.#", "3"),
					resource.TestCheckResourceAttr("data.border0_user_emails_to_ids.unit_test", "ids.0", mockUsersResponse.List[0].ID),
					resource.TestCheckResourceAttr("data.border0_user_emails_to_ids.unit_test", "ids.1", mockUsersResponse.List[2].ID),
					resource.TestCheckResourceAttr("data.border0_user_emails_to_ids.unit_test", "ids.2", mockUsersResponse.List[4].ID),
				),
			},
		},
	})
}
