package resourcepath

import (
	"iter"
)

func Matches(pattern, path string) bool {
	patternElements := Elements(pattern)
	next, stop := iter.Pull(patternElements)
	defer stop()
	pathElements := Elements(path)
	for e := range pathElements {
		p, ok := next()
		if !ok {
			return false
		}
		if p.IsVariable() {
			return false // variables not allwed in paths
		}
		if e != p {
			return false
		}
	}
	// check for remainers
	if _, ok := next(); ok {
		return false
	}
	return true
}
