package resourcepath

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Sscan scans a resource path, storing successive segments into successive variables
// as determined by the provided pattern.
func Sscan(path, pattern string, value reflect.Value) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse resource path '%s' with pattern '%s': %w", path, pattern, err)
		}
	}()
	var pathScanner, patternScanner Scanner
	patternScanner.Init(pattern)
	variables := make([]reflect.Value, 0)
	for patternScanner.Scan() {
		if !patternScanner.Segment().IsVariable() {
			continue
		}
		v := value.FieldByNameFunc(func(s string) bool {
			return strings.ToLower(s) == patternScanner.Segment().Literal().ResourceID()
		})
		if !v.CanSet() || v.Kind() != reflect.String {
			return fmt.Errorf("failed to unmarshal path")
		}
		variables = append(variables, v)
	}
	pathScanner.Init(path)
	patternScanner.Init(pattern)
	var i int
	for patternScanner.Scan() {
		if patternScanner.Full() {
			return fmt.Errorf("invalid pattern")
		}
		if !pathScanner.Scan() {
			return fmt.Errorf("segment %s: %w", patternScanner.Segment(), io.ErrUnexpectedEOF)
		}
		pathSegment, patternSegment := pathScanner.Segment(), patternScanner.Segment()
		if !patternSegment.IsVariable() {
			if patternSegment.Literal() != pathSegment.Literal() {
				return fmt.Errorf("segment %s: got %s", patternSegment, pathSegment)
			}
			continue
		}
		variables[i].SetString(pathSegment.Literal().ResourceID())
		i++
	}
	if pathScanner.Scan() {
		return fmt.Errorf("got trailing segments in path")
	}
	return nil
}
