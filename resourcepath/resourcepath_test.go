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
	t.Parallel()
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
				elements: map[string]string{
					"organization": "test-org",
				},
			},
		},
		{
			name:    "valid",
			path:    "organizations/test-org/users/test-user",
			pattern: userPattern,
			resp: &ResourcePath{
				elements: map[string]string{
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
				elements: map[string]string{
					"organization": "test-org",
				},
			},
		},
		{
			name:    "valid wildcard",
			path:    "organizations/-/users/test-user",
			pattern: userPattern,
			resp: &ResourcePath{
				elements: map[string]string{
					"organization": "-",
					"user":         "test-user",
				},
			},
		},
		{
			name:    "path shorter than pattern",
			pattern: userPattern,
			path:    "organizations/test-org",
			err:     "element users: unexpected EOF",
		},
		{
			name:    "singleton missing trailing literal",
			pattern: singleton,
			path:    "organizations/test-org",
			err:     "element logs: unexpected EOF",
		},
		{
			name:    "path longer than pattern",
			pattern: orgPattern,
			path:    "organizations/test-org/users/test-user",
			err:     "got trailing elements in path",
		},
		{
			name:    "trailing slash",
			pattern: orgPattern,
			path:    "organizations/test-org/",
			err:     "got trailing elements in path",
		},
		{
			name:    "leading slash",
			pattern: orgPattern,
			path:    "/organizations/test-org",
			err:     "element organizations: got ",
		},
		{
			name:    "wrong collection",
			pattern: orgPattern,
			path:    "orgs/test-org",
			err:     "element organizations: got orgs",
		},
		{
			name:    "collection is case sensitive",
			pattern: orgPattern,
			path:    "Organizations/test-org",
			err:     "element organizations: got Organizations",
		},
		{
			name: "empty pattern",
			path: "organizations/test-org",
			err:  "pattern can't be empty",
		},
		{
			name:    "empty path",
			pattern: userPattern,
			err:     "path can't be empty",
		},
		{
			name:    "missing variable",
			pattern: orgPattern,
			path:    "organizations/",
			err:     "element {organization}: empty value",
		},
		{
			name:    "missing variable in the middle",
			pattern: userPattern,
			path:    "organizations//users/test-user",
			err:     "element {organization}: empty value",
		},
		{
			name:    "valid unreserved characters",
			path:    "organizations/Test.org_1~x",
			pattern: orgPattern,
			resp: &ResourcePath{
				elements: map[string]string{
					"organization": "Test.org_1~x",
				},
			},
		},
		{
			name:    "valid variable syntax in value",
			pattern: orgPattern,
			path:    "organizations/{test-org}",
			resp: &ResourcePath{
				elements: map[string]string{
					"organization": "{test-org}",
				},
			},
		},
		{
			name:    "valid space and percent in value",
			pattern: orgPattern,
			path:    "organizations/test org%20",
			resp: &ResourcePath{
				elements: map[string]string{
					"organization": "test org%20",
				},
			},
		},
		{
			name:    "valid non-ascii value",
			pattern: orgPattern,
			path:    "organizations/tést",
			resp: &ResourcePath{
				elements: map[string]string{
					"organization": "tést",
				},
			},
		},
		{
			name:    "variable syntax in collection",
			pattern: orgPattern,
			path:    "{organizations}/test-org",
			err:     "element organizations: got {organizations}",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path, err := ParseString(tt.path, tt.pattern)
			if tt.err != "" {
				assert.Error(t, err, tt.err)
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, path, tt.resp, protocmp.Transform(), cmp.AllowUnexported(ResourcePath{}))
			}
		})
	}
}
