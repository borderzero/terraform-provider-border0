package border0

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sshBlock builds the []any that the SDK hands to parseSSHPermissions for a single
// `ssh { ... }` block, with nested blocks represented as sets of one element.
func sshBlock(nested map[string]map[string]any) []any {
	perm := map[string]any{"allowed": true}
	for name, attrs := range nested {
		perm[name] = schema.NewSet(func(any) int { return 1 }, []any{attrs})
	}
	return []any{perm}
}

func Test_parseSSHPermissions_AgentForwarding(t *testing.T) {
	tests := []struct {
		name    string
		block   []any
		wantKey bool
	}{
		{
			name:    "enabled",
			block:   sshBlock(map[string]map[string]any{"sftp": {"allowed": true}, "agent_forwarding": {"allowed": true}}),
			wantKey: true,
		},
		{
			name:    "allowed = false",
			block:   sshBlock(map[string]map[string]any{"sftp": {"allowed": true}, "agent_forwarding": {"allowed": false}}),
			wantKey: false,
		},
		{
			name:    "block omitted",
			block:   sshBlock(map[string]map[string]any{"sftp": {"allowed": true}}),
			wantKey: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms := parseSSHPermissions(tt.block)
			require.NotNil(t, perms, "ssh itself is enabled in every case")
			assert.NotNil(t, perms.SFTP)

			raw, err := json.Marshal(perms)
			require.NoError(t, err)
			var got map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(raw, &got))

			if tt.wantKey {
				assert.NotNil(t, perms.AgentForwarding)
				assert.JSONEq(t, `{}`, string(got["agent_forwarding"]))
			} else {
				assert.Nil(t, perms.AgentForwarding)
				assert.NotContains(t, got, "agent_forwarding")
			}
			assert.JSONEq(t, `{}`, string(got["sftp"]))
		})
	}
}
