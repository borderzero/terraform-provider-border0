package schemautil

import (
	"github.com/borderzero/terraform-provider-border0/internal/diagnostics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetValues sets multiple values in a resource schema.
func SetValues(d *schema.ResourceData, values map[string]any) diag.Diagnostics {
	for key, value := range values {
		if err := d.Set(key, value); err != nil {
			return diagnostics.Error(err, "Failed to set %s", key)
		}
	}
	return nil
}
