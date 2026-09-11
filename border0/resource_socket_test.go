package border0_test

import (
	"context"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/border0-go/client/enum"
	"github.com/borderzero/border0-go/types/service"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

var httpSocketConfig = `
resource "border0_socket" "unit_test_http" {
  name = "unit-test-http-socket"
  description = "socket created from terraform unit test"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
}
`

var httpSocketConfig_update = `
resource "border0_socket" "unit_test_http" {
  name = "unit-test-http-socket"
  description = "update socket description"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
}
`

var sshSocketConfig_awsEC2InstanceConnect = `
resource "border0_socket" "unit_test_http" {
  name = "unit-test-http-socket"
  description = "socket created from terraform unit test"
  socket_type = "http"
  tags = {
    "test_key_1" = "test_value_1"
  }
}
`

func Test_Resource_Border0Socket_HTTPBasic(t *testing.T) {
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
		// read = client.Socket() + client.Socket() + client.SocketConnectors() + client.SocketUpstreamConfigs()
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
				Config: httpSocketConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "name", "unit-test-http-socket"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "description", "socket created from terraform unit test"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "socket_type", "http"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_http", "tags.test_key_1", "test_value_1"),
					resource.TestCheckResourceAttrSet("border0_socket.unit_test_http", "id"),
				),
			},
			{
				Config: httpSocketConfig_update,
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

// kubernetes_configuration is optional, so this socket deliberately omits it
var kubernetesSocketConfig = `
resource "border0_socket" "unit_test_kubernetes" {
  name          = "unit-test-kubernetes-socket"
  socket_type   = "kubernetes"
  connector_ids = ["unit-test-connector-id"]
}
`

// a kubernetes socket with no configuration still gets this back from the api,
// which used to show up as a diff on every single plan
func Test_Resource_Border0Socket_KubernetesWithoutConfiguration(t *testing.T) {
	socketOutput := border0client.Socket{
		SocketID:     "unit-test-kubernetes-socket-id",
		Name:         "unit-test-kubernetes-socket",
		SocketType:   "kubernetes",
		UpstreamType: "kubernetes",
	}
	connectorsOutput := border0client.SocketConnectors{
		List: []border0client.SocketConnector{
			{ConnectorID: "unit-test-connector-id", SocketID: socketOutput.SocketID},
		},
	}
	upstreamConfigsOutput := kubernetesUpstreamConfigs("")

	// call order and counts do not matter here, the empty plan is what we care about
	clientMock := mocks.APIClientRequester{}
	clientMock.EXPECT().CreateSocket(matchContext, mock.Anything).Return(&socketOutput, nil)
	clientMock.EXPECT().Socket(matchContext, socketOutput.SocketID).Return(&socketOutput, nil)
	clientMock.EXPECT().SocketConnectors(matchContext, socketOutput.SocketID).Return(&connectorsOutput, nil)
	clientMock.EXPECT().SocketUpstreamConfigs(matchContext, socketOutput.SocketID).Return(upstreamConfigsOutput, nil)
	clientMock.EXPECT().DeleteSocket(matchContext, socketOutput.SocketID).Return(nil)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				// no ExpectNonEmptyPlan, so leftover drift fails this test
				Config: kubernetesSocketConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_socket.unit_test_kubernetes", "socket_type", "kubernetes"),
					resource.TestCheckResourceAttr("border0_socket.unit_test_kubernetes", "kubernetes_configuration.#", "0"),
					resource.TestCheckResourceAttrSet("border0_socket.unit_test_kubernetes", "id"),
				),
			},
		},
	})
}

var kubernetesSocketConfig_withServer = `
resource "border0_socket" "unit_test_kubernetes" {
  name          = "unit-test-kubernetes-socket"
  socket_type   = "kubernetes"
  connector_ids = ["unit-test-connector-id"]

  kubernetes_configuration {
    server = "https://10.0.0.1:6443"
  }
}
`

// deleting the block must really delete the settings. fails if the block is Computed.
func Test_Resource_Border0Socket_KubernetesConfigurationRemoved(t *testing.T) {
	socketOutput := border0client.Socket{
		SocketID:     "unit-test-kubernetes-socket-id",
		Name:         "unit-test-kubernetes-socket",
		SocketType:   "kubernetes",
		UpstreamType: "kubernetes",
	}
	connectorsOutput := border0client.SocketConnectors{
		List: []border0client.SocketConnector{
			{ConnectorID: "unit-test-connector-id", SocketID: socketOutput.SocketID},
		},
	}

	// the mock remembers what the api holds, so step 2 sees the cleared server
	upstreamConfigs := kubernetesUpstreamConfigs("https://10.0.0.1:6443")

	clientMock := mocks.APIClientRequester{}
	clientMock.EXPECT().CreateSocket(matchContext, mock.Anything).Return(&socketOutput, nil)
	clientMock.EXPECT().Socket(matchContext, socketOutput.SocketID).Return(&socketOutput, nil)
	clientMock.EXPECT().SocketConnectors(matchContext, socketOutput.SocketID).Return(&connectorsOutput, nil)
	clientMock.EXPECT().SocketUpstreamConfigs(matchContext, socketOutput.SocketID).
		RunAndReturn(func(context.Context, string) (*border0client.SocketUpstreamConfigs, error) {
			return upstreamConfigs, nil
		})
	clientMock.EXPECT().DeleteSocket(matchContext, socketOutput.SocketID).Return(nil)

	// the point of the test. removing the block must reach the api as an empty server
	clientMock.EXPECT().
		UpdateSocket(matchContext, socketOutput.SocketID, mock.MatchedBy(clearsKubernetesServer)).
		RunAndReturn(func(context.Context, string, *border0client.Socket) (*border0client.Socket, error) {
			upstreamConfigs = kubernetesUpstreamConfigs("")
			return &socketOutput, nil
		}).Once()

	// without this, a skipped update would let the test pass
	t.Cleanup(func() { clientMock.AssertExpectations(t) })

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: kubernetesSocketConfig_withServer,
				Check: resource.TestCheckResourceAttr(
					"border0_socket.unit_test_kubernetes",
					"kubernetes_configuration.0.server", "https://10.0.0.1:6443"),
			},
			{
				// block is gone, so the server must be gone too
				Config: kubernetesSocketConfig,
			},
		},
	})
}

var kubernetesSocketConfig_emptyBlock = `
resource "border0_socket" "unit_test_kubernetes" {
  name          = "unit-test-kubernetes-socket"
  socket_type   = "kubernetes"
  connector_ids = ["unit-test-connector-id"]

  kubernetes_configuration {}
}
`

// an empty block says nothing to the api, but it was written on purpose, so it
// must stay in state. if it vanishes the plan wants to add it back forever.
func Test_Resource_Border0Socket_KubernetesEmptyConfiguration(t *testing.T) {
	socketOutput := border0client.Socket{
		SocketID:     "unit-test-kubernetes-socket-id",
		Name:         "unit-test-kubernetes-socket",
		SocketType:   "kubernetes",
		UpstreamType: "kubernetes",
	}
	connectorsOutput := border0client.SocketConnectors{
		List: []border0client.SocketConnector{
			{ConnectorID: "unit-test-connector-id", SocketID: socketOutput.SocketID},
		},
	}

	// call order and counts do not matter here, the empty plan is what we care about
	clientMock := mocks.APIClientRequester{}
	clientMock.EXPECT().CreateSocket(matchContext, mock.Anything).Return(&socketOutput, nil)
	clientMock.EXPECT().Socket(matchContext, socketOutput.SocketID).Return(&socketOutput, nil)
	clientMock.EXPECT().SocketConnectors(matchContext, socketOutput.SocketID).Return(&connectorsOutput, nil)
	clientMock.EXPECT().SocketUpstreamConfigs(matchContext, socketOutput.SocketID).Return(kubernetesUpstreamConfigs(""), nil)
	clientMock.EXPECT().DeleteSocket(matchContext, socketOutput.SocketID).Return(nil)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				// no ExpectNonEmptyPlan, so a vanishing block fails this test
				Config: kubernetesSocketConfig_emptyBlock,
				Check: resource.TestCheckResourceAttr(
					"border0_socket.unit_test_kubernetes", "kubernetes_configuration.#", "1"),
			},
		},
	})
}

// kubernetesUpstreamConfigs is what the api sends back for the given server.
func kubernetesUpstreamConfigs(server string) *border0client.SocketUpstreamConfigs {
	return &border0client.SocketUpstreamConfigs{
		List: []border0client.SocketUpstreamConfig{
			{
				Config: service.Configuration{
					ServiceType: "kubernetes",
					KubernetesServiceConfiguration: &service.KubernetesServiceConfiguration{
						KubernetesServiceType: service.KubernetesServiceTypeStandard,
						StandardKubernetesServiceConfiguration: &service.StandardKubernetesServiceConfiguration{
							Server: server,
						},
					},
				},
			},
		},
	}
}

// clearsKubernetesServer reports whether a socket update wipes the server.
func clearsKubernetesServer(socket *border0client.Socket) bool {
	if socket.UpstreamConfig == nil || socket.UpstreamConfig.KubernetesServiceConfiguration == nil {
		return false
	}
	standard := socket.UpstreamConfig.KubernetesServiceConfiguration.StandardKubernetesServiceConfiguration
	return standard != nil && standard.Server == ""
}
