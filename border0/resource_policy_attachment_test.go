package border0_test

import (
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var policyAttachmentConfig = `
resource "border0_policy_attachment" "unit_test" {
  policy_id = "unit-test-policy-id"
  socket_id = "unit-test-socket-id"
}
`

func Test_Resource_Border0PolicyAttachment(t *testing.T) {
	policyID := "unit-test-policy-id"
	socketID := "unit-test-socket-id"

	policy := border0client.Policy{
		SocketIDs: []string{socketID},
		// other policy fields are omitted in this test
	}

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// read = client.Policy(
		// create = client.AttachPolicyToSocket()
		// delete = client.RemovePolicyFromSocket()

		// terraform apply (create + read + read)
		clientMock.EXPECT().AttachPolicyToSocket(matchContext, policyID, socketID).Return(nil).Call,
		clientMock.EXPECT().Policy(matchContext, policyID).Return(&policy, nil).Call,
		clientMock.EXPECT().Policy(matchContext, policyID).Return(&policy, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().Policy(matchContext, policyID).Return(&policy, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().RemovePolicyFromSocket(matchContext, policyID, socketID).Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: policyAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_policy_attachment.unit_test", "policy_id", "unit-test-policy-id"),
					resource.TestCheckResourceAttr("border0_policy_attachment.unit_test", "socket_id", "unit-test-socket-id"),
					resource.TestCheckResourceAttrSet("border0_policy_attachment.unit_test", "id"),
				),
			},
			{
				ResourceName:      "border0_policy_attachment.unit_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
