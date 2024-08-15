package resourcepath

import (
	"testing"

	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

type UserPathType string

type UserPath struct {
	Groups       string `path:"groups/{group}/users/{user}"`
	Organization string `path:"organizations/{organization}/users/{user},organizations/{organization}/groups/{group}/users/{user}"`
	User         string `path:"users/{user}"`
}

func Test_ResourcePath(t *testing.T) {
	for _, tt := range []struct {
		name string
		path string
		resp UserPath
		err  string
	}{
		{
			name: "valid",
			path: "users/test-user",
			resp: UserPath{
				User: "test-user",
			},
		},
		{
			name: "valid",
			path: "organizations/test-org/users/test-user",
			resp: UserPath{
				User:         "test-user",
				Organization: "test-org",
			},
		},
		{
			name: "bad",
			path: "organizatios/test-org/users/test-user",
			err:  "no matching patterns for path: organizatios/test-org/users/test-user",
		},
		{
			name: "bad",
			path: "groups/test-org/users/test-user",
			err:  "parse resource path 'groups/test-org/users/test-user' with pattern 'groups/{group}/users/{user}': failed to unmarshal path",
		},
	} {
		var path UserPath
		err := UnmarshalString(tt.path, &path)
		if tt.err != "" {
			assert.Error(t, err, tt.err)
		} else {
			assert.NilError(t, err)
			assert.DeepEqual(t, path, tt.resp, protocmp.Transform())
		}
	}
}
