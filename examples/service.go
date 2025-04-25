package examples

import (
	"context"
	"time"

	"github.com/blaberg/aep-go/pagination"
	booksv1 "github.com/blaberg/aep-go/proto/gen/example/books/v1"
	bookv1 "github.com/blaberg/aep-go/proto/gen/example/books/v1"
	"github.com/blaberg/aep-go/resourceid"
	"github.com/blaberg/aep-go/validate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ booksv1.BookServiceServer = &Service{}

// Service implements the BookService
type Service struct {
	paginator *pagination.Paginator
	storage   *Storage
}

// CreateBook implements the CreateBook RPC
func (s *Service) CreateBook(ctx context.Context, req *bookv1.CreateBookRequest) (*bookv1.Book, error) {
	// Generate book ID if not provided
	bookID := req.Id
	if bookID == "" {
		bookID = resourceid.New()
	}
	if err := validate.ResourceID(bookID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book ID: %v", err)
	}

	// Create the book's resource path
	path := bookv1.NewBookPath(bookID)

	// Create the book
	now := timestamppb.New(time.Now())
	book := &bookv1.Book{
		Path:        path.String(),
		DisplayName: req.Book.DisplayName,
		CreateTime:  now,
		UpdateTime:  now,
	}

	// Store the book
	s.storage.Create(book)

	return book, nil
}

// GetBook implements the GetBook RPC
func (s *Service) GetBook(ctx context.Context, req *bookv1.GetBookRequest) (*bookv1.Book, error) {
	// Validate path format
	path, err := bookv1.ParseBookResourcePath(req.GetPath())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid resource path format: %v", err)
	}

	// Get the book
	book, ok := s.storage.Get(path.GetBook())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "book not found: %s", req.Path)
	}

	return book, nil
}

// ListBooks implements the ListBooks RPC
func (s *Service) ListBooks(ctx context.Context, req *bookv1.ListBooksRequest) (*bookv1.ListBooksResponse, error) {
	// Validate parent format
	path, err := bookv1.ParseAuthorBookResourcePath(req.GetParent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent format: %v", err)
	}

	// Parse the page token
	token, err := s.paginator.ParsePageToken(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token: %v", err)
	}

	// Get books with pagination
	books, hasMore, err := s.storage.List(path.GetAuthor(), token.Offset, req.MaxPageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list books: %v", err)
	}
	if books == nil {
		return &bookv1.ListBooksResponse{
			Results:       []*bookv1.Book{},
			NextPageToken: "",
		}, nil
	}

	// Generate the next page token
	nextToken := token.Next(hasMore, req.MaxPageSize)
	nextPageToken := ""
	if nextToken != nil {
		nextPageToken = nextToken.String()
	}

	return &bookv1.ListBooksResponse{
		Results:       books,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *Service) DeleteBook(
	context.Context,
	*booksv1.DeleteBookRequest,
) (*emptypb.Empty, error) {
	panic("implement me")
}

func (s *Service) UpdateBook(
	context.Context,
	*booksv1.UpdateBookRequest,
) (*booksv1.Book, error) {
	panic("implement me")
}
