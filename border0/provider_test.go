package border0_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	border0client "github.com/borderzero/border0-go/client"
	"github.com/borderzero/terraform-provider-border0/border0"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"
)

func Test_Provider(t *testing.T) {
	if err := border0.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func Test_Provider_Implemented(t *testing.T) {
	var _ *schema.Provider = border0.Provider()
}

var matchContext = mock.MatchedBy(func(ctx context.Context) bool {
	return true
})

func mockCallsInOrder(calls ...*mock.Call) {
	for i := len(calls) - 1; i >= 0; i-- {
		calls[i].Once()
		if i > 0 {
			calls[i].NotBefore(calls[i-1])
		}
	}
}

func testProviderFactories(t *testing.T, api border0client.Requester) map[string]func() (*schema.Provider, error) {
	t.Helper()

	return map[string]func() (*schema.Provider, error){
		"border0": func() (*schema.Provider, error) {
			return border0.Provider(func(p *schema.Provider) {
				p.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
					return api, nil
				}
				p.Schema = nil // no need to include any of the global configuration
			}), nil
		},
	}
}

func testMatchResourceAttrJSON(name, key, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		actual, ok := rs.Primary.Attributes[key]
		if !ok {
			return fmt.Errorf("field not found: %s", key)
		}

		var expectedJSONAsAny, actualJSONAsAny any
		if err := json.Unmarshal([]byte(expected), &expectedJSONAsAny); err != nil {
			return fmt.Errorf("expected ('%s') needs to be valid json. JSON parsing error: %s", key, err.Error())
		}
		if err := json.Unmarshal([]byte(actual), &actualJSONAsAny); err != nil {
			return fmt.Errorf("input ('%s') needs to be valid json. JSON parsing error: %s", key, err.Error())
		}

		if !reflect.DeepEqual(expectedJSONAsAny, actualJSONAsAny) {
			return fmt.Errorf("field ('%s') expected and actual JSON values are not equal", key)
		}

		return nil
	}
}
