package border0_test

import (
	"fmt"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func Test_Resource_Border0ServiceAccountToken_NeverExpires(t *testing.T) {
	createdAt, err := border0client.FlexibleTimeFrom("2020-01-02T15:04:05Z")
	require.NoError(t, err)

	serviceAccount := &border0client.ServiceAccount{
		Name:        "unit-test-service-account",
		Description: "Test Description",
		Role:        "admin",
		Active:      true,
	}

	createInput := border0client.ServiceAccountToken{
		Name: "unit-test-sacc-token-never-expires",
	}

	singleTokenListOutput := border0client.ServiceAccountToken{
		Name:      createInput.Name,
		ID:        "unit-test-sacc-token-id",
		CreatedAt: createdAt,
	}

	createOutput := border0client.ServiceAccountToken{
		Name:      singleTokenListOutput.Name,
		ID:        singleTokenListOutput.ID,
		CreatedAt: singleTokenListOutput.CreatedAt,
		Token:     "unit-test-sacc-token",
	}

	listOutput := border0client.ServiceAccountTokens{
		List: []border0client.ServiceAccountToken{singleTokenListOutput},
	}

	serviceAccountTokenNeverExpiresConfig := fmt.Sprintf(`
		resource "border0_service_account_token" "unit_test_never_expires" {
			service_account_name = "%s"
			name                 = "%s"
		}`,
		serviceAccount.Name,
		createInput.Name,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateServiceAccountToken(matchContext, serviceAccount.Name, &createInput).Return(&createOutput, nil).Call,
		clientMock.EXPECT().ServiceAccountTokens(matchContext, serviceAccount.Name).Return(&listOutput, nil).Call,
		clientMock.EXPECT().ServiceAccountTokens(matchContext, serviceAccount.Name).Return(&listOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().ServiceAccountTokens(matchContext, serviceAccount.Name).Return(&listOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteServiceAccountToken(matchContext, serviceAccount.Name, createOutput.ID).Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: serviceAccountTokenNeverExpiresConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_service_account_token.unit_test_never_expires", "service_account_name", serviceAccount.Name),
					resource.TestCheckResourceAttr("border0_service_account_token.unit_test_never_expires", "name", createOutput.Name),
					resource.TestCheckResourceAttrSet("border0_service_account_token.unit_test_never_expires", "id"),
				),
			},
			{
				ResourceName:      "border0_service_account_token.unit_test_never_expires",
				ImportState:       true,
				ImportStateVerify: false, // skip verification because the token a computed value from the API
			},
		},
	})
}

func Test_Resource_Border0ServiceAccountToken_Expires(t *testing.T) {
	createdAt, err := border0client.FlexibleTimeFrom("2020-01-02T15:04:05Z")
	require.NoError(t, err)

	expiresAt, err := border0client.FlexibleTimeFrom("2023-12-31T23:59:59Z")
	require.NoError(t, err)

	serviceAccount := &border0client.ServiceAccount{
		Name:        "unit-test-service-account",
		Description: "Test Description",
		Role:        "admin",
		Active:      true,
	}

	createInput := border0client.ServiceAccountToken{
		Name:      "unit-test-sacc-token-never-expires",
		ExpiresAt: expiresAt,
	}

	singleTokenListOutput := border0client.ServiceAccountToken{
		Name:      createInput.Name,
		ID:        "unit-test-sacc-token-id",
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}

	createOutput := border0client.ServiceAccountToken{
		Name:      singleTokenListOutput.Name,
		ID:        singleTokenListOutput.ID,
		Token:     "unit-test-sacc-token",
		CreatedAt: singleTokenListOutput.CreatedAt,
		ExpiresAt: expiresAt,
	}

	listOutput := border0client.ServiceAccountTokens{
		List: []border0client.ServiceAccountToken{singleTokenListOutput},
	}

	serviceAccountTokenExpiresConfig := fmt.Sprintf(`
		resource "border0_service_account_token" "unit_test_expires" {
			service_account_name = "%s"
			name                 = "%s"
			expires_at           = "%s"
		}`,
		serviceAccount.Name,
		createInput.Name,
		expiresAt,
	)

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateServiceAccountToken(matchContext, serviceAccount.Name, &createInput).Return(&createOutput, nil).Call,
		clientMock.EXPECT().ServiceAccountTokens(matchContext, serviceAccount.Name).Return(&listOutput, nil).Call,
		clientMock.EXPECT().ServiceAccountTokens(matchContext, serviceAccount.Name).Return(&listOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().ServiceAccountTokens(matchContext, serviceAccount.Name).Return(&listOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteServiceAccountToken(matchContext, serviceAccount.Name, createOutput.ID).Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: serviceAccountTokenExpiresConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_service_account_token.unit_test_expires", "service_account_name", serviceAccount.Name),
					resource.TestCheckResourceAttr("border0_service_account_token.unit_test_expires", "name", createOutput.Name),
					resource.TestCheckResourceAttr("border0_service_account_token.unit_test_expires", "expires_at", createOutput.ExpiresAt.String()),
					resource.TestCheckResourceAttrSet("border0_service_account_token.unit_test_expires", "id"),
				),
			},
			{
				ResourceName:      "border0_service_account_token.unit_test_expires",
				ImportState:       true,
				ImportStateVerify: false, // skip verification because the token a computed value from the API
			},
		},
	})
}
