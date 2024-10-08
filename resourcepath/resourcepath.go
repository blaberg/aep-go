package resourcepath

import (
	"fmt"
	"io"
)

// todo: det slutgiliga målet är att vi vill kunna generara rsurrser som använder den här funktionen.
// De ska genereras som {message}ResourcePath. Jag vill att de ska ha GetOrganization() funktioner som returnerar strings
// för varje segmnet av resurs namnet. Tror det blir en snyggare lösning än det vi har nu
// alla validate och hasWilcard functioner ska ckså finnas på det
// generellt tror jag vi kan göra paketet mycket simplare

type ResourcePath struct {
	segments map[string]string
}

func (p ResourcePath) Get(segment string) string {
	return p.segments[segment]
}

func ParseString(path, pattern string) (*ResourcePath, error) {
	if pattern == "" {
		return nil, fmt.Errorf("pattern can't be empty")
	}
	segments := make(map[string]string)
	var pathScanner, patternScanner Scanner
	patternScanner.Init(pattern)
	pathScanner.Init(path)
	for patternScanner.Scan() {
		if patternScanner.Full() {
			return nil, fmt.Errorf("invalid pattern")
		}
		if !pathScanner.Scan() {
			return nil, fmt.Errorf("segment %s: %w", patternScanner.Segment(), io.ErrUnexpectedEOF)
		}
		pathSegment, patternSegment := pathScanner.Segment(), patternScanner.Segment()
		if !patternSegment.IsVariable() {
			if patternSegment.Literal() != pathSegment.Literal() {
				return nil, fmt.Errorf("segment %s: got %s", patternSegment, pathSegment)
			}
			if patternSegment.Literal() == Wildcard {
				return nil, fmt.Errorf("the pattern can't contain wildcards")
			}
			continue
		}
		segments[patternSegment.Literal().ResourceID()] = pathSegment.Literal().ResourceID()
	}
	if pathScanner.Scan() {
		return nil, fmt.Errorf("got trailing segments in path")
	}
	return &ResourcePath{
		segments: segments,
	}, nil
}
