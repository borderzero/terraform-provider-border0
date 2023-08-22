package schemautil

import (
	"fmt"

	"github.com/borderzero/border0-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetValues sets multiple values in a resource schema.
func SetValues(d *schema.ResourceData, values map[string]any) diag.Diagnostics {
	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			return DiagnosticsError(err, "Failed to set %s", key)
		}
	}
	return nil
}

// DiagnosticsError returns a diag.Diagnostics with an error and a formatted message.
func DiagnosticsError(err error, message string, args ...interface{}) diag.Diagnostics {
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
