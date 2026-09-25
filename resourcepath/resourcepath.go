package resourcepath

import (
	"fmt"
	"io"
	"iter"
	"slices"
	"strings"
)

// ResourcePath represents an AEP resource path.
type ResourcePath struct {
	pattern  string
	elements map[string]string
}

// Get returns the value of the element.
// If the element is not found, an empty string is returned.
func (p ResourcePath) Get(element string) string {
	return p.elements[element]
}

// Pattern returns the pattern the resource path was created with.
func (p ResourcePath) Pattern() string {
	return p.pattern
}

// String returns the resource path, with each variable in the pattern replaced
// by its value. Variables without a value are replaced by an empty string.
func (p ResourcePath) String() string {
	var parts []string
	for e := range Elements(p.pattern) {
		if e.IsVariable() {
			parts = append(parts, p.Get(string(e.GetLiteral())))
			continue
		}
		parts = append(parts, string(e))
	}
	return strings.Join(parts, "/")
}

// NewResourcePath creates a new ResourcePath from a pattern and the values of its variables.
//
// It returns an error if a variable is repeated in the pattern, a variable in the
// pattern has no value, a value is empty or contains a "/", or a value is given for
// a variable that is not in the pattern.
// The String of the returned ResourcePath can always be parsed with ParseString.
func NewResourcePath(pattern string, values map[string]string) (*ResourcePath, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern can't be empty")
	}
	elements := make(map[string]string, len(values))
	for e := range Elements(pattern) {
		if !e.IsVariable() {
			continue
		}
		variable := string(e.GetLiteral())
		value, ok := values[variable]
		switch {
		case elements[variable] != "":
			return nil, fmt.Errorf("element %s: repeated variable", e)
		case !ok:
			return nil, fmt.Errorf("element %s: missing value", e)
		case value == "":
			return nil, fmt.Errorf("element %s: empty value", e)
		case strings.Contains(value, "/"):
			return nil, fmt.Errorf("element %s: value %q contains \"/\"", e, value)
		}
		elements[variable] = value
	}
	if len(elements) != len(values) {
		var unknown []string
		for variable := range values {
			if _, ok := elements[variable]; !ok {
				unknown = append(unknown, variable)
			}
		}
		slices.Sort(unknown)
		return nil, fmt.Errorf("element {%s}: not in pattern", unknown[0])
	}
	return &ResourcePath{
		elements: elements,
		pattern:  pattern,
	}, nil
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
			if pattrElem.GetLiteral() != Literal(pathElem) {
				return nil, fmt.Errorf("element %s: got %s", pattrElem, pathElem)
			}
			continue
		}
		variable := string(pattrElem.GetLiteral())
		if _, ok := elements[variable]; ok {
			return nil, fmt.Errorf("element %s: repeated variable", pattrElem)
		}
		if len(pathElem) == 0 {
			return nil, fmt.Errorf("element %s: empty value", pattrElem)
		}
		elements[variable] = string(pathElem)
	}
	if _, ok := next(); ok {
		return nil, fmt.Errorf("got trailing elements in path")
	}
	return &ResourcePath{
		elements: elements,
		pattern:  pattern,
	}, nil
}
