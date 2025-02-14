package resourcepath

import (
	"fmt"
	"io"
	"iter"
)

// todo: det slutgiliga målet är att vi vill kunna generara rsurrser som använder den här funktionen.
// De ska genereras som {message}ResourcePath. Jag vill att de ska ha GetOrganization() funktioner som returnerar strings
// för varje segmnet av resurs namnet. Tror det blir en snyggare lösning än det vi har nu
// alla validate och hasWilcard functioner ska ckså finnas på det
// generellt tror jag vi kan göra paketet mycket simplare

type ResourcePath struct {
	elements map[string]string
}

func (p ResourcePath) Get(element string) string {
	return p.elements[element]
}

func NewResourcePath(elements map[string]string) *ResourcePath {
	return &ResourcePath{
		elements: elements,
	}
}

func ParseString(path, pattern string) (*ResourcePath, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern can't be empty")
	}
	elements := make(map[string]string)
	if path == "" {
		return nil, fmt.Errorf("path can't be empty")
	}
	pathItr, patternItr := Elements(path), Elements(pattern)
	next, stop := iter.Pull(pathItr)
	defer stop()
	for pattrElem := range patternItr {
		pathElem, ok := next()
		if !ok {
			return nil, fmt.Errorf("element %s: %w", pattrElem, io.ErrUnexpectedEOF)
		}
		if !pattrElem.IsVariable() {
			if pattrElem.GetLiteral() != pathElem.GetLiteral() {
				return nil, fmt.Errorf("element %s: got %s", pattrElem, pathElem)
			}
			continue
		}
		elements[string(pattrElem.GetLiteral())] = string(pathElem.GetLiteral())
	}
	if _, ok := next(); ok {
		return nil, fmt.Errorf("got trailing elements in path")
	}
	return &ResourcePath{
		elements: elements,
	}, nil
}
