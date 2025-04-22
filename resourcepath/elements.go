package resourcepath

import (
	"iter"
	"strings"
)

// Elements returns an iterator over the elements of the resource path.
func Elements(path string) iter.Seq[Element] {
	return func(yield func(Element) bool) {
		elements := strings.Split(path, "/")
		for _, element := range elements {
			if !yield(Element(element)) {
				return
			}
		}
	}
}

// EBNF
//
//	element  = variable | literal ;
//	variable = "{" literal "}" ;
type Element string
type Literal string

// IsVariable returns true if the element is a variable.
func (e Element) IsVariable() bool {
	return len(e) > 2 &&
		strings.HasPrefix(string(e), "{") &&
		strings.HasSuffix(string(e), "}")
}

// GetLiteral returns the literal value of the element.
func (e Element) GetLiteral() Literal {
	if e.IsVariable() {
		return Literal(e[1 : len(e)-1])
	}
	return Literal(e)
}

// IsWildcard returns true if the element is a wildcard.
func (e Element) IsWildcard() bool {
	return e == "-"
}
