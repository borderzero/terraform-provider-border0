package schemautil

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	idPartsDelimeter = ":"
)

// LoadMultipartID loads the schema resource id consisting of multiple parts onto the given variables.
func LoadMultipartID(d *schema.ResourceData, parts ...*string) diag.Diagnostics {
	idString := d.Id()

	// handle the empty id case
	if idString == "" {
		for i := range parts {
			*parts[i] = ""
		}
		return nil
	}

	idParts := strings.Split(idString, idPartsDelimeter)
	if len(idParts) != len(parts) {
		return diag.Errorf("the number of parts on the schema resource id (%d) did not match the number of parts expected (%d)", len(idParts), len(parts))
	}
	for i, idPart := range idParts {
		*parts[i] = idPart
	}
	return nil
}

// SetMultipartID sets the id on the schema resource data concatenating multiple parts.
func SetMultipartID(d *schema.ResourceData, parts ...string) {
	d.SetId(strings.Join(parts, idPartsDelimeter))
}
