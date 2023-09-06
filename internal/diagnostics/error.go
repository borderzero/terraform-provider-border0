package diagnostics

import (
	"fmt"

	"github.com/borderzero/border0-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Error returns a diag.Diagnostics with an error and a formatted message.
func Error(err error, message string, args ...interface{}) diag.Diagnostics {
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

	if clientError, ok := err.(client.Error); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  clientError.Error(),
		})
	}

	return diags
}
