package border0

import (
	"context"
	"encoding/json"
	"strconv"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicyV2Document() *schema.Resource {
	return &schema.Resource{
		Description: "`border0_policy_v2_document` data source can be used to generate a policy document in JSON format for use with `border0_policy` resource.",
		ReadContext: dataSourcePolicyV2DocumentRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The permissions that you want to allow.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The Database permissions that you want to allow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"use_allowed_databases_list": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Use allowed databases list.",
									},
									"allowed_databases": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of allowed databases.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"database": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name of the database.",
												},
												"allowed_query_types": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of allowed query types.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"use_allowed_query_types_list": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Use allowed query types list.",
												},
											},
										},
									},
									"max_session_duration_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum session duration in seconds.",
									},
									"allowed": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether database access is allowed.",
									},
								},
							},
						},
						"ssh": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The SSH permissions that you want to allow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether ssh access is allowed.",
									},
									"shell": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "SSH Shell permission.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"allowed": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether ssh shell access is allowed.",
												},
											},
										},
									},
									"exec": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "SSH Exec permission.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"use_commands_list": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Use allowed commands list.",
												},
												"commands": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of allowed commands.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"allowed": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether ssh exec access is allowed.",
												},
											},
										},
									},
									"sftp": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "SSH SFTP permission.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"allowed": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether ssh sftp access is allowed.",
												},
											},
										},
									},
									"tcp_forwarding": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "SSH TCP Forwarding permission.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"use_allowed_connections_list": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Use allowed connections list.",
												},
												"allowed_connections": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of allowed TCP forwarding connections.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"destination_address": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Destination address.",
															},
															"destination_port": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Destination port.",
															},
														},
													},
												},
												"allowed": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether ssh sftp access is allowed.",
												},
											},
										},
									},
									"kubectl_exec": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "SSH Kubectl Exec permission.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"use_allowed_namespaces_list": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Use allowed namespaces list.",
												},
												"allowed_namespaces": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of allowed namespaces.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"namespace": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Namespace name.",
															},
															"pod_selector": {
																Type:        schema.TypeMap,
																Optional:    true,
																Description: "Pod selector map.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"use_pod_selector": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Use pod selector.",
															},
														},
													},
												},
												"allowed": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether kubernetes exec access is allowed.",
												},
											},
										},
									},
									"docker_exec": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "SSH Docker Exec permission.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"use_allowed_containers_list": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Use allowed containers list.",
												},
												"allowed_containers": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of allowed containers.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"allowed": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether docker exec access is allowed.",
												},
											},
										},
									},
									"max_session_duration_seconds": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum session duration in seconds.",
									},
									"use_allowed_usernames_list": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Use allowed usernames list.",
									},
									"allowed_usernames": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of allowed usernames.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"http": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The HTTP permissions that you want to allow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether http access is allowed.",
									},
								},
							},
						},
						"kubernetes": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The Kubernetes permissions that you want to allow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether kubernetes access is allowed.",
									},
								},
							},
						},
						"tls": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The TLS permissions that you want to allow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether tls access is allowed.",
									},
								},
							},
						},
						"vnc": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The VNC permissions that you want to allow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether vnc access is allowed.",
									},
								},
							},
						},
						"rdp": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The RDP permissions that you want to allow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether rdp access is allowed.",
									},
								},
							},
						},
						"network": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The Network permissions that you want to allow.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether network access is allowed.",
									},
								},
							},
						},
					},
				},
			},
			"condition": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"who": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"email": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The email address of the user who is allowed to perform the actions.",
									},
									"group": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The group uuid of the group which is allowed to perform the actions.",
									},
									"service_account": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The service account name which is allowed to perform the actions.",
									},
								},
							},
							Description: "Who is allowed to perform the actions.",
						},
						"where": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_ip": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The IP address that the request must originate from to be allowed to perform the actions.",
									},
									"country": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The country that the request must originate from to be allowed to perform the actions.",
									},
									"country_not": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
										Description: "The country that the request must _NOT_ originate from to be allowed to perform the actions.",
									},
								},
							},
							Description: "Where the request must originate from to be allowed to perform the actions.",
						},
						"when": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"after": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the request must be made after to be allowed to perform the actions.",
									},
									"before": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the request must be made before to be allowed to perform the actions.",
									},
									"time_of_day_after": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the request must be made after to be allowed to perform the actions.",
									},
									"time_of_day_before": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the request must be made before to be allowed to perform the actions.",
									},
								},
							},
							Description: "When the request must be made to be allowed to perform the actions.",
						},
					},
				},
				Description: "The conditions under which you want to allow the actions.",
			},
		},
	}
}

func dataSourcePolicyV2DocumentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var policyData border0client.PolicyDataV2

	if v, ok := d.GetOk("permissions"); ok {
		if permissions := v.(*schema.Set).List(); len(permissions) > 0 {
			permMap := permissions[0].(map[string]interface{})

			if v, ok := permMap["database"]; ok {
				policyData.Permissions.Database = parseDatabasePermissions(v.(*schema.Set).List())
			}
			if v, ok := permMap["ssh"]; ok {
				policyData.Permissions.SSH = parseSSHPermissions(v.(*schema.Set).List())
			}
			if v, ok := permMap["http"]; ok {
				if httpPerms := v.(*schema.Set).List(); len(httpPerms) > 0 {
					httpPerm := httpPerms[0].(map[string]interface{})
					if v, ok := httpPerm["allowed"]; ok {
						if allowed, ok := v.(bool); ok && allowed {
							policyData.Permissions.HTTP = &border0client.HTTPPermissions{}
						}
					}
				}
			}
			if v, ok := permMap["kubernetes"]; ok {
				if kubernetesPerms := v.(*schema.Set).List(); len(kubernetesPerms) > 0 {
					kubernetesPerm := kubernetesPerms[0].(map[string]interface{})
					if v, ok := kubernetesPerm["allowed"]; ok {
						if allowed, ok := v.(bool); ok && allowed {
							policyData.Permissions.Kubernetes = &border0client.KubernetesPermissions{}
						}
					}
				}
			}
			if v, ok := permMap["tls"]; ok {
				if tlsPerms := v.(*schema.Set).List(); len(tlsPerms) > 0 {
					tlsPerm := tlsPerms[0].(map[string]interface{})
					if v, ok := tlsPerm["allowed"]; ok {
						if allowed, ok := v.(bool); ok && allowed {
							policyData.Permissions.TLS = &border0client.TLSPermissions{}
						}
					}
				}
			}
			if v, ok := permMap["vnc"]; ok {
				if vncPerms := v.(*schema.Set).List(); len(vncPerms) > 0 {
					vncPerm := vncPerms[0].(map[string]interface{})
					if v, ok := vncPerm["allowed"]; ok {
						if allowed, ok := v.(bool); ok && allowed {
							policyData.Permissions.VNC = &border0client.VNCPermissions{}
						}
					}
				}
			}
			if v, ok := permMap["rdp"]; ok {
				if rdpPerms := v.(*schema.Set).List(); len(rdpPerms) > 0 {
					rdpPerm := rdpPerms[0].(map[string]interface{})
					if v, ok := rdpPerm["allowed"]; ok {
						if allowed, ok := v.(bool); ok && allowed {
							policyData.Permissions.RDP = &border0client.RDPPermissions{}
						}
					}
				}
			}
			if v, ok := permMap["network"]; ok {
				if networkPerms := v.(*schema.Set).List(); len(networkPerms) > 0 {
					networkPerm := networkPerms[0].(map[string]interface{})
					if v, ok := networkPerm["allowed"]; ok {
						if allowed, ok := v.(bool); ok && allowed {
							policyData.Permissions.Network = &border0client.NetworkPermissions{}
						}
					}
				}
			}
		}
	}

	if v, ok := d.GetOk("condition"); ok {
		if conditions := v.(*schema.Set).List(); len(conditions) > 0 {
			condition := conditions[0].(map[string]interface{})
			if v, ok := condition["who"]; ok {
				if whos := v.(*schema.Set).List(); len(whos) > 0 {
					who := whos[0].(map[string]interface{})
					if v, ok := who["email"]; ok {
						policyData.Condition.Who.Email = policyDecodeStringList(v.(*schema.Set).List())
					}
					if v, ok := who["group"]; ok {
						policyData.Condition.Who.Group = policyDecodeStringList(v.(*schema.Set).List())
					}
					if v, ok := who["service_account"]; ok {
						policyData.Condition.Who.ServiceAccount = policyDecodeStringList(v.(*schema.Set).List())
					}
				}
			}
			if v, ok := condition["where"]; ok {
				if wheres := v.(*schema.Set).List(); len(wheres) > 0 {
					where := wheres[0].(map[string]interface{})
					if v, ok := where["allowed_ip"]; ok {
						policyData.Condition.Where.AllowedIP = policyDecodeStringList(v.(*schema.Set).List())
					}
					if v, ok := where["country"]; ok {
						policyData.Condition.Where.Country = policyDecodeStringList(v.(*schema.Set).List())
					}
					if v, ok := where["country_not"]; ok {
						policyData.Condition.Where.CountryNot = policyDecodeStringList(v.(*schema.Set).List())
					}
				}
			}
			if v, ok := condition["when"]; ok {
				if whens := v.(*schema.Set).List(); len(whens) > 0 {
					when := whens[0].(map[string]interface{})
					if v, ok := when["after"]; ok {
						policyData.Condition.When.After = v.(string)
					}
					if v, ok := when["before"]; ok {
						policyData.Condition.When.Before = v.(string)
					}
					if v, ok := when["time_of_day_after"]; ok {
						policyData.Condition.When.TimeOfDayAfter = v.(string)
					}
					if v, ok := when["time_of_day_before"]; ok {
						policyData.Condition.When.TimeOfDayBefore = v.(string)
					}
				}
			}
		}
	}

	jsonPolicyData, err := json.MarshalIndent(policyData, "", "  ")
	if err != nil {
		return diagnostics.Error(err, "Failed to marshal policy data")
	}
	jsonString := string(jsonPolicyData)

	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(stringHashcode(jsonString)))
	return nil
}

func parseDatabasePermissions(dbPerms []interface{}) *border0client.DatabasePermissions {
	var maxSessionDurationSeconds *int
	var allowedDatabases *[]border0client.DatabasePermission
	var allowed bool

	for _, dbPerm := range dbPerms {
		permMap := dbPerm.(map[string]interface{})

		if v, ok := permMap["allowed"]; ok {
			allowed = v.(bool)
		}

		if v, ok := permMap["use_allowed_databases_list"]; ok {
			if useAllowedDatabasesList, ok := v.(bool); ok && useAllowedDatabasesList {
				if v, ok := permMap["allowed_databases"]; ok {
					databases := []border0client.DatabasePermission{}
					for _, ad := range v.([]interface{}) {
						adMap := ad.(map[string]interface{})
						allowedDatabase := border0client.DatabasePermission{
							Database: adMap["database"].(string),
						}

						if aql, ok := adMap["use_allowed_query_types_list"]; ok {
							if useAllowedQueryTypesList, ok := aql.(bool); ok && useAllowedQueryTypesList {
								if aq, ok := adMap["allowed_query_types"]; ok {
									queryTypes := []string{}
									for _, qt := range aq.([]interface{}) {
										queryTypes = append(queryTypes, qt.(string))
									}
									allowedDatabase.AllowedQueryTypes = &queryTypes
								}
							}
						}

						databases = append(databases, allowedDatabase)
					}
					allowedDatabases = &databases
				}
			}
		}

		if v, ok := permMap["max_session_duration_seconds"]; ok {
			duration := v.(int)
			if duration > 0 {
				maxSessionDurationSeconds = &duration
			}
		}
	}

	if !allowed {
		return nil
	}

	return &border0client.DatabasePermissions{
		AllowedDatabases:          allowedDatabases,
		MaxSessionDurationSeconds: maxSessionDurationSeconds,
	}
}

func parseSSHPermissions(sshPerms []interface{}) *border0client.SSHPermissions {
	var shellPermission *border0client.SSHShellPermission
	var execPermission *border0client.SSHExecPermission
	var sftpPermission *border0client.SSHSFTPPermission
	var tcpForwardingPermission *border0client.SSHTCPForwardingPermission
	var kubectlExecPermission *border0client.SSHKubectlExecPermission
	var dockerExecPermission *border0client.SSHDockerExecPermission
	var maxSessionDurationSeconds *int
	var allowedUsernames *[]string
	var allowed bool

	for _, sshPerm := range sshPerms {
		permMap := sshPerm.(map[string]interface{})
		if v, ok := permMap["allowed"]; ok {
			allowed = v.(bool)
		}

		if v, ok := permMap["shell"]; ok {
			var execAllowed bool
			if shells := v.(*schema.Set).List(); len(shells) > 0 {
				shell := shells[0].(map[string]interface{})
				if v, ok := shell["allowed"]; ok {
					execAllowed, _ = v.(bool)
				}
			}

			if execAllowed {
				shellPermission = &border0client.SSHShellPermission{}
			}
		}

		if v, ok := permMap["exec"]; ok {
			var useCommandList, execAllowed bool
			var commands []string

			if execs := v.(*schema.Set).List(); len(execs) > 0 {
				exec := execs[0].(map[string]interface{})
				if v, ok := exec["allowed"]; ok {
					execAllowed, _ = v.(bool)
				}

				if v, ok := exec["use_commands_list"]; ok {
					useCommandList, _ = v.(bool)
				}

				if v, ok := exec["commands"]; ok {
					for _, cmd := range v.([]interface{}) {
						commands = append(commands, cmd.(string))
					}
				}
			}

			if execAllowed {
				execPermission = &border0client.SSHExecPermission{}

				if useCommandList {
					execPermission.Commands = &commands
				}
			}
		}

		if v, ok := permMap["sftp"]; ok {
			var sftpAllowed bool
			if sftps := v.(*schema.Set).List(); len(sftps) > 0 {
				sftp := sftps[0].(map[string]interface{})
				if v, ok := sftp["allowed"]; ok {
					sftpAllowed, _ = v.(bool)
				}
			}

			if sftpAllowed {
				sftpPermission = &border0client.SSHSFTPPermission{}
			}
		}

		if v, ok := permMap["tcp_forwarding"]; ok {
			var tcpForwardingAllowed, useAllowedConnectionsList bool
			var allowedConnections *[]border0client.SSHTcpForwardingConnection
			if tcpForwardings := v.(*schema.Set).List(); len(tcpForwardings) > 0 {
				tcpForwarding := tcpForwardings[0].(map[string]interface{})
				if v, ok := tcpForwarding["allowed"]; ok {
					tcpForwardingAllowed, _ = v.(bool)
				}

				if v, ok := tcpForwarding["use_allowed_connections_list"]; ok {
					useAllowedConnectionsList, _ = v.(bool)
				}

				if v, ok := tcpForwarding["allowed_connections"]; ok {
					allowedConnections = parseSSHTCPForwardingConnections(v.([]interface{}))
				}

				if tcpForwardingAllowed {
					tcpForwardingPermission = &border0client.SSHTCPForwardingPermission{}

					if useAllowedConnectionsList {
						tcpForwardingPermission.AllowedConnections = allowedConnections
					}
				}
			}
		}

		if v, ok := permMap["kubectl_exec"]; ok {
			var kubectlExecAllowed, useAllowedNamespacesList bool
			var allowedNamespaces *[]border0client.KubectlExecNamespace
			if kubectlExecs := v.(*schema.Set).List(); len(kubectlExecs) > 0 {
				kubectlExec := kubectlExecs[0].(map[string]interface{})
				if v, ok := kubectlExec["allowed"]; ok {
					kubectlExecAllowed, _ = v.(bool)
				}

				if v, ok := kubectlExec["use_allowed_namespaces_list"]; ok {
					useAllowedNamespacesList, _ = v.(bool)
				}

				if v, ok := kubectlExec["allowed_namespaces"]; ok {
					allowedNamespaces = parseSSHKubectlExecNamespaces(v.([]interface{}))
				}

				if kubectlExecAllowed {
					kubectlExecPermission = &border0client.SSHKubectlExecPermission{}

					if useAllowedNamespacesList {
						kubectlExecPermission.AllowedNamespaces = allowedNamespaces
					}
				}
			}
		}

		if v, ok := permMap["docker_exec"]; ok {
			var dockerExecAllowed, useAllowedContainerList bool
			var allowedContainers []string
			if dockerExecs := v.(*schema.Set).List(); len(dockerExecs) > 0 {
				dockerExec := dockerExecs[0].(map[string]interface{})
				if v, ok := dockerExec["allowed"]; ok {
					dockerExecAllowed, _ = v.(bool)
				}

				if v, ok := dockerExec["use_allowed_containers_list"]; ok {
					useAllowedContainerList, _ = v.(bool)
				}

				if v, ok := dockerExec["allowed_containers"]; ok {
					for _, cmd := range v.([]interface{}) {
						allowedContainers = append(allowedContainers, cmd.(string))
					}
				}

				if dockerExecAllowed {
					dockerExecPermission = &border0client.SSHDockerExecPermission{}

					if useAllowedContainerList {
						dockerExecPermission.AllowedContainers = &allowedContainers
					}
				}
			}
		}

		if v, ok := permMap["max_session_duration_seconds"]; ok {
			duration := v.(int)
			if duration > 0 {
				maxSessionDurationSeconds = &duration
			}
		}

		if v, ok := permMap["use_allowed_usernames_list"]; ok {
			if useAllowedUsernamesList, ok := v.(bool); ok && useAllowedUsernamesList {
				if v, ok := permMap["allowed_usernames"]; ok {
					usernames := []string{}
					for _, username := range v.([]interface{}) {
						usernames = append(usernames, username.(string))
					}
					allowedUsernames = &usernames
				}
			}
		}
	}

	if !allowed {
		return nil
	}

	return &border0client.SSHPermissions{
		Shell:                     shellPermission,
		Exec:                      execPermission,
		SFTP:                      sftpPermission,
		TCPForwarding:             tcpForwardingPermission,
		KubectlExec:               kubectlExecPermission,
		DockerExec:                dockerExecPermission,
		MaxSessionDurationSeconds: maxSessionDurationSeconds,
		AllowedUsernames:          allowedUsernames,
	}
}

func parseSSHTCPForwardingConnections(allowedConnections []interface{}) *[]border0client.SSHTcpForwardingConnection {
	var connections []border0client.SSHTcpForwardingConnection

	for _, conn := range allowedConnections {
		connMap := conn.(map[string]interface{})

		var destAddress, destPort string
		if addr, ok := connMap["destination_address"]; ok && addr != nil {
			destAddress = addr.(string)
		}
		if port, ok := connMap["destination_port"]; ok && port != nil {
			destPort = port.(string)
		}
		connection := border0client.SSHTcpForwardingConnection{
			DestinationAddress: &destAddress,
			DestinationPort:    &destPort,
		}

		connections = append(connections, connection)
	}

	return &connections
}

func parseSSHKubectlExecNamespaces(allowedNamespaces []interface{}) *[]border0client.KubectlExecNamespace {
	var namespaces []border0client.KubectlExecNamespace

	for _, ns := range allowedNamespaces {
		nsMap := ns.(map[string]interface{})
		namespace := border0client.KubectlExecNamespace{
			Namespace: nsMap["namespace"].(string),
		}

		if ps, ok := nsMap["use_pod_selector"]; ok {
			if usePodSelector, ok := ps.(bool); ok && usePodSelector {
				if ps, ok := nsMap["pod_selector"]; ok {
					podSelector := map[string]string{}
					for key, value := range ps.(map[string]interface{}) {
						podSelector[key] = value.(string)
					}
					namespace.PodSelector = &podSelector
				}
			}
		}

		namespaces = append(namespaces, namespace)
	}

	return &namespaces
}
