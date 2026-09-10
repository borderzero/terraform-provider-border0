package border0_test

import (
	"fmt"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_Resource_Border0ApprovalFlow(t *testing.T) {
	socketID := "socketid-uuid-uuid-uuid-uuiduuiduuid"
	requesterGroupID := "reqgroup-uuid-uuid-uuid-uuiduuiduuid"
	approverUserID := "apprusrs-uuid-uuid-uuid-uuiduuiduuid"

	createInput := border0client.ApprovalWorkflow{
		Name:              "unit test flow",
		Description:       "unit test flow description",
		SocketIDs:         []string{socketID},
		RequesterGroupIDs: []string{requesterGroupID},
		ApproverUserIDs:   []string{approverUserID},
		AllowSelfApproval: false,
	}
	createOutput := createInput
	createOutput.ID = "flowuuid-uuid-uuid-uuid-uuiduuiduuid"

	updateInput := createOutput
	updateInput.Name = "unit test flow after"
	updateInput.AllowSelfApproval = true
	updateOutput := updateInput

	configStep1 := fmt.Sprintf(`
		resource "border0_approval_flow" "unit_test" {
			name                = "%s"
			description         = "%s"
			socket_ids          = [ "%s" ]
			requester_group_ids = [ "%s" ]
			approver_user_ids   = [ "%s" ]
		}`,
		createInput.Name,
		createInput.Description,
		socketID,
		requesterGroupID,
		approverUserID,
	)

	configStep2 := fmt.Sprintf(`
		resource "border0_approval_flow" "unit_test" {
			name                = "%s"
			description         = "%s"
			socket_ids          = [ "%s" ]
			requester_group_ids = [ "%s" ]
			approver_user_ids   = [ "%s" ]
			allow_self_approval = true
		}`,
		updateInput.Name,
		updateInput.Description,
		socketID,
		requesterGroupID,
		approverUserID,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateApprovalWorkflow(matchContext, &createInput).Return(&createOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, createOutput.ID).Return(&createOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, createOutput.ID).Return(&createOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().ApprovalWorkflow(matchContext, createOutput.ID).Return(&createOutput, nil).Call,

		// terraform apply (update + read + read)
		clientMock.EXPECT().UpdateApprovalWorkflow(matchContext, &updateInput).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().ApprovalWorkflow(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteApprovalWorkflow(matchContext, updateOutput.ID).Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: configStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "name", createOutput.Name),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "description", createOutput.Description),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "socket_ids.#", "1"),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "socket_ids.0", socketID),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "requester_group_ids.0", requesterGroupID),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "approver_user_ids.0", approverUserID),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "allow_self_approval", "false"),
					resource.TestCheckResourceAttrSet("border0_approval_flow.unit_test", "id"),
				),
			},
			{
				Config: configStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "name", updateOutput.Name),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test", "allow_self_approval", "true"),
					resource.TestCheckResourceAttrSet("border0_approval_flow.unit_test", "id"),
				),
			},
			{
				ResourceName:      "border0_approval_flow.unit_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func Test_Resource_Border0ApprovalFlow_SocketTags(t *testing.T) {
	requesterUserID := "requsers-uuid-uuid-uuid-uuiduuiduuid"
	approverGroupID := "apprgrps-uuid-uuid-uuid-uuiduuiduuid"

	createInput := border0client.ApprovalWorkflow{
		Name:              "unit test tags flow",
		Description:       "unit test tags flow description",
		SocketTags:        map[string]string{"env": "prod"},
		RequesterUserIDs:  []string{requesterUserID},
		ApproverGroupIDs:  []string{approverGroupID},
		AllowSelfApproval: false,
	}
	createOutput := createInput
	createOutput.ID = "tagsflow-uuid-uuid-uuid-uuiduuiduuid"

	updateInput := createOutput
	updateInput.SocketTags = map[string]string{"env": "staging", "team": "sre"}
	updateOutput := updateInput

	// switch from socket tags to socket ids, this should clear the tags from the state
	socketID := "socketid-uuid-uuid-uuid-uuiduuiduuid"
	noTagsInput := updateOutput
	noTagsInput.SocketTags = nil
	noTagsInput.SocketIDs = []string{socketID}
	noTagsOutput := noTagsInput

	configOneTag := fmt.Sprintf(`
		resource "border0_approval_flow" "unit_test_tags" {
			name               = "%s"
			description        = "%s"
			socket_tags        = { env = "prod" }
			requester_user_ids = [ "%s" ]
			approver_group_ids = [ "%s" ]
		}`,
		createInput.Name,
		createInput.Description,
		requesterUserID,
		approverGroupID,
	)

	configTwoTags := fmt.Sprintf(`
		resource "border0_approval_flow" "unit_test_tags" {
			name               = "%s"
			description        = "%s"
			socket_tags        = { env = "staging", team = "sre" }
			requester_user_ids = [ "%s" ]
			approver_group_ids = [ "%s" ]
		}`,
		updateInput.Name,
		updateInput.Description,
		requesterUserID,
		approverGroupID,
	)

	configNoTags := fmt.Sprintf(`
		resource "border0_approval_flow" "unit_test_tags" {
			name               = "%s"
			description        = "%s"
			socket_ids         = [ "%s" ]
			requester_user_ids = [ "%s" ]
			approver_group_ids = [ "%s" ]
		}`,
		noTagsInput.Name,
		noTagsInput.Description,
		socketID,
		requesterUserID,
		approverGroupID,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateApprovalWorkflow(matchContext, &createInput).Return(&createOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, createOutput.ID).Return(&createOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, createOutput.ID).Return(&createOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().ApprovalWorkflow(matchContext, createOutput.ID).Return(&createOutput, nil).Call,

		// terraform apply (update + read + read)
		clientMock.EXPECT().UpdateApprovalWorkflow(matchContext, &updateInput).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().ApprovalWorkflow(matchContext, updateOutput.ID).Return(&updateOutput, nil).Call,

		// terraform apply (update + read + read), socket tags cleared
		clientMock.EXPECT().UpdateApprovalWorkflow(matchContext, &noTagsInput).Return(&noTagsOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, noTagsOutput.ID).Return(&noTagsOutput, nil).Call,
		clientMock.EXPECT().ApprovalWorkflow(matchContext, noTagsOutput.ID).Return(&noTagsOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteApprovalWorkflow(matchContext, noTagsOutput.ID).Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: configOneTag,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_tags.%", "1"),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_tags.env", "prod"),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_ids.#", "0"),
					resource.TestCheckResourceAttrSet("border0_approval_flow.unit_test_tags", "id"),
				),
			},
			{
				Config: configTwoTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_tags.%", "2"),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_tags.env", "staging"),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_tags.team", "sre"),
				),
			},
			{
				Config: configNoTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_tags.%", "0"),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_ids.#", "1"),
					resource.TestCheckResourceAttr("border0_approval_flow.unit_test_tags", "socket_ids.0", socketID),
				),
			},
		},
	})
}
