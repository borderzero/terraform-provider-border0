package border0_test

import (
	"testing"

	"github.com/borderzero/terraform-provider-border0/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var policyDocumentV2Config = `
data "border0_policy_v2_document" "unit_test" {
	permissions {
		database {
			allowed = true
			max_session_duration_seconds = 3600
			use_allowed_databases_list = true
			allowed_databases {
				database = "test"
				use_allowed_query_types_list = true
				allowed_query_types = ["ReadOnly"]
			}
			allowed_databases {
				database = "aap"
			}
		}
		ssh {
			allowed = true
			max_session_duration_seconds = 3600
			use_allowed_usernames_list = true
			allowed_usernames = [ "bas" ]
			shell {
				allowed = true
			}
			exec {
				allowed = true
				use_commands_list = true
				commands = ["ls" ]
			}
			sftp {
				allowed = true
			}
			tcp_forwarding {
				allowed = true
				use_allowed_connections_list = true
				allowed_connections {
					destination_address = "*"
					destination_port = "*"
				}
				allowed_connections {
					destination_address = "nu.nl"
					destination_port = "443"
				}
			}
			kubectl_exec {
				allowed = true
				use_allowed_namespaces_list= true
				allowed_namespaces {
					namespace = "test"
					use_pod_selector = true
					pod_selector = {
						aap2 = "test"
					}
				}
			}
			docker_exec {
				allowed = true
				use_allowed_containers_list = true
				allowed_containers = ["test"]
			}
		}
		vpn {
			allowed = true
		}
		http {
			allowed = true
		}
		tls {
			allowed = true
		}
		rdp {
			allowed = true
		}
		vnc {
			allowed = true
		}
	}
	condition {
		who {
			email = [ "johndoe@example.com" ]
			group = [ "db5c2352-b689-4135-babc-e97a8893128b" ]
			service_account = [ "test-sa" ]
		}
		where {
			allowed_ip = [ "0.0.0.0/0", "::/0" ]
			country = [ "NL", "CA", "US", "BR", "FR" ]
			country_not = [ "BE" ]
		}
		when {
			after = "2022-10-13T05:12:27Z"
			time_of_day_after = "00:00 UTC"
			time_of_day_before = "23:59 UTC"
		}
	}
}
`

func Test_DataSource_PolicyDocumentV2(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:        true,
		ProviderFactories: testProviderFactories(t, new(mocks.APIClientRequester)),
		Steps: []resource.TestStep{
			{
				Config: policyDocumentV2Config,
				Check: resource.ComposeTestCheckFunc(
					// permissions
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.database.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.database.0.max_session_duration_seconds", "3600"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.database.0.use_allowed_databases_list", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.database.0.allowed_databases.0.database", "test"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.database.0.allowed_databases.0.use_allowed_query_types_list", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.database.0.allowed_databases.0.allowed_query_types.0", "ReadOnly"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.database.0.allowed_databases.1.database", "aap"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.database.0.allowed_databases.1.use_allowed_query_types_list", "false"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.max_session_duration_seconds", "3600"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.use_allowed_usernames_list", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.allowed_usernames.0", "bas"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.shell.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.exec.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.exec.0.use_commands_list", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.exec.0.commands.0", "ls"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.sftp.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.tcp_forwarding.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.tcp_forwarding.0.use_allowed_connections_list", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.tcp_forwarding.0.allowed_connections.0.destination_address", "*"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.tcp_forwarding.0.allowed_connections.0.destination_port", "*"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.tcp_forwarding.0.allowed_connections.1.destination_address", "nu.nl"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.tcp_forwarding.0.allowed_connections.1.destination_port", "443"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.kubectl_exec.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.kubectl_exec.0.use_allowed_namespaces_list", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.kubectl_exec.0.allowed_namespaces.0.namespace", "test"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.kubectl_exec.0.allowed_namespaces.0.use_pod_selector", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.kubectl_exec.0.allowed_namespaces.0.pod_selector.aap2", "test"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.docker_exec.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.docker_exec.0.use_allowed_containers_list", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.ssh.0.docker_exec.0.allowed_containers.0", "test"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.vpn.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.http.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.tls.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.rdp.0.allowed", "true"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "permissions.0.vnc.0.allowed", "true"),

					// conditions
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.who.0.email.0", "johndoe@example.com"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.who.0.group.0", "db5c2352-b689-4135-babc-e97a8893128b"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.who.0.service_account.0", "test-sa"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.where.0.allowed_ip.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.where.0.allowed_ip.1", "::/0"),
					// country list item gets sorted
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.where.0.country.0", "BR"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.where.0.country.1", "CA"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.where.0.country.2", "FR"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.where.0.country.3", "NL"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.where.0.country.4", "US"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.where.0.country_not.0", "BE"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.when.0.after", "2022-10-13T05:12:27Z"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.when.0.time_of_day_after", "00:00 UTC"),
					resource.TestCheckResourceAttr("data.border0_policy_v2_document.unit_test", "condition.0.when.0.time_of_day_before", "23:59 UTC"),
					resource.TestCheckResourceAttrSet("data.border0_policy_v2_document.unit_test", "json"),
				),
			},
		},
	})
}
