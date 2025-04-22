# AEP-Go Example Service

This example demonstrates how to use the AEP-Go library to build a gRPC service that follows AEP standards. The service implements:

- Resource path validation and parsing
- Resource ID generation
- Pagination with token-based validation

## Service Overview

The example implements a simple Book Service with the following operations:
- `CreateBook`: Creates a new book with an optional or auto-generated ID
- `GetBook`: Retrieves a book by its resource path
- `ListBooks`: Lists books with pagination support

## Project Structure

```
.
├── service.go         # Service implementation
└── go.mod          # Go module file
```
