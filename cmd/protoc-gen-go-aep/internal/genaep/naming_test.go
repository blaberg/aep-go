package genaep

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestPascalCase(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		input string
		want  string
	}{
		{input: "book", want: "Book"},
		{input: "book_id", want: "BookId"},
		{input: "user_event_id", want: "UserEventId"},
		{input: "user-events", want: "UserEvents"},
		{input: "book2_id", want: "Book2Id"},
		{input: "book__id", want: "BookId"},
		{input: "", want: ""},
	} {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, pascalCase(tt.input), tt.want)
		})
	}
}

func TestParamName(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		input string
		want  string
	}{
		{input: "book", want: "book"},
		{input: "book_id", want: "bookId"},
		{input: "user_event_id", want: "userEventId"},
		{input: "id", want: "id"},
		// Go keywords.
		{input: "type", want: "type_"},
		{input: "func", want: "func_"},
		// Predeclared identifiers.
		{input: "string", want: "string_"},
		{input: "error", want: "error_"},
		{input: "nil", want: "nil_"},
		{input: "len", want: "len_"},
		// Identifiers used in the generated code.
		{input: "path", want: "path_"},
		{input: "err", want: "err_"},
		{input: "resourcepath", want: "resourcepath_"},
	} {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, paramName(tt.input), tt.want)
		})
	}
}

func TestPatternNames(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		pattern     string
		name        string
		constName   string
		constructor string
	}{
		{
			pattern:     "books/{book_id}",
			name:        "Book",
			constName:   "BookPattern",
			constructor: "NewBookResourcePath",
		},
		{
			pattern:     "authors/{author_id}/books/{book_id}",
			name:        "AuthorBook",
			constName:   "AuthorBookPattern",
			constructor: "NewAuthorBookResourcePath",
		},
		{
			pattern:     "users/{user_id}/user-events/{user_event_id}",
			name:        "UserUserEvent",
			constName:   "UserUserEventPattern",
			constructor: "NewUserUserEventResourcePath",
		},
		{
			pattern:     "books/{book}",
			name:        "Book",
			constName:   "BookPattern",
			constructor: "NewBookResourcePath",
		},
		{
			pattern:     "users/{user_id}/config",
			name:        "UserConfig",
			constName:   "UserConfigPattern",
			constructor: "NewUserConfigResourcePath",
		},
		{
			pattern:     "users/{user_id}/notification-settings",
			name:        "UserNotificationSettings",
			constName:   "UserNotificationSettingsPattern",
			constructor: "NewUserNotificationSettingsResourcePath",
		},
		{
			pattern:     "config",
			name:        "Config",
			constName:   "ConfigPattern",
			constructor: "NewConfigResourcePath",
		},
	} {
		t.Run(tt.pattern, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, patternName(tt.pattern), tt.name)
			assert.Equal(t, constName(tt.pattern), tt.constName)
			assert.Equal(t, constructorName(tt.pattern), tt.constructor)
		})
	}
}

func TestGetterName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, getterName("book_id"), "GetBookId")
	assert.Equal(t, getterName("user_event_id"), "GetUserEventId")
}
