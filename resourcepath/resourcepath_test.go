package resourcepath

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
		name     string
		path     string
		pattern  string
		elements map[string]string
		err      string
	}{
		{
			name:    "valid organization path",
			path:    "organizations/test-org",
			pattern: orgPattern,
			elements: map[string]string{
				"organization": "test-org",
			},
		},
		{
			name:    "valid",
			path:    "organizations/test-org/users/test-user",
			pattern: userPattern,
			elements: map[string]string{
				"user":         "test-user",
				"organization": "test-org",
			},
		},
		{
			name:    "valid singleton",
			path:    "organizations/test-org/logs",
			pattern: singleton,
			elements: map[string]string{
				"organization": "test-org",
			},
		},
		{
			name:    "valid wildcard",
			path:    "organizations/-/users/test-user",
			pattern: userPattern,
			elements: map[string]string{
				"organization": "-",
				"user":         "test-user",
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
			elements: map[string]string{
				"organization": "Test.org_1~x",
			},
		},
		{
			name:    "valid variable syntax in value",
			pattern: orgPattern,
			path:    "organizations/{test-org}",
			elements: map[string]string{
				"organization": "{test-org}",
			},
		},
		{
			name:    "valid space and percent in value",
			pattern: orgPattern,
			path:    "organizations/test org%20",
			elements: map[string]string{
				"organization": "test org%20",
			},
		},
		{
			name:    "valid non-ascii value",
			pattern: orgPattern,
			path:    "organizations/tést",
			elements: map[string]string{
				"organization": "tést",
			},
		},
		{
			name:    "repeated variable",
			pattern: "a/{x}/b/{x}",
			path:    "a/v/b/w",
			err:     "element {x}: repeated variable",
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
				want := &ResourcePath{pattern: tt.pattern, elements: tt.elements}
				assert.DeepEqual(t, path, want, cmp.AllowUnexported(ResourcePath{}))
				assert.Equal(t, path.Pattern(), tt.pattern)
				assert.Equal(t, path.String(), tt.path)
			}
		})
	}
}

func TestNewResourcePath(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name    string
		pattern string
		values  map[string]string
		want    string
		err     string
	}{
		{
			name:    "single variable",
			pattern: orgPattern,
			values:  map[string]string{"organization": "test-org"},
			want:    "organizations/test-org",
		},
		{
			name:    "multiple variables",
			pattern: userPattern,
			values: map[string]string{
				"organization": "test-org",
				"user":         "test-user",
			},
			want: "organizations/test-org/users/test-user",
		},
		{
			name:    "trailing literal",
			pattern: singleton,
			values:  map[string]string{"organization": "test-org"},
			want:    "organizations/test-org/logs",
		},
		{
			name:    "wildcard",
			pattern: userPattern,
			values: map[string]string{
				"organization": "-",
				"user":         "test-user",
			},
			want: "organizations/-/users/test-user",
		},
		{
			name:   "empty pattern",
			values: map[string]string{"organization": "test-org"},
			err:    "pattern can't be empty",
		},
		{
			name:    "missing value",
			pattern: userPattern,
			values:  map[string]string{"organization": "test-org"},
			err:     "element {user}: missing value",
		},
		{
			name:    "nil values",
			pattern: orgPattern,
			err:     "element {organization}: missing value",
		},
		{
			name:    "empty value",
			pattern: orgPattern,
			values:  map[string]string{"organization": ""},
			err:     "element {organization}: empty value",
		},
		{
			name:    "value with slash",
			pattern: orgPattern,
			values:  map[string]string{"organization": "test/org"},
			err:     `element {organization}: value "test/org" contains "/"`,
		},
		{
			name:    "repeated variable",
			pattern: "a/{x}/b/{x}",
			values:  map[string]string{"x": "v"},
			err:     "element {x}: repeated variable",
		},
		{
			name:    "no variables",
			pattern: "config",
			values:  map[string]string{},
			want:    "config",
		},
		{
			name:    "variables not in pattern",
			pattern: "shelves/{shelf}",
			values: map[string]string{
				"shelf": "top",
				"desks": "big",
				"chair": "comfy",
			},
			err: "element {chair}: not in pattern",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path, err := NewResourcePath(tt.pattern, tt.values)
			if tt.err != "" {
				assert.Error(t, err, tt.err)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, path.Pattern(), tt.pattern)
			assert.Equal(t, path.String(), tt.want)
			// Every path created with NewResourcePath can be parsed back.
			parsed, err := ParseString(path.String(), tt.pattern)
			assert.NilError(t, err)
			assert.DeepEqual(t, parsed, path, cmp.AllowUnexported(ResourcePath{}))
		})
	}
}

func TestNewResourcePath_CopiesValues(t *testing.T) {
	t.Parallel()
	values := map[string]string{"organization": "test-org"}
	path, err := NewResourcePath(orgPattern, values)
	assert.NilError(t, err)
	values["organization"] = "changed"
	assert.Equal(t, path.String(), "organizations/test-org")
}

func TestResourcePath_ZeroValue(t *testing.T) {
	t.Parallel()
	var path ResourcePath
	assert.Equal(t, path.Pattern(), "")
	assert.Equal(t, path.String(), "")
	assert.Equal(t, path.Get("organization"), "")
}
