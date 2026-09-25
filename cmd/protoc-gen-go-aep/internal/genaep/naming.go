package genaep

import (
	"go/token"
	"go/types"
	"strings"

	"github.com/blaberg/aep-go/resourcepath"
)

// localIdentifiers are identifiers used in the body of generated functions,
// which parameters must not shadow.
var localIdentifiers = map[string]bool{
	"path":         true,
	"err":          true,
	"resourcepath": true,
}

// pascalCase converts a snake_case or kebab-case name to PascalCase,
// e.g. "user_event_id" -> "UserEventId" and "user-events" -> "UserEvents".
func pascalCase(s string) string {
	var b strings.Builder
	for _, part := range strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' }) {
		b.WriteString(strings.ToUpper(part[:1]))
		b.WriteString(part[1:])
	}
	return b.String()
}

// patternName returns the name of a pattern, used for its constant and constructor.
// It is built from the variables with the "_id" suffix removed, e.g.
// "authors/{author_id}/books/{book_id}" -> "AuthorBook". If the pattern ends with
// a literal, as singletons do, the literal is appended, e.g.
// "users/{user_id}/config" -> "UserConfig".
func patternName(pattern string) string {
	var b strings.Builder
	var last resourcepath.Element
	for e := range resourcepath.Elements(pattern) {
		last = e
		if e.IsVariable() {
			b.WriteString(pascalCase(strings.TrimSuffix(string(e.GetLiteral()), "_id")))
		}
	}
	if !last.IsVariable() {
		b.WriteString(pascalCase(string(last)))
	}
	return b.String()
}

// constName returns the name of the generated constant for a pattern, e.g. "AuthorBookPattern".
func constName(pattern string) string {
	return patternName(pattern) + "Pattern"
}

// constructorName returns the name of the generated constructor for a pattern,
// e.g. "NewAuthorBookResourcePath".
func constructorName(pattern string) string {
	return "New" + patternName(pattern) + "ResourcePath"
}

// getterName returns the name of the generated getter for a variable, e.g. "GetBookId".
func getterName(variable string) string {
	return "Get" + pascalCase(variable)
}

// paramName converts a variable name to a lowerCamelCase parameter name,
// e.g. "user_event_id" -> "userEventId". Names that would shadow a Go keyword,
// a predeclared identifier or an identifier used in the generated code get a
// trailing underscore, e.g. "type" -> "type_" and "error" -> "error_".
func paramName(variable string) string {
	name := pascalCase(variable)
	name = strings.ToLower(name[:1]) + name[1:]
	if token.IsKeyword(name) || types.Universe.Lookup(name) != nil || localIdentifiers[name] {
		return name + "_"
	}
	return name
}
