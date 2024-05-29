package border0_test

import (
	"fmt"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func Test_Resource_Border0Group(t *testing.T) {
	memberIDA, err := uuid.GenerateUUID()
	assert.NoError(t, err)

	groupCreateInput := border0client.Group{
		DisplayName: "Unit Test Before",
	}
	groupCreateOutput := border0client.Group{
		ID:          "uuiduuid-uuid-uuid-uuid-uuiduuiduuid",
		DisplayName: groupCreateInput.DisplayName,
	}

	groupMembershipsUpdateInput := border0client.Group{
		ID:          groupCreateOutput.ID,
		DisplayName: groupCreateOutput.DisplayName,
	}
	groupMembershipsUpdateInputMembers := []string{memberIDA}

	groupMembershipsUpdateOutput := border0client.Group{
		ID:          groupMembershipsUpdateInput.ID,
		DisplayName: groupMembershipsUpdateInput.DisplayName,
		Members:     []border0client.User{{ID: groupMembershipsUpdateInputMembers[0]}},
	}

	updateInput := border0client.Group{
		ID:          groupCreateOutput.ID,
		DisplayName: "Unit Test After",
	}
	updateOutput := border0client.Group{
		ID:          updateInput.ID,
		DisplayName: updateInput.DisplayName,
	}

	userConfigStep1 := fmt.Sprintf(`
		resource "border0_group" "unit_test" {
			display_name    = "%s"
			members         = [ "%s" ]
		}`,
		groupCreateInput.DisplayName,
		memberIDA,
	)

	userConfigStep2 := fmt.Sprintf(`
		resource "border0_group" "unit_test" {
			display_name    = "%s"
			members         = []
		}`,
		updateOutput.DisplayName,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateGroup(matchContext, &groupCreateInput).Return(&groupCreateOutput, nil).Call,
		clientMock.EXPECT().UpdateGroupMemberships(matchContext, &groupMembershipsUpdateInput, groupMembershipsUpdateInputMembers).Return(&groupMembershipsUpdateOutput, nil).Call,
		clientMock.EXPECT().Group(matchContext, groupCreateOutput.ID).Return(&groupMembershipsUpdateOutput, nil).Call,
		clientMock.EXPECT().Group(matchContext, groupCreateOutput.ID).Return(&groupMembershipsUpdateOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().Group(matchContext, groupCreateOutput.ID).Return(&groupMembershipsUpdateOutput, nil).Call,

		// terraform apply (update + read + read)
		clientMock.EXPECT().UpdateGroupMemberships(matchContext, &border0client.Group{ID: groupCreateOutput.ID}, []string{}).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().UpdateGroup(matchContext, &updateInput).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Group(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Group(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().Group(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteGroup(matchContext, updateOutput.ID).Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: userConfigStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_group.unit_test", "display_name", groupCreateOutput.DisplayName),
					resource.TestCheckResourceAttr("border0_group.unit_test", "members.#", "1"),
					resource.TestCheckResourceAttr("border0_group.unit_test", "members.0", memberIDA),
					resource.TestCheckResourceAttrSet("border0_group.unit_test", "id"),
				),
			},
			{
				Config: userConfigStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_group.unit_test", "display_name", updateOutput.DisplayName),
					resource.TestCheckResourceAttr("border0_group.unit_test", "members.#", "0"),
					resource.TestCheckResourceAttrSet("border0_group.unit_test", "id"),
				),
			},
			{
				ResourceName:      "border0_group.unit_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
