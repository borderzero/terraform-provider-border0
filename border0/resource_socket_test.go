package border0_test

import (
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/client/enum"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialHTTPSocketConfig = `
resource "border0_socket" "unit_test_http" {
  name = "unit-test-http-socket"
  description = "socket created from terraform unit test"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
}
`

var updateHTTPSocketConfig = `
resource "border0_socket" "unit_test_http" {
  name = "unit-test-http-socket"
  description = "update socket description"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
}
`

func Test_Resource_Border0Socket_HTTP(t *testing.T) {
	initialInput := border0client.Socket{
		Name:        "unit-test-http-socket",
		Description: "socket created from terraform unit test",
		SocketType:  enum.SocketTypeHTTP,
		Tags: map[string]string{
			"test_key_1": "test_value_1",
		},
		UpstreamType: "http",
	}
	initialOutput := border0client.Socket{
		SocketID:    "unit-test-http-socket-id",
		Name:        "unit-test-http-socket",
		Description: "socket created from terraform unit test",
		SocketType:  enum.SocketTypeHTTP,
		Tags: map[string]string{
			"test_key_1": "test_value_1",
		},
		UpstreamType: "http",
	}
	updateInput := border0client.Socket{
		Name:        "unit-test-http-socket",
		Description: "update socket description",
		SocketType:  enum.SocketTypeHTTP,
		Tags: map[string]string{
			"test_key_1": "test_value_1",
		},
		UpstreamType: "http",
	}
	updateOutput := border0client.Socket{
		SocketID:    "unit-test-http-socket-id",
		Name:        "unit-test-http-socket",
		Description: "update socket description",
		SocketType:  enum.SocketTypeHTTP,
		Tags: map[string]string{
			"test_key_1": "test_value_1",
		},
		UpstreamType: "http",
	}

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// read = client.Socket() + client.SocketConnectors() + client.SocketUpstreamConfigs()
		// create = client.CreateSocket()
		// update = client.Socket() + client.UpdateSocket()
		// delete = client.DeleteSocket()

		// terraform apply (create + read + read)
		clientMock.EXPECT().CreateSocket(matchContext, &initialInput).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Socket(matchContext, "unit-test-http-socket-id").Return(&initialOutput, nil).Call,
		clientMock.EXPECT().SocketConnectors(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketConnectors), nil).Call,
		clientMock.EXPECT().SocketUpstreamConfigs(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketUpstreamConfigs), nil).Call,
		clientMock.EXPECT().Socket(matchContext, "unit-test-http-socket-id").Return(&initialOutput, nil).Call,
		clientMock.EXPECT().SocketConnectors(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketConnectors), nil).Call,
		clientMock.EXPECT().SocketUpstreamConfigs(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketUpstreamConfigs), nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().Socket(matchContext, "unit-test-http-socket-id").Return(&initialOutput, nil).Call,
		clientMock.EXPECT().SocketConnectors(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketConnectors), nil).Call,
		clientMock.EXPECT().SocketUpstreamConfigs(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketUpstreamConfigs), nil).Call,

		// terraform apply (update + read + read)
		// update needs to fetch socket before updating socket
		clientMock.EXPECT().Socket(matchContext, "unit-test-http-socket-id").Return(&updateOutput, nil).Call,
		clientMock.EXPECT().UpdateSocket(matchContext, "unit-test-http-socket-id", &updateInput).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Socket(matchContext, "unit-test-http-socket-id").Return(&updateOutput, nil).Call,
		clientMock.EXPECT().SocketConnectors(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketConnectors), nil).Call,
		clientMock.EXPECT().SocketUpstreamConfigs(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketUpstreamConfigs), nil).Call,
		clientMock.EXPECT().Socket(matchContext, "unit-test-http-socket-id").Return(&updateOutput, nil).Call,
		clientMock.EXPECT().SocketConnectors(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketConnectors), nil).Call,
		clientMock.EXPECT().SocketUpstreamConfigs(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketUpstreamConfigs), nil).Call,

		// terraform import (read)
		clientMock.EXPECT().Socket(matchContext, "unit-test-http-socket-id").Return(&updateOutput, nil).Call,
		clientMock.EXPECT().SocketConnectors(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketConnectors), nil).Call,
		clientMock.EXPECT().SocketUpstreamConfigs(matchContext, "unit-test-http-socket-id").Return(new(border0client.SocketUpstreamConfigs), nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeleteSocket(matchContext, "unit-test-http-socket-id").Return(nil).Call,
	)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: initialHTTPSocketConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "name", "unit-test-http-socket"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "description", "socket created from terraform unit test"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "socket_type", "http"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "tags.test_key_1", "test_value_1"),
					resource.TestCheckResourceAttrSet("border0_socket.unit_test_http", "id"),
				),
			},
			{
				Config: updateHTTPSocketConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "name", "unit-test-http-socket"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "description", "update socket description"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "socket_type", "http"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "tags.test_key_1", "test_value_1"),
				),
			},
			{
				ResourceName:      "border0_socket.unit_test_http",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
