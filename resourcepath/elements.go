package resourcepath

import (
	"iter"
	"strings"
)

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

func (e Element) IsVariable() bool {
	return len(e) > 2 &&
		strings.HasPrefix(string(e), "{") &&
		strings.HasSuffix(string(e), "}")
}

func (e Element) GetLiteral() Literal {
	if e.IsVariable() {
		return Literal(e[1 : len(e)-1])
	}
	return Literal(e)
}

func (e Element) IsWildcard() bool {
	return e == "-"
}
