package genaep

import (
	"fmt"
	"regexp"

	"github.com/blaberg/aep-go/resourcepath"
)

// variableRegexp matches pattern variables that can be turned into Go identifiers, e.g. "book_id".
var variableRegexp = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// resource is the model of a resource message that code is generated from.
type resource struct {
	// Name is the Go name of the message, e.g. "Book".
	Name string
	// Patterns are the resource patterns, in the order they are declared.
	Patterns []pattern
	// Variables is the union of the variables in all patterns, in order of first appearance.
	Variables []variable
}

type pattern struct {
	// Value is the pattern, e.g. "authors/{author_id}/books/{book_id}".
	Value string
	// Const is the name of the generated constant, e.g. "AuthorBookPattern".
	Const string
	// Constructor is the name of the generated constructor, e.g. "NewAuthorBookResourcePath".
	Constructor string
	// Variables are the variables in the pattern, in order.
	Variables []variable
}

type variable struct {
	// Name is the variable in the pattern, e.g. "author_id".
	Name string
	// Param is the name of the constructor parameter, e.g. "authorId".
	Param string
	// Getter is the name of the generated getter, e.g. "GetAuthorId".
	Getter string
}

func (r *resource) pathType() string    { return r.Name + "ResourcePath" }
func (r *resource) patternType() string { return r.Name + "ResourcePattern" }
func (r *resource) parseFunc() string   { return "Parse" + r.Name + "ResourcePath" }

// newResource creates the model for a resource message.
func newResource(name string, patterns []string) (*resource, error) {
	r := &resource{Name: name}
	seen := make(map[string]bool)
	for _, value := range patterns {
		p := pattern{
			Value:       value,
			Const:       constName(value),
			Constructor: constructorName(value),
		}
		for e := range resourcepath.Elements(value) {
			if !e.IsVariable() {
				continue
			}
			v := string(e.GetLiteral())
			if !variableRegexp.MatchString(v) {
				return nil, fmt.Errorf("%s: pattern %q: variable %s must match %s", name, value, e, variableRegexp)
			}
			p.Variables = append(p.Variables, variable{Name: v, Param: paramName(v), Getter: getterName(v)})
			if !seen[v] {
				seen[v] = true
				r.Variables = append(r.Variables, p.Variables[len(p.Variables)-1])
			}
		}
		r.Patterns = append(r.Patterns, p)
	}
	return r, nil
}
