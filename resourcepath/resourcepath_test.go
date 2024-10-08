package resourcepath

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

const (
	orgPattern  = "organizations/{organization}"
	userPattern = "organizations/{organization}/users/{user}"
	singleton   = orgPattern + "/logs"
)

func Test_ResourcePath(t *testing.T) {
	for _, tt := range []struct {
		name    string
		path    string
		pattern string
		resp    *ResourcePath
		err     string
	}{
		{
			name:    "valid organization path",
			path:    "organizations/test-org",
			pattern: orgPattern,
			resp: &ResourcePath{
				segments: map[string]string{
					"organization": "test-org",
				},
			},
		},
		{
			name:    "valid",
			path:    "organizations/test-org/users/test-user",
			pattern: userPattern,
			resp: &ResourcePath{
				segments: map[string]string{
					"user":         "test-user",
					"organization": "test-org",
				},
			},
		},
		{
			name:    "valid singleton",
			path:    "organizations/test-org/logs",
			pattern: singleton,
			resp: &ResourcePath{
				segments: map[string]string{
					"organization": "test-org",
				},
			},
		},
		{
			name:    "invalida pattern",
			pattern: userPattern,
			path:    "organizations/test-org",
			err:     "segment users: unexpected EOF",
		},
		{
			name: "empty pattern",
			path: "organizations/test-org",
			err:  "pattern can't be empty",
		},
		{
			name:    "empty path",
			pattern: userPattern,
			err:     "segment organizations: unexpected EOF",
		},
		{
			name: "empty pattern",
			path: "organizations/test-org",
			err:  "pattern can't be empty",
		},
		{
			name: "empty pattern",
			path: "organizations/test-org",
			err:  "pattern can't be empty",
		},
	} {
		path, err := ParseString(tt.path, tt.pattern)
		if tt.err != "" {
			assert.Error(t, err, tt.err)
		} else {
			assert.NilError(t, err)
			assert.DeepEqual(t, path, tt.resp, protocmp.Transform(), cmp.AllowUnexported(ResourcePath{}))
		}
	}
}
