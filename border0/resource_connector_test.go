package border0_test

import (
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialConnectorConfig = `
resource "border0_connector" "unit_test" {
  name = "unit-test-connector-1"
  description = "connector created from terraform unit test"
}
`
var updateConnectorConfig = `
resource "border0_connector" "unit_test" {
  name = "unit-test-connector-1"
  description = "update connector description"
}
`

func Test_Resource_Border0Connector(t *testing.T) {
	initialInput := border0client.Connector{
		Name:        "unit-test-connector-1",
		Description: "connector created from terraform unit test",
	}
	initialOutput := border0client.Connector{
		ConnectorID: "unit-test-id-1",
		Name:        "unit-test-connector-1",
		Description: "connector created from terraform unit test",
	}
	updateInputOutput := border0client.Connector{
		ConnectorID: "unit-test-id-1",
		Name:        "unit-test-connector-1",
		Description: "update connector description",
	}

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// read = client.Connector()
		// create = client.CreateConnector()
		// update = client.UpdateConnector()
		// delete = client.DeleteConnector()

		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateConnector(matchContext, &initialInput).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Connector(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Connector(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().Connector(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,

		// terraform aplly (update + read + read)
		clientMock.EXPECT().UpdateConnector(matchContext, &updateInputOutput).Return(&updateInputOutput, nil).Call,
		clientMock.EXPECT().Connector(matchContext, "unit-test-id-1").Return(&updateInputOutput, nil).Call,
		clientMock.EXPECT().Connector(matchContext, "unit-test-id-1").Return(&updateInputOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().Connector(matchContext, "unit-test-id-1").Return(&updateInputOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteConnector(matchContext, "unit-test-id-1").Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: initialConnectorConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_connector.unit_test", "name", "unit-test-connector-1"),
					resource.TestCheckResourceAttr("border0_connector.unit_test", "description", "connector created from terraform unit test"),
					resource.TestCheckResourceAttrSet("border0_connector.unit_test", "id"),
				),
			},
			{
				Config: updateConnectorConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_connector.unit_test", "name", "unit-test-connector-1"),
					resource.TestCheckResourceAttr("border0_connector.unit_test", "description", "update connector description"),
				),
			},
			{
				ResourceName:      "border0_connector.unit_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
