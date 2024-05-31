package border0

import (
	"context"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderOption is a function that can be passed to `Provider()` to configures it.
type ProviderOption func(p *schema.Provider)

// Provider returns a Border0 implementation and definition of terraform `schema.Provider`.
func Provider(options ...ProviderOption) *schema.Provider {
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"border0_socket":                resourceSocket(),
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
			"border0_user_emails_to_ids": dataSourceUserEmailsToIDs(),
		},
	}

	for _, option := range options {
		option(provider)
	}

	return provider
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	apiURL := d.Get("api_url").(string)

	if token == "" {
		return nil, diag.Errorf("border0 provider credential is empty - set `token`")
	}

	return border0client.New(
		border0client.WithAuthToken(token),
		border0client.WithBaseURL(apiURL),
	), nil
}
