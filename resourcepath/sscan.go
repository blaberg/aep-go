package resourcepath

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Sscan scans a resource name, storing successive segments into successive variables
// as determined by the provided pattern.
func Sscan(name, pattern string, value reflect.Value) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse resource path '%s' with pattern '%s': %w", name, pattern, err)
		}
	}()
	var nameScanner, patternScanner Scanner
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
	nameScanner.Init(name)
	patternScanner.Init(pattern)
	var i int
	for patternScanner.Scan() {
		if patternScanner.Full() {
			return fmt.Errorf("invalid pattern")
		}
		if !nameScanner.Scan() {
			return fmt.Errorf("segment %s: %w", patternScanner.Segment(), io.ErrUnexpectedEOF)
		}
		nameSegment, patternSegment := nameScanner.Segment(), patternScanner.Segment()
		if !patternSegment.IsVariable() {
			if patternSegment.Literal() != nameSegment.Literal() {
				return fmt.Errorf("segment %s: got %s", patternSegment, nameSegment)
			}
			continue
		}
		variables[i].SetString(nameSegment.Literal().ResourceID())
		i++
	}
	if nameScanner.Scan() {
		return fmt.Errorf("got trailing segments in name")
	}
	return nil
}
