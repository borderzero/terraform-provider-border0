package border0

import (
	"context"
	"time"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/internal/lib/sem"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// NOTE: currently only applies to socket, but we plan on
	// having this apply to all resources using the same semaphore.
	//
	// Terraform's parallelism is 10 by default but can be set to any
	// value using the "-parallelism" flag e.g. -parallelism=200...
	// So we cap it at 100 here in case it's set to a higher value.
	maxParallelism = 100

	defaultTimeout = time.Second * 30
)

// ProviderOption is a function that can be passed to `Provider()` to configures it.
type ProviderOption func(p *schema.Provider)

// Provider returns a Border0 implementation and definition of terraform `schema.Provider`.
func Provider(options ...ProviderOption) *schema.Provider {
	semaphore := sem.New(maxParallelism)

	provider := &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BORDER0_TOKEN", ""),
				Required:    true,
				Description: "The auth token used to authenticate with the Border0 API. Can also be set with the `BORDER0_TOKEN` environment variable. If you need to generate a Border0 access token, go to [Border0 Admin Portal](https://portal.border0.com) -> Organization Settings -> Access Tokens, create a token in `Member` permission groups.",
				Sensitive:   true,
			},
			"api_url": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BORDER0_API", "https://api.border0.com/api/v1"),
				Optional:    true,
				Description: "The URL of the Border0 API. Can also be set with the `BORDER0_API` environment variable. Defaults to `https://api.border0.com/api/v1`.",
			},
			"http_client_timeout": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BORDER0_HTTP_CLIENT_TIMEOUT", "30s"),
				Optional:    true,
				Description: "The timeout for each HTTP request. Can also be set with the `BORDER0_HTTP_CLIENT_TIMEOUT` environment variable. Defaults to `30s`.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"border0_socket":                resourceSocket(semaphore),
			"border0_policy":                resourcePolicy(),
			"border0_policy_attachment":     resourcePolicyAttachment(),
			"border0_connector":             resourceConnector(),
			"border0_connector_token":       resourceConnectorToken(),
			"border0_user":                  resourceUser(),
			"border0_group":                 resourceGroup(),
			"border0_service_account":       resourceServiceAccount(),
			"border0_service_account_token": resourceServiceAccountToken(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"border0_policy_document":    dataSourcePolicyDocument(),
			"border0_policy_v2_document": dataSourcePolicyV2Document(),
			"border0_user_emails_to_ids": dataSourceUserEmailsToIDs(),
			"border0_group_names_to_ids": dataSourceGroupNamesToIDs(),
		},
	}

	for _, option := range options {
		option(provider)
	}

	return provider
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	token := d.Get("token").(string)
	if token == "" {
		return nil, diag.Errorf("border0 provider credential is empty - set `token`")
	}

	opts := []border0client.Option{border0client.WithAuthToken(token)}

	if apiURLAny := d.Get("api_url"); apiURLAny != nil {
		apiURL, ok := apiURLAny.(string)
		if !ok {
			return nil, diag.Errorf("`api_url` is set but is not a string")
		}
		opts = append(opts, border0client.WithBaseURL(apiURL))
	}

	if timeoutAny := d.Get("http_client_timeout"); timeoutAny != nil {
		timeoutStr, ok := timeoutAny.(string)
		if !ok {
			return nil, diag.Errorf("`http_client_timeout` is set but is not a string")
		}
		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, diag.Errorf("`http_client_timeout` is set but is not a valid duration e.g. 10s: %v", err)
		}
		opts = append(opts, border0client.WithTimeout(timeout))
	} else {
		opts = append(opts, border0client.WithTimeout(defaultTimeout))
	}

	return border0client.New(opts...), nil
}
