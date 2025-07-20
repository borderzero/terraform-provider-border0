package border0_test

import (
	"fmt"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_DataSource_GroupsToIDs(t *testing.T) {
	mockGroupsResponse := border0client.Groups{
		List: []border0client.Group{
			{
				ID:          "mock1-id",
				DisplayName: "mock1",
			},
			{
				ID:          "mock2-id",
				DisplayName: "mock2",
			},
			{
				ID:          "mock3-id",
				DisplayName: "mock3",
			},
			{
				ID:          "mock4-id",
				DisplayName: "mock4",
			},
			{
				ID:          "mock5-id",
				DisplayName: "mock5",
			},
		},
	}

	config := fmt.Sprintf(`
		data "border0_group_names_to_ids" "unit_test" {
			names = [ "%s", "%s", "%s" ]
		}`,
		mockGroupsResponse.List[0].DisplayName,
		mockGroupsResponse.List[2].DisplayName,
		mockGroupsResponse.List[4].DisplayName,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// refresh for startup
		clientMock.EXPECT().Groups(matchContext).Return(&mockGroupsResponse, nil).Call,

		// refresh for apply, apply, and post-apply
		clientMock.EXPECT().Groups(matchContext).Return(&mockGroupsResponse, nil).Call,
		clientMock.EXPECT().Groups(matchContext).Return(&mockGroupsResponse, nil).Call,
		clientMock.EXPECT().Groups(matchContext).Return(&mockGroupsResponse, nil).Call,

		// refresh for cleanup
		clientMock.EXPECT().Groups(matchContext).Return(&mockGroupsResponse, nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.border0_group_names_to_ids.unit_test", "names.#", "3"),
					resource.TestCheckResourceAttr("data.border0_group_names_to_ids.unit_test", "names.0", mockGroupsResponse.List[0].DisplayName),
					resource.TestCheckResourceAttr("data.border0_group_names_to_ids.unit_test", "names.1", mockGroupsResponse.List[2].DisplayName),
					resource.TestCheckResourceAttr("data.border0_group_names_to_ids.unit_test", "names.2", mockGroupsResponse.List[4].DisplayName),
					resource.TestCheckResourceAttr("data.border0_group_names_to_ids.unit_test", "ids.#", "3"),
					resource.TestCheckResourceAttr("data.border0_group_names_to_ids.unit_test", "ids.0", mockGroupsResponse.List[0].ID),
					resource.TestCheckResourceAttr("data.border0_group_names_to_ids.unit_test", "ids.1", mockGroupsResponse.List[2].ID),
					resource.TestCheckResourceAttr("data.border0_group_names_to_ids.unit_test", "ids.2", mockGroupsResponse.List[4].ID),
				),
			},
		},
	})
}
