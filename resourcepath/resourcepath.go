package resourcepath

import (
	"fmt"
	"io"
	"iter"
)

// ResourcePath represents an AEP resource path.
type ResourcePath struct {
	elements map[string]string
}

// Get returns the value of the element.
// If the element is not found, an empty string is returned.
func (p ResourcePath) Get(element string) string {
	return p.elements[element]
}

// NewResourcePath creates a new ResourcePath.
func NewResourcePath(elements map[string]string) *ResourcePath {
	return &ResourcePath{
		elements: elements,
	}
}

// ParseString parses a path and a pattern and returns a ResourcePath.
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
