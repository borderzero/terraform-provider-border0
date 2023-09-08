package border0_test

import (
	"testing"

	"github.com/autarch/testify/require"
	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var connectorTokenNeverExpiresConfig = `
resource "border0_connector_token" "unit_test_never_expires" {
  connector_id = "unit-test-connector-id"
  name = "unit-test-connector-token-never-expires"
}
`

var connmectorTokenExpiresConfig = `
resource "border0_connector_token" "unit_test_expires" {
  connector_id = "unit-test-connector-id"
  name = "unit-test-connector-token-never-expires"
  expires_at = "2023-12-31T23:59:59Z"
}
`

func Test_Resource_Border0ConnectorToken_NeverExpires(t *testing.T) {
	createdAt, err := border0client.FlexibleTimeFrom("2020-01-02T15:04:05Z")
	require.NoError(t, err)

	input := border0client.ConnectorToken{
		ConnectorID: "unit-test-connector-id",
		Name:        "unit-test-connector-token-never-expires",
	}
	output := border0client.ConnectorToken{
		ConnectorID: "unit-test-connector-id",
		Name:        "unit-test-connector-token-never-expires",
		ID:          "unit-test-connector-token-id",
		Token:       "unit-test-connector-token",
		CreatedBy:   "bilbo.baggins@border0.com",
		CreatedAt:   createdAt,
	}
	tokensOutput := border0client.ConnectorTokens{
		List: []border0client.ConnectorToken{output},
	}

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// read = client.ConnectorTokens()
		// create = client.CreateConnectorToken()
		// delete = client.DeleteConnectorToken()

		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateConnectorToken(matchContext, &input).Return(&output, nil).Call,
		clientMock.EXPECT().ConnectorTokens(matchContext, "unit-test-connector-id").Return(&tokensOutput, nil).Call,
		clientMock.EXPECT().ConnectorTokens(matchContext, "unit-test-connector-id").Return(&tokensOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().ConnectorTokens(matchContext, "unit-test-connector-id").Return(&tokensOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteConnectorToken(matchContext, "unit-test-connector-id", "unit-test-connector-token-id").Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: connectorTokenNeverExpiresConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_connector_token.unit_test_never_expires", "name", "unit-test-connector-token-never-expires"),
					resource.TestCheckResourceAttr("border0_connector_token.unit_test_never_expires", "connector_id", "unit-test-connector-id"),
					resource.TestCheckResourceAttrSet("border0_connector_token.unit_test_never_expires", "id"),
				),
			},
			{
				ResourceName:      "border0_connector_token.unit_test_never_expires",
				ImportState:       true,
				ImportStateVerify: false, // skip verification because the token a computed value from the API
			},
		},
	})
}

func Test_Resource_Border0ConnectorToken_Expires(t *testing.T) {
	createdAt, err := border0client.FlexibleTimeFrom("2020-01-02T15:04:05Z")
	require.NoError(t, err)

	expiresAt, err := border0client.FlexibleTimeFrom("2023-12-31T23:59:59Z")
	require.NoError(t, err)

	input := border0client.ConnectorToken{
		ConnectorID: "unit-test-connector-id",
		Name:        "unit-test-connector-token-never-expires",
		ExpiresAt:   expiresAt,
	}
	output := border0client.ConnectorToken{
		ConnectorID: "unit-test-connector-id",
		Name:        "unit-test-connector-token-never-expires",
		ExpiresAt:   expiresAt,
		ID:          "unit-test-connector-token-id",
		Token:       "unit-test-connector-token",
		CreatedBy:   "bilbo.baggins@border0.com",
		CreatedAt:   createdAt,
	}
	tokensOutput := border0client.ConnectorTokens{
		List: []border0client.ConnectorToken{output},
	}

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// read = client.ConnectorTokens()
		// create = client.CreateConnectorToken()
		// delete = client.DeleteConnectorToken()

		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateConnectorToken(matchContext, &input).Return(&output, nil).Call,
		clientMock.EXPECT().ConnectorTokens(matchContext, "unit-test-connector-id").Return(&tokensOutput, nil).Call,
		clientMock.EXPECT().ConnectorTokens(matchContext, "unit-test-connector-id").Return(&tokensOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().ConnectorTokens(matchContext, "unit-test-connector-id").Return(&tokensOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteConnectorToken(matchContext, "unit-test-connector-id", "unit-test-connector-token-id").Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: connmectorTokenExpiresConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_connector_token.unit_test_expires", "name", "unit-test-connector-token-never-expires"),
					resource.TestCheckResourceAttr("border0_connector_token.unit_test_expires", "connector_id", "unit-test-connector-id"),
					resource.TestCheckResourceAttrSet("border0_connector_token.unit_test_expires", "id"),
				),
			},
			{
				ResourceName:      "border0_connector_token.unit_test_expires",
				ImportState:       true,
				ImportStateVerify: false, // skip verification because the token a computed value from the API
			},
		},
	})
}
