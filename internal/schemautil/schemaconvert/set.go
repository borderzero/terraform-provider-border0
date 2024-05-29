package schemaconvert

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// SetToSlice converts a schema.Set to a slice of the set's type.
func SetToSlice[T any](s *schema.Set) []T {
	slice := []T{}
	for _, elem := range s.List() {
		slice = append(slice, elem.(T))
	}
	return slice
}
