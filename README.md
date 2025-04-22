# AEP Go

Go SDK implementation of [API Extension Proposals](https://aep.dev/)

## Overview

This project contains a collection of helper functions to make it easier to adopt
AEP. It also contains a code generator that can be used together with Protool Buffers,
to make the proces even simpler

## Installation

```bash
go get github.com/blaberg/aep-go
```

## Usage

### List functionality

```go
func (s *Service) ListBooks(ctx context.Context, req *bookv1.ListBooksRequest) (*bookv1.ListBooksResponse, error) {
	// Validate parent format
	_, err := resourcepath.ParseString(req.Parent, "publishers/{publisher}")
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent format: %v", err)
	}

	// Parse the page token
	token, err := s.paginator.ParsePageToken(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token: %v", err)
	}

	// Get books with pagination
	books, hasMore := s.storage.List(req.Parent, token.Offset, req.MaxPageSize)
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
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the [LICENSE NAME] - see the [LICENSE](LICENSE) file for details.

