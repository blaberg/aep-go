package resourcepath

import (
	"strings"
)

// Scanner scans a resource path.
type Scanner struct {
	path                     string
	start, end               int
	serviceStart, serviceEnd int
	full                     bool
}

// Init initializes the scanner.
func (s *Scanner) Init(path string) {
	s.path = path
	s.start, s.end = 0, 0
	s.full = false
}

// Scan to the next segment.
func (s *Scanner) Scan() bool {
	switch s.end {
	case len(s.path):
		return false
	case 0:
		// Special case for full resource paths.
		if strings.HasPrefix(s.path, "//") {
			s.full = true
			s.start, s.end = 2, 2
			nextSlash := strings.IndexByte(s.path[s.start:], '/')
			if nextSlash == -1 {
				s.serviceStart, s.serviceEnd = s.start, len(s.path)
				s.start, s.end = len(s.path), len(s.path)
				return false
			}
			s.serviceStart, s.serviceEnd = s.start, s.start+nextSlash
			s.start, s.end = s.start+nextSlash+1, s.start+nextSlash+1
		} else if strings.HasPrefix(s.path, "/") {
			s.start = s.end + 1 // start past beginning slash
		}
	default:
		s.start = s.end + 1 // start past latest slash
	}
	if nextSlash := strings.IndexByte(s.path[s.start:], '/'); nextSlash == -1 {
		s.end = len(s.path)
	} else {
		s.end = s.start + nextSlash
	}
	return true
}

// Start returns the start index (inclusive) of the current segment.
func (s *Scanner) Start() int {
	return s.start
}

// End returns the end index (exclusive) of the current segment.
func (s *Scanner) End() int {
	return s.end
}

// Segment returns the current segment.
func (s *Scanner) Segment() Segment {
	return Segment(s.path[s.start:s.end])
}

// Full returns true if the scanner has detected a full resource path.
func (s *Scanner) Full() bool {
	return s.full
}

// ServicePath returns the service path, when the scanner has detected a full resource path.
func (s *Scanner) ServicePath() string {
	return s.path[s.serviceStart:s.serviceEnd]
}
