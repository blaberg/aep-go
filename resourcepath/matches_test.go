package resourcepath

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMatches(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name     string
		path     string
		pattern  string
		expected bool
	}{
		{
			name:     "valid pattern",
			path:     "shippers/1/sites/1",
			pattern:  "shippers/{shipper}/sites/{site}",
			expected: true,
		},

		{
			name:     "path longer than pattern",
			path:     "shippers/1/sites/1/settings",
			pattern:  "shippers/{shipper}/sites/{site}",
			expected: false,
		},

		{
			name:     "empty pattern",
			pattern:  "",
			path:     "shippers/1/sites/1",
			expected: false,
		},

		{
			name:     "empty pattern and empty path",
			pattern:  "",
			path:     "",
			expected: false,
		},

		{
			name:     "singleton",
			path:     "shippers/1/sites/1/settings",
			pattern:  "shippers/{shipper}/sites/{site}/settings",
			expected: true,
		},

		{
			name:     "wildcard pattern",
			path:     "shippers/1/sites/1",
			pattern:  "shippers/-/sites/-",
			expected: false,
		},

		{
			name:     "full parent",
			path:     "//freight-example.einride.tech/shippers/1/sites/1",
			pattern:  "shippers/{shipper}/sites/{site}",
			expected: true,
		},

		{
			name:     "full pattern",
			path:     "shippers/1",
			pattern:  "//freight-example.einride.tech/shippers/{shipper}",
			expected: false,
		},

		{
			name:     "slash prefix in the path",
			path:     "/shippers/1",
			pattern:  "shippers/{shipper}",
			expected: true,
		},

		{
			name:     "slash prefix in the pattern",
			path:     "shippers/1",
			pattern:  "/shippers/{shipper}",
			expected: true,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Check(
				t,
				Match(tt.pattern, tt.path) == tt.expected,
				fmt.Sprintf("expected Match(%q, %q)=%t", tt.pattern, tt.path, tt.expected),
			)
		})
	}
}
