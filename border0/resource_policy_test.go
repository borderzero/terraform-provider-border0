package border0_test

import (
	"encoding/json"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

var initialPolicyConfig = `
resource "border0_policy" "unit_test" {
  name = "unit-test-policy-1"
  version = "v1"
  description = "policy created from terraform unit test"
  policy_data = jsonencode({
    "version": "v1",
    "action": [ "database", "ssh", "http", "tls" ],
    "condition": {
      "who": {
        "email": [ "johndoe@example.com" ],
        "group": [ "db5c2352-b689-4135-babc-e97a8893128b" ],
        "domain": [ "example.com" ]
      },
      "where": {
        "allowed_ip": [ "0.0.0.0/0", "::/0" ],
        "country": [ "NL", "CA", "US", "BR", "FR" ],
        "country_not": [ "BE" ]
      },
      "when": {
        "after": "2022-10-13T05:12:26Z",
        "before": null,
        "time_of_day_after": "00:00 UTC",
        "time_of_day_before": "23:59 UTC"
      }
    }
  })
}
`

var updatePolicyConfig = `
resource "border0_policy" "unit_test" {
  name = "unit-test-policy-1"
  description = "update policy description"
  policy_data = jsonencode({
    "version": "v1",
    "action": [ "database", "ssh", "http", "tls" ],
    "condition": {
      "who": {
        "email": [ "johndoe@example.com", "another@example.com" ],
        "group": [ "db5c2352-b689-4135-babc-e97a8893128b" ],
        "domain": [ "example.com" ]
      },
      "where": {
        "allowed_ip": [ "0.0.0.0/0", "::/0" ],
        "country": [ "NL", "CA", "US", "BR", "FR" ],
        "country_not": [ "BE" ]
      },
      "when": {
        "after": "2022-10-13T05:12:26Z",
        "before": null,
        "time_of_day_after": "00:00 UTC",
        "time_of_day_before": "23:59 UTC"
      }
    }
  })
}
`

var initialPolicyConfigV2 = `
resource "border0_policy" "unit_test_v2" {
  name = "unit-test-policy-v2"
  description = "policy created from terraform unit test"
  version = "v2"
  policy_data = jsonencode({
	"permissions": {
        "database": {
            "allowed_databases": [
                {
                    "database": "videos",
                    "allowed_query_types": [
                        "ReadOnly",
                    ]
                }
            ]
        },
        "http": {},
        "rdp": {},
        "ssh": {
            "docker_exec": {
                "allowed_containers": [
                    "api-api-1"
                ]
            },
            "exec": {},
            "kubectl_exec": {},
            "sftp": {},
            "shell": {},
            "tcp_forwarding": {}
        },
        "tls": {},
        "vnc": {},
        "vpn": {}
    },
    "condition": {
      "who": {
        "email": [ "johndoe@example.com" ],
        "group": [ "db5c2352-b689-4135-babc-e97a8893128b" ],
        "service_account": [ "test-sa" ]
      },
      "where": {
        "allowed_ip": [ "0.0.0.0/0", "::/0" ],
        "country": [ "NL", "CA", "US", "BR", "FR" ],
        "country_not": [ "BE" ]
      },
      "when": {
        "after": "2022-10-13T05:12:26Z",
        "before": null,
        "time_of_day_after": "00:00 UTC",
        "time_of_day_before": "23:59 UTC"
      }
    }
  })
}
`

var updatePolicyConfigV2 = `
resource "border0_policy" "unit_test_v2" {
  name = "unit-test-policy-v2"
  version = "v2"
  description = "update policy description"
  policy_data = jsonencode({
	"permissions": {
        "database": {
            "allowed_databases": [
                {
                    "database": "books",
                    "allowed_query_types": [
                        "ReadOnly",
                        "USE"
                    ]
                }
            ]
        },
        "http": {},
        "rdp": {},
        "ssh": {
            "docker_exec": {
                "allowed_containers": [
                    "api-api-2"
                ]
            },
            "exec": {},
            "kubectl_exec": {},
            "sftp": {},
            "shell": {},
            "tcp_forwarding": {}
        },
        "vpn": {}
    },
    "condition": {
      "who": {
        "email": [ "johndoe@example.com", "another@example.com" ],
        "group": [ "db5c2352-b689-4135-babc-e97a8893128b" ],
        "service_account": [ "test-sa" ]
      },
      "where": {
        "allowed_ip": [ "0.0.0.0/0", "::/0" ],
        "country": [ "NL", "CA", "US", "BR", "FR" ],
        "country_not": [ "BE" ]
      },
      "when": {
        "after": "2022-10-13T05:12:26Z",
        "before": null,
        "time_of_day_after": "00:00 UTC",
        "time_of_day_before": "23:59 UTC"
      }
    }
  })
}
`

func Test_Resource_Border0Policy(t *testing.T) {
	initialPolicyData := border0client.PolicyData{
		Version: "v1",
		Action:  []string{"database", "ssh", "http", "tls"},
		Condition: border0client.PolicyCondition{
			Who: border0client.PolicyWho{
				Email:  []string{"johndoe@example.com"},
				Group:  []string{"db5c2352-b689-4135-babc-e97a8893128b"},
				Domain: []string{"example.com"},
			},
			Where: border0client.PolicyWhere{
				AllowedIP:  []string{"0.0.0.0/0", "::/0"},
				Country:    []string{"NL", "CA", "US", "BR", "FR"},
				CountryNot: []string{"BE"},
			},
			When: border0client.PolicyWhen{
				After:           "2022-10-13T05:12:26Z",
				Before:          "",
				TimeOfDayAfter:  "00:00 UTC",
				TimeOfDayBefore: "23:59 UTC",
			},
		},
	}
	initialInput := border0client.Policy{
		Name:        "unit-test-policy-1",
		Description: "policy created from terraform unit test",
		Version:     "v1",
		PolicyData:  initialPolicyData,
	}
	initialOutput := border0client.Policy{
		ID:          "unit-test-id-1",
		Version:     "v1",
		Name:        "unit-test-policy-1",
		Description: "policy created from terraform unit test",
		PolicyData:  initialPolicyData,
	}

	updatePolicyData := border0client.PolicyData{
		Version: "v1",
		Action:  []string{"database", "ssh", "http", "tls"},
		Condition: border0client.PolicyCondition{
			Who: border0client.PolicyWho{
				Email: []string{
					"johndoe@example.com",
					"another@example.com",
				},
				Group:  []string{"db5c2352-b689-4135-babc-e97a8893128b"},
				Domain: []string{"example.com"},
			},
			Where: border0client.PolicyWhere{
				AllowedIP:  []string{"0.0.0.0/0", "::/0"},
				Country:    []string{"NL", "CA", "US", "BR", "FR"},
				CountryNot: []string{"BE"},
			},
			When: border0client.PolicyWhen{
				After:           "2022-10-13T05:12:26Z",
				Before:          "",
				TimeOfDayAfter:  "00:00 UTC",
				TimeOfDayBefore: "23:59 UTC",
			},
		},
	}
	updateInput := border0client.Policy{
		Name:        "unit-test-policy-1",
		Description: "update policy description",
		PolicyData:  updatePolicyData,
	}
	updateOutput := border0client.Policy{
		ID:          "unit-test-id-1",
		Version:     "v1",
		Name:        "unit-test-policy-1",
		Description: "update policy description",
		PolicyData:  updatePolicyData,
	}

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// read = client.Policy()
		// create = client.CreatePolicy()
		// update = client.UpdatePolicy()
		// delete = client.DeletePolicy()

		// terraform apply (create + read + read)
		clientMock.EXPECT().CreatePolicy(matchContext, &initialInput).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&initialOutput, nil).Call,

		// terraform aplly (update + read + read)
		clientMock.EXPECT().UpdatePolicy(matchContext, "unit-test-id-1", &updateInput).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&updateOutput, nil).Call,

		// terraform import (read)
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-1").Return(&updateOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeletePolicy(matchContext, "unit-test-id-1").Return(nil).Call,
	)

	initialPolicyDataJSON, err := json.Marshal(initialPolicyData)
	require.NoError(t, err)

	updatePolicyDataJSON, err := json.Marshal(updatePolicyData)
	require.NoError(t, err)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: initialPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_policy.unit_test", "name", "unit-test-policy-1"),
					resource.TestCheckResourceAttr("border0_policy.unit_test", "description", "policy created from terraform unit test"),
					resource.TestCheckResourceAttrSet("border0_policy.unit_test", "id"),
					testMatchResourceAttrJSON("border0_policy.unit_test", "policy_data", string(initialPolicyDataJSON)),
				),
			},
			{
				Config: updatePolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_policy.unit_test", "name", "unit-test-policy-1"),
					resource.TestCheckResourceAttr("border0_policy.unit_test", "description", "update policy description"),
					resource.TestCheckResourceAttr("border0_policy.unit_test", "version", "v1"),
					testMatchResourceAttrJSON("border0_policy.unit_test", "policy_data", string(updatePolicyDataJSON)),
				),
			},
			{
				ResourceName:      "border0_policy.unit_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func Test_Resource_Border0PolicyV2(t *testing.T) {
	initialPolicyData := border0client.PolicyDataV2{
		Permissions: border0client.PolicyPermissions{
			Database: &border0client.DatabasePermissions{
				AllowedDatabases: &[]border0client.DatabasePermission{
					{
						Database:          "videos",
						AllowedQueryTypes: &[]string{"ReadOnly"},
					},
				},
			},
			HTTP: &border0client.HTTPPermissions{},
			RDP:  &border0client.RDPPermissions{},
			SSH: &border0client.SSHPermissions{
				DockerExec: &border0client.SSHDockerExecPermission{
					AllowedContainers: &[]string{"api-api-1"},
				},
				Exec:          &border0client.SSHExecPermission{},
				KubectlExec:   &border0client.SSHKubectlExecPermission{},
				SFTP:          &border0client.SSHSFTPPermission{},
				Shell:         &border0client.SSHShellPermission{},
				TCPForwarding: &border0client.SSHTCPForwardingPermission{},
			},
			TLS: &border0client.TLSPermissions{},
			VNC: &border0client.VNCPermissions{},
			VPN: &border0client.VPNPermissions{},
		},
		Condition: border0client.PolicyConditionV2{
			Who: border0client.PolicyWhoV2{
				Email:          []string{"johndoe@example.com"},
				Group:          []string{"db5c2352-b689-4135-babc-e97a8893128b"},
				ServiceAccount: []string{"test-sa"},
			},
			Where: border0client.PolicyWhere{
				AllowedIP:  []string{"0.0.0.0/0", "::/0"},
				Country:    []string{"NL", "CA", "US", "BR", "FR"},
				CountryNot: []string{"BE"},
			},
			When: border0client.PolicyWhen{
				After:           "2022-10-13T05:12:26Z",
				Before:          "",
				TimeOfDayAfter:  "00:00 UTC",
				TimeOfDayBefore: "23:59 UTC",
			},
		},
	}

	initialInput := border0client.Policy{
		Name:        "unit-test-policy-v2",
		Version:     "v2",
		Description: "policy created from terraform unit test",
		PolicyData:  initialPolicyData,
	}
	initialOutput := border0client.Policy{
		ID:          "unit-test-id-v2",
		Version:     "v2",
		Name:        "unit-test-policy-v2",
		Description: "policy created from terraform unit test",
		PolicyData:  initialPolicyData,
	}

	updatePolicyData := border0client.PolicyDataV2{
		Permissions: border0client.PolicyPermissions{
			Database: &border0client.DatabasePermissions{
				AllowedDatabases: &[]border0client.DatabasePermission{
					{
						Database:          "books",
						AllowedQueryTypes: &[]string{"ReadOnly", "USE"},
					},
				},
			},
			HTTP: &border0client.HTTPPermissions{},
			RDP:  &border0client.RDPPermissions{},
			SSH: &border0client.SSHPermissions{
				DockerExec: &border0client.SSHDockerExecPermission{
					AllowedContainers: &[]string{"api-api-2"},
				},
				Exec:          &border0client.SSHExecPermission{},
				KubectlExec:   &border0client.SSHKubectlExecPermission{},
				SFTP:          &border0client.SSHSFTPPermission{},
				Shell:         &border0client.SSHShellPermission{},
				TCPForwarding: &border0client.SSHTCPForwardingPermission{},
			},
			VPN: &border0client.VPNPermissions{},
		},
		Condition: border0client.PolicyConditionV2{
			Who: border0client.PolicyWhoV2{
				Email: []string{
					"johndoe@example.com",
					"another@example.com",
				},
				Group:          []string{"db5c2352-b689-4135-babc-e97a8893128b"},
				ServiceAccount: []string{"test-sa"},
			},
			Where: border0client.PolicyWhere{
				AllowedIP:  []string{"0.0.0.0/0", "::/0"},
				Country:    []string{"NL", "CA", "US", "BR", "FR"},
				CountryNot: []string{"BE"},
			},
			When: border0client.PolicyWhen{
				After:           "2022-10-13T05:12:26Z",
				Before:          "",
				TimeOfDayAfter:  "00:00 UTC",
				TimeOfDayBefore: "23:59 UTC",
			},
		},
	}
	updateInput := border0client.Policy{
		Name:        "unit-test-policy-v2",
		Description: "update policy description",
		PolicyData:  updatePolicyData,
	}
	updateOutput := border0client.Policy{
		ID:          "unit-test-id-v2",
		Name:        "unit-test-policy-v2",
		Version:     "v2",
		Description: "update policy description",
		PolicyData:  updatePolicyData,
	}

	clientMock := mocks.APIClientRequester{}
	mockCallsInOrder(
		// read = client.Policy()
		// create = client.CreatePolicy()
		// update = client.UpdatePolicy()
		// delete = client.DeletePolicy()

		// terraform apply (create + read + read)
		clientMock.EXPECT().CreatePolicy(matchContext, &initialInput).Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-v2").Return(&initialOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-v2").Return(&initialOutput, nil).Call,

		// this read is needed because of the update
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-v2").Return(&initialOutput, nil).Call,

		// terraform aplly (update + read + read)
		clientMock.EXPECT().UpdatePolicy(matchContext, "unit-test-id-v2", &updateInput).Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-v2").Return(&updateOutput, nil).Call,
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-v2").Return(&updateOutput, nil).Call,

		// // terraform import (read)
		clientMock.EXPECT().Policy(matchContext, "unit-test-id-v2").Return(&updateOutput, nil).Call,

		// terraform destroy (delete)
		clientMock.EXPECT().DeletePolicy(matchContext, "unit-test-id-v2").Return(nil).Call,
	)

	initialPolicyDataJSON, err := json.Marshal(initialPolicyData)
	require.NoError(t, err)

	updatePolicyDataJSON, err := json.Marshal(updatePolicyData)
	require.NoError(t, err)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, &clientMock),
		Steps: []resource.TestStep{
			{
				Config: initialPolicyConfigV2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_policy.unit_test_v2", "name", "unit-test-policy-v2"),
					resource.TestCheckResourceAttr("border0_policy.unit_test_v2", "description", "policy created from terraform unit test"),
					resource.TestCheckResourceAttrSet("border0_policy.unit_test_v2", "id"),
					resource.TestCheckResourceAttr("border0_policy.unit_test_v2", "version", "v2"),
					testMatchResourceAttrJSON("border0_policy.unit_test_v2", "policy_data", string(initialPolicyDataJSON)),
				),
			},
			{
				Config: updatePolicyConfigV2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("border0_policy.unit_test_v2", "name", "unit-test-policy-v2"),
					resource.TestCheckResourceAttr("border0_policy.unit_test_v2", "description", "update policy description"),
					testMatchResourceAttrJSON("border0_policy.unit_test_v2", "policy_data", string(updatePolicyDataJSON)),
				),
			},
			{
				ResourceName:      "border0_policy.unit_test_v2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
