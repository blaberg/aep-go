package examples

import (
	"sync"

	booksv1 "github.com/blaberg/aep-go/proto/gen/example/books/v1"
)

// Storage provides in-memory storage for books, in the order they were created.
type Storage struct {
	mu    sync.RWMutex
	books []*booksv1.Book
}

// NewStorage creates a new storage instance
func NewStorage() *Storage {
	return &Storage{
		books: make([]*booksv1.Book, 0),
	}
}

// Create stores a new book, and reports whether it was created.
// It returns false if a book with the same path already exists.
func (s *Storage) Create(book *booksv1.Book) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.index(book.GetPath()) >= 0 {
		return false
	}
	s.books = append(s.books, book)
	return true
}

// Get retrieves a book by its path
func (s *Storage) Get(path string) (*booksv1.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i := s.index(path); i >= 0 {
		return s.books[i], true
	}
	return nil, false
}

// List returns the books with the given parent, where an empty parent means
// top-level books, with pagination.
func (s *Storage) List(parent string, offset int64, pageSize int32) ([]*booksv1.Book, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filteredBooks []*booksv1.Book
	for _, book := range s.books {
		path, err := booksv1.ParseBookResourcePath(book.GetPath())
		if err != nil {
			return nil, false, err
		}
		if parentOf(path) == parent {
			filteredBooks = append(filteredBooks, book)
		}
	}

	// Apply pagination
	start := int(offset)
	if start < 0 || start >= len(filteredBooks) {
		return nil, false, nil
	}
	end := min(start+int(pageSize), len(filteredBooks))
	return filteredBooks[start:end], end < len(filteredBooks), nil
}

// Update replaces the book with the same path, and reports whether it existed.
func (s *Storage) Update(book *booksv1.Book) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i := s.index(book.GetPath()); i >= 0 {
		s.books[i] = book
		return true
	}
	return false
}

// Delete removes a book by its path, and reports whether the book existed.
func (s *Storage) Delete(path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i := s.index(path); i >= 0 {
		s.books = append(s.books[:i], s.books[i+1:]...)
		return true
	}
	return false
}

// index returns the index of the book with the given path, or -1 if there is none.
// The caller must hold the lock.
func (s *Storage) index(path string) int {
	for i, book := range s.books {
		if book.GetPath() == path {
			return i
		}
	}
	return -1
}

// parentOf returns the parent of a book: the author for a book written by an author,
// or an empty string for a top-level book.
func parentOf(path *booksv1.BookResourcePath) string {
	switch path.Pattern() {
	case booksv1.AuthorBookPattern:
		return "authors/" + path.GetAuthorId()
	default:
		return ""
	}
}
