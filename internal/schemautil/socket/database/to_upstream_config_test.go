package database

import (
	"testing"

	"github.com/borderzero/border0-go/types/service"
)

// SDKv2 hands every schema key to the converter, unset ones as zero values,
// so the IAM branch must not touch the (nil) TlsAuth config.
func TestGcpCloudSqlToUpstreamConfigIamAuth(t *testing.T) {
	data := map[string]any{
		"tls_auth":                   false,
		"cloudsql_connector_enabled": false,
		"cloudsql_iam_auth":          true,
		"username":                   "pam@project.iam",
		"password":                   "",
		"certificate":                "",
		"private_key":                "",
		"ca_certificate":             "",
		"cloudsql_instance_id":       "project:us-central1:instance",
		"gcp_credentials":            "{}",
	}

	config := new(service.GcpCloudSqlDatabaseServiceConfiguration)
	if diags := gcpCloudSqlToUpstreamConfig(data, config); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	iam := config.GcpCloudSQLConnectorIAMAuth
	if iam == nil {
		t.Fatal("expected IAM auth configuration to be set")
	}
	if iam.Username != "pam@project.iam" || iam.InstanceId != "project:us-central1:instance" || iam.GcpCredentialsJson != "{}" {
		t.Fatalf("unexpected IAM auth configuration: %+v", iam)
	}
	if config.TlsAuth != nil {
		t.Fatalf("expected no TLS auth configuration, got %+v", config.TlsAuth)
	}
}
