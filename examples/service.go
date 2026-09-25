package examples

import (
	"context"
	"time"

	"github.com/blaberg/aep-go/pagination"
	booksv1 "github.com/blaberg/aep-go/proto/gen/example/books/v1"
	"github.com/blaberg/aep-go/resourceid"
	"github.com/blaberg/aep-go/resourcepath"
	"github.com/blaberg/aep-go/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ booksv1.BookServiceServer = &Service{}

// authorPattern is the pattern of the parent of a book written by an author.
// There is no Author message to generate code for, so the resourcepath package is used directly.
const authorPattern = "authors/{author_id}"

// maxPageSize is the default and maximum number of books returned by ListBooks.
const maxPageSize = 100

// Service implements the BookService.
//
// A book is either written by an author, with a path like "authors/tolkien/books/the-hobbit",
// or a top-level book without a single author, such as an anthology, with a path like
// "books/beowulf". Books written by an author can only be changed by that author, while
// top-level books can only be changed by an admin.
type Service struct {
	paginator *pagination.Paginator
	storage   *Storage
}

// NewService creates a new Service with empty storage.
func NewService() *Service {
	return &Service{
		paginator: pagination.NewPaginator(),
		storage:   NewStorage(),
	}
}

// CreateBook implements the CreateBook RPC.
// Books created with an author as parent are written by that author,
// while books created without a parent are top-level books.
func (s *Service) CreateBook(ctx context.Context, req *booksv1.CreateBookRequest) (*booksv1.Book, error) {
	// Generate book ID if not provided
	bookID := req.GetId()
	if bookID == "" {
		bookID = resourceid.New()
	}
	if err := validate.ResourceID(bookID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book ID: %v", err)
	}

	// Create the book's resource path, with the constructor for the parent's pattern
	var path *booksv1.BookResourcePath
	var err error
	if req.GetParent() == "" {
		path, err = booksv1.NewBookResourcePath(bookID)
	} else {
		parent, parseErr := resourcepath.ParseString(req.GetParent(), authorPattern)
		if parseErr != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", parseErr)
		}
		path, err = booksv1.NewAuthorBookResourcePath(parent.Get("author_id"), bookID)
	}
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book path: %v", err)
	}

	// Create the book
	now := timestamppb.New(time.Now())
	book := &booksv1.Book{
		Path:        path.String(),
		DisplayName: req.GetBook().GetDisplayName(),
		CreateTime:  now,
		UpdateTime:  now,
	}
	if !s.storage.Create(book) {
		return nil, status.Errorf(codes.AlreadyExists, "book %q already exists", book.GetPath())
	}
	return book, nil
}

// GetBook implements the GetBook RPC.
func (s *Service) GetBook(ctx context.Context, req *booksv1.GetBookRequest) (*booksv1.Book, error) {
	path, err := booksv1.ParseBookResourcePath(req.GetPath())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid path: %v", err)
	}
	book, ok := s.storage.Get(path.String())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetPath())
	}
	return book, nil
}

// ListBooks implements the ListBooks RPC.
// With an author as parent it lists the books written by that author,
// and without a parent it lists the top-level books.
func (s *Service) ListBooks(ctx context.Context, req *booksv1.ListBooksRequest) (*booksv1.ListBooksResponse, error) {
	// Validate the parent
	if req.GetParent() != "" {
		if _, err := resourcepath.ParseString(req.GetParent(), authorPattern); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
		}
	}

	// Parse the page token
	token, err := s.paginator.ParsePageToken(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token: %v", err)
	}

	// Use the maximum page size if none, or a too large one, is given
	pageSize := req.GetMaxPageSize()
	if pageSize <= 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	// Get books with pagination
	books, hasMore, err := s.storage.List(req.GetParent(), token.Offset, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list books: %v", err)
	}

	// Generate the next page token
	nextPageToken := ""
	if next := token.Next(hasMore, pageSize); next != nil {
		nextPageToken = next.String()
	}
	return &booksv1.ListBooksResponse{
		Results:       books,
		NextPageToken: nextPageToken,
	}, nil
}

// UpdateBook implements the UpdateBook RPC.
func (s *Service) UpdateBook(ctx context.Context, req *booksv1.UpdateBookRequest) (*booksv1.Book, error) {
	path, err := booksv1.ParseBookResourcePath(req.GetPath())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid path: %v", err)
	}
	if err := authorize(ctx, path); err != nil {
		return nil, err
	}
	book, ok := s.storage.Get(path.String())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetPath())
	}

	// Update the fields in the update mask, or all updatable fields if there is none
	fields := req.GetUpdateMask().GetPaths()
	if len(fields) == 0 {
		fields = []string{"display_name"}
	}
	updated := proto.Clone(book).(*booksv1.Book)
	for _, field := range fields {
		switch field {
		case "display_name":
			updated.DisplayName = req.GetBook().GetDisplayName()
		default:
			return nil, status.Errorf(codes.InvalidArgument, "field %q can't be updated", field)
		}
	}
	updated.UpdateTime = timestamppb.New(time.Now())
	if !s.storage.Update(updated) {
		return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetPath())
	}
	return updated, nil
}

// DeleteBook implements the DeleteBook RPC.
func (s *Service) DeleteBook(ctx context.Context, req *booksv1.DeleteBookRequest) (*emptypb.Empty, error) {
	path, err := booksv1.ParseBookResourcePath(req.GetPath())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid path: %v", err)
	}
	if err := authorize(ctx, path); err != nil {
		return nil, err
	}
	if !s.storage.Delete(path.String()) {
		return nil, status.Errorf(codes.NotFound, "book %q not found", req.GetPath())
	}
	return &emptypb.Empty{}, nil
}

// authorize checks that the caller may change the book.
// Who that is depends on which pattern the path has.
func authorize(ctx context.Context, path *booksv1.BookResourcePath) error {
	caller := callerFromContext(ctx)
	switch path.Pattern() {
	case booksv1.AuthorBookPattern:
		if caller != path.GetAuthorId() {
			return status.Errorf(codes.PermissionDenied, "only author %q can change this book", path.GetAuthorId())
		}
	case booksv1.BookPattern:
		if caller != "admin" {
			return status.Error(codes.PermissionDenied, "only admins can change top-level books")
		}
	}
	return nil
}

// callerFromContext returns the caller of the request, read from the "x-user" metadata.
// This is a stand-in for real authentication, to keep the example short.
func callerFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if users := md.Get("x-user"); len(users) > 0 {
		return users[0]
	}
	return ""
}
