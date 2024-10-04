package resourcepath

// Match reports whether the specified resource path matches the specified resource path pattern.
func Match(pattern, path string) bool {
	var pathScanner, patternScanner Scanner
	pathScanner.Init(path)
	patternScanner.Init(pattern)
	for patternScanner.Scan() {
		if !pathScanner.Scan() {
			return false
		}
		pathSegment := pathScanner.Segment()
		if pathSegment.IsVariable() {
			return false
		}
		patternSegment := patternScanner.Segment()
		if patternSegment.IsWildcard() {
			return false // edge case - wildcard not allowed in patterns
		}
		if patternSegment.IsVariable() {
			if pathSegment == "" {
				return false
			}
		} else if pathSegment != patternSegment {
			return false
		}
	}
	switch {
	case
		pathScanner.Scan(),             // path has more segments than pattern, no match
		patternScanner.Segment() == "", // edge case - empty pattern never matches
		patternScanner.Full():          // edge case - full resource path not allowed in patterns
		return false
	}
	return true
}
