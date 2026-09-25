package examples

import (
	"context"
	"testing"

	booksv1 "github.com/blaberg/aep-go/proto/gen/example/books/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"gotest.tools/v3/assert"
)

// newTestService returns a service with two books by tolkien, one by lewis and one top-level book.
func newTestService(t *testing.T) *Service {
	t.Helper()
	s := NewService()
	for _, req := range []*booksv1.CreateBookRequest{
		{Parent: "authors/tolkien", Id: "the-hobbit", Book: &booksv1.Book{DisplayName: "The Hobbit"}},
		{Parent: "authors/tolkien", Id: "the-silmarillion", Book: &booksv1.Book{DisplayName: "The Silmarillion"}},
		{Parent: "authors/lewis", Id: "narnia", Book: &booksv1.Book{DisplayName: "Narnia"}},
		{Id: "beowulf", Book: &booksv1.Book{DisplayName: "Beowulf"}},
	} {
		_, err := s.CreateBook(context.Background(), req)
		assert.NilError(t, err)
	}
	return s
}

// asCaller returns a context for a request made by caller.
func asCaller(caller string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs("x-user", caller))
}

func TestCreateBook(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		parent string
		id     string
		want   string
		code   codes.Code
	}{
		{
			name: "top-level book",
			id:   "the-odyssey",
			want: "books/the-odyssey",
		},
		{
			name:   "book by author",
			parent: "authors/tolkien",
			id:     "the-lord-of-the-rings",
			want:   "authors/tolkien/books/the-lord-of-the-rings",
		},
		{
			name:   "invalid parent",
			parent: "publishers/penguin",
			id:     "the-odyssey",
			code:   codes.InvalidArgument,
		},
		{
			name:   "already exists",
			parent: "authors/tolkien",
			id:     "the-hobbit",
			code:   codes.AlreadyExists,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newTestService(t)
			book, err := s.CreateBook(context.Background(), &booksv1.CreateBookRequest{
				Parent: tt.parent,
				Id:     tt.id,
				Book:   &booksv1.Book{DisplayName: "A Book"},
			})
			assert.Equal(t, status.Code(err), tt.code)
			if tt.code == codes.OK {
				assert.Equal(t, book.GetPath(), tt.want)
			}
		})
	}
}

func TestGetBook(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		path string
		want string
		code codes.Code
	}{
		{
			name: "book by author",
			path: "authors/tolkien/books/the-hobbit",
			want: "The Hobbit",
		},
		{
			name: "top-level book",
			path: "books/beowulf",
			want: "Beowulf",
		},
		{
			name: "book by other author",
			path: "authors/lewis/books/the-hobbit",
			code: codes.NotFound,
		},
		{
			name: "book by author as top-level book",
			path: "books/the-hobbit",
			code: codes.NotFound,
		},
		{
			name: "invalid path",
			path: "shelves/1",
			code: codes.InvalidArgument,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newTestService(t)
			book, err := s.GetBook(context.Background(), &booksv1.GetBookRequest{Path: tt.path})
			assert.Equal(t, status.Code(err), tt.code)
			if tt.code == codes.OK {
				assert.Equal(t, book.GetDisplayName(), tt.want)
			}
		})
	}
}

func TestListBooks(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		parent string
		want   []string
		code   codes.Code
	}{
		{
			name:   "books by author",
			parent: "authors/tolkien",
			want:   []string{"authors/tolkien/books/the-hobbit", "authors/tolkien/books/the-silmarillion"},
		},
		{
			name: "top-level books",
			want: []string{"books/beowulf"},
		},
		{
			name:   "author without books",
			parent: "authors/austen",
		},
		{
			name:   "invalid parent",
			parent: "publishers/penguin",
			code:   codes.InvalidArgument,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newTestService(t)
			resp, err := s.ListBooks(context.Background(), &booksv1.ListBooksRequest{Parent: tt.parent})
			assert.Equal(t, status.Code(err), tt.code)
			var got []string
			for _, book := range resp.GetResults() {
				got = append(got, book.GetPath())
			}
			assert.DeepEqual(t, got, tt.want)
		})
	}
}

func TestListBooks_Pagination(t *testing.T) {
	t.Parallel()
	s := newTestService(t)
	req := &booksv1.ListBooksRequest{Parent: "authors/tolkien", MaxPageSize: 1}
	var got []string
	for {
		resp, err := s.ListBooks(context.Background(), req)
		assert.NilError(t, err)
		for _, book := range resp.GetResults() {
			got = append(got, book.GetPath())
		}
		if resp.GetNextPageToken() == "" {
			break
		}
		req.PageToken = resp.GetNextPageToken()
	}
	assert.DeepEqual(t, got, []string{"authors/tolkien/books/the-hobbit", "authors/tolkien/books/the-silmarillion"})
}

func TestUpdateBook(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		caller string
		path   string
		mask   []string
		code   codes.Code
	}{
		{
			name:   "author updates own book",
			caller: "tolkien",
			path:   "authors/tolkien/books/the-hobbit",
		},
		{
			name:   "with update mask",
			caller: "tolkien",
			path:   "authors/tolkien/books/the-hobbit",
			mask:   []string{"display_name"},
		},
		{
			name:   "other author updates book",
			caller: "lewis",
			path:   "authors/tolkien/books/the-hobbit",
			code:   codes.PermissionDenied,
		},
		{
			name:   "admin updates top-level book",
			caller: "admin",
			path:   "books/beowulf",
		},
		{
			name:   "field that can't be updated",
			caller: "tolkien",
			path:   "authors/tolkien/books/the-hobbit",
			mask:   []string{"create_time"},
			code:   codes.InvalidArgument,
		},
		{
			name:   "missing book",
			caller: "tolkien",
			path:   "authors/tolkien/books/the-lord-of-the-rings",
			code:   codes.NotFound,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newTestService(t)
			req := &booksv1.UpdateBookRequest{
				Path: tt.path,
				Book: &booksv1.Book{DisplayName: "New Title"},
			}
			if tt.mask != nil {
				req.UpdateMask = &fieldmaskpb.FieldMask{Paths: tt.mask}
			}
			_, err := s.UpdateBook(asCaller(tt.caller), req)
			assert.Equal(t, status.Code(err), tt.code)
			if tt.code == codes.OK {
				book, err := s.GetBook(context.Background(), &booksv1.GetBookRequest{Path: tt.path})
				assert.NilError(t, err)
				assert.Equal(t, book.GetDisplayName(), "New Title")
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		caller string
		path   string
		code   codes.Code
	}{
		{
			name:   "author deletes own book",
			caller: "tolkien",
			path:   "authors/tolkien/books/the-hobbit",
		},
		{
			name:   "other author deletes book",
			caller: "lewis",
			path:   "authors/tolkien/books/the-hobbit",
			code:   codes.PermissionDenied,
		},
		{
			name:   "admin deletes top-level book",
			caller: "admin",
			path:   "books/beowulf",
		},
		{
			name:   "author deletes top-level book",
			caller: "tolkien",
			path:   "books/beowulf",
			code:   codes.PermissionDenied,
		},
		{
			name:   "missing book",
			caller: "tolkien",
			path:   "authors/tolkien/books/the-lord-of-the-rings",
			code:   codes.NotFound,
		},
		{
			name:   "invalid path",
			caller: "admin",
			path:   "shelves/1",
			code:   codes.InvalidArgument,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := newTestService(t)
			_, err := s.DeleteBook(asCaller(tt.caller), &booksv1.DeleteBookRequest{Path: tt.path})
			assert.Equal(t, status.Code(err), tt.code)
			if tt.code == codes.OK {
				_, err := s.GetBook(context.Background(), &booksv1.GetBookRequest{Path: tt.path})
				assert.Equal(t, status.Code(err), codes.NotFound)
			}
		})
	}
}
