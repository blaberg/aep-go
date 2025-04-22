package examples

import (
	"sync"

	bookv1 "github.com/blaberg/aep-go/proto/gen/example/books/v1"
)

// Storage provides in-memory storage for books
type Storage struct {
	mu    sync.RWMutex
	books []*bookv1.Book
}

// NewStorage creates a new storage instance
func NewStorage() *Storage {
	return &Storage{
		books: make([]*bookv1.Book, 0),
	}
}

// Create stores a new book
func (s *Storage) Create(book *bookv1.Book) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.books = append(s.books, book)
}

// Get retrieves a book by its path
func (s *Storage) Get(path string) (*bookv1.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, book := range s.books {
		if book.Path == path {
			return book, true
		}
	}
	return nil, false
}

// List returns books that match the given parent path with pagination
func (s *Storage) List(parent string, offset int64, pageSize int32) ([]*bookv1.Book, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filteredBooks []*bookv1.Book
	for _, book := range s.books {
		path, err := bookv1.ParseAuthorBookResourcePath(book.Path)
		if err != nil {
			return nil, false, err
		}
		if path.GetAuthor() == parent {
			filteredBooks = append(filteredBooks, book)
		}
	}

	// Apply pagination
	start := int(offset)
	end := start + int(pageSize)
	if end > len(filteredBooks) {
		end = len(filteredBooks)
	}
	if start >= len(filteredBooks) {
		return nil, false, nil
	}

	// Check if there are more pages
	hasMore := end < len(filteredBooks)

	return filteredBooks[start:end], hasMore, nil
}

// Delete removes a book by its path
func (s *Storage) Delete(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, book := range s.books {
		if book.Path == path {
			// Remove the book by swapping with the last element and truncating
			s.books[i] = s.books[len(s.books)-1]
			s.books = s.books[:len(s.books)-1]
			return
		}
	}
}
