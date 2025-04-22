package validate

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestResourceID(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		id   string
		err  string
	}{
		{
			name: "minimal valid ID",
			id:   "a",
		},
		{
			name: "typical valid ID",
			id:   "valid-resource-id",
		},
		{
			name: "valid 63-char ID",
			id:   "a" + strings.Repeat("b", 61) + "c",
		},
		{
			name: "empty string",
			id:   "",
			err:  "resource ID can not be empty",
		},
		{
			name: "too long string",
			id:   strings.Repeat("a", 64),
			err:  "resource ID SHOULD not be longer than 63 characters",
		},
		{
			name: "starts with uppercase letter",
			id:   "Avalid",
			err:  "resource ID must start with a lowercase letter",
		},
		{
			name: "starts with digit",
			id:   "1invalid",
			err:  "resource ID must start with a lowercase letter",
		},
		{
			name: "starts with hyphen",
			id:   "-abc",
			err:  "resource ID must start with a lowercase letter",
		},
		{
			name: "ends with hyphen",
			id:   "valid-",
			err:  "resource ID cannot end with a hyphen",
		},
		{
			name: "invalid character: underscore",
			id:   "inva_lid",
			err:  "resource ID can only contain lowercase letters, numbers, and hyphens",
		},
		{
			name: "invalid uppercase in middle",
			id:   "valId",
			err:  "resource ID can only contain lowercase letters, numbers, and hyphens",
		},
		{
			name: "UUID formatted ID",
			id:   "cc6d5224-1b7d-4f0b-b3da-d6f6759398e4",
			err:  "resource ID can not ba a UUID",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ResourceID(tt.id)
			if tt.err != "" {
				assert.Error(t, err, tt.err)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}
