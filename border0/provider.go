package border0

import (
	"context"
	"fmt"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderOption func(p *schema.Provider)

func Provider(options ...ProviderOption) *schema.Provider {
	provider := &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BORDER0_AUTH_TOKEN", ""),
				Required:    true,
				Description: "The auth token used to authenticate with the Border0 API.",
				Sensitive:   true,
			},
			"base_url": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BORDER0_BASE_URL", "https://api.border0.com/api/v1"),
				Optional:    true,
				Description: "The URL of the Border0 API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"border0_socket": resourceSocket(),
		},
	}

	for _, option := range options {
		option(provider)
	}

	return provider
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	authToken := d.Get("auth_token").(string)
	baseURL := d.Get("base_url").(string)

	if authToken == "" {
		return nil, diag.Errorf("border0 provider credential is empty - set `auth_token`")
	}

	return border0client.New(
		border0client.WithAuthToken(authToken),
		border0client.WithBaseURL(baseURL),
	), nil
}

func setValues(d *schema.ResourceData, values map[string]any) diag.Diagnostics {
	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			return diagnosticsError(err, "Failed to set %s", key)
		}
	}
	return nil
}

func diagnosticsError(err error, message string, args ...interface{}) diag.Diagnostics {
	var detail string
	if err != nil {
		detail = err.Error()
	}

	diags := []diag.Diagnostic{
		{
			Severity: diag.Error,
			Summary:  fmt.Sprintf(message, args...),
			Detail:   detail,
		},
	}

	if clientError, ok := err.(border0client.Error); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  clientError.Error(),
		})
	}

	return diags
}
