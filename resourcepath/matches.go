package resourcepath

func Matches(pattern, path string) bool {
	var patternScanner, pathScanner Scanner
	patternScanner.Init(pattern)
	pathScanner.Init(path)
	for pathScanner.Scan() {
		if !patternScanner.Scan() {
			return false
		}
		if patternScanner.Segment().IsVariable() {
			continue // variables not allwed in paths
		}
		if pathScanner.Segment() != patternScanner.Segment() {
			return false
		}
	}
	return true
}
