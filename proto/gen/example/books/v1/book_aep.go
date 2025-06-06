// Code generated by protoc-gen-go-aep. DO NOT EDIT.
//
// versions:
// 	protoc-gen-go-aep development
// 	protoc (unknown)
// source: example/books/v1/book.proto

package booksv1

import (
	fmt "fmt"
	resourcepath "github.com/blaberg/aep-go/resourcepath"
	strings "strings"
)

type AuthorBookResourcePath struct {
	path *resourcepath.ResourcePath
}

func ParseAuthorBookResourcePath(p string) (*AuthorBookResourcePath, error) {
	path, err := resourcepath.ParseString(p, "authors/{author}/books/{book}")
	if err != nil {
		return nil, err
	}
	return &AuthorBookResourcePath{
		path: path,
	}, nil
}

func NewAuthorBookPath(
	author string,
	book string,
) *AuthorBookResourcePath {
	segments := map[string]string{
		"author": author,
		"book":   book,
	}
	return &AuthorBookResourcePath{
		path: resourcepath.NewResourcePath(segments),
	}
}

func (p *AuthorBookResourcePath) String() string {
	return strings.Join(
		[]string{
			"authors",
			p.path.Get("author"),
			"books",
			p.path.Get("book"),
		},
		"/",
	)
}

func (p *AuthorBookResourcePath) GetAuthor() string {
	return p.path.Get("author")
}

func (p *AuthorBookResourcePath) GetBook() string {
	return p.path.Get("book")
}

type BookResourcePath struct {
	path *resourcepath.ResourcePath
}

func ParseBookResourcePath(p string) (*BookResourcePath, error) {
	path, err := resourcepath.ParseString(p, "books/{book}")
	if err != nil {
		return nil, err
	}
	return &BookResourcePath{
		path: path,
	}, nil
}

func NewBookPath(
	book string,
) *BookResourcePath {
	segments := map[string]string{
		"book": book,
	}
	return &BookResourcePath{
		path: resourcepath.NewResourcePath(segments),
	}
}

func (p *BookResourcePath) String() string {
	return strings.Join(
		[]string{
			"books",
			p.path.Get("book"),
		},
		"/",
	)
}

func (p *BookResourcePath) GetBook() string {
	return p.path.Get("book")
}

type isMultipattern interface {
	isMultipattern()
}

func (*AuthorBookResourcePath) isMultipattern() {}

func (*BookResourcePath) isMultipattern() {}

func ParseMultipattern(p string) (isMultipattern, error) {
	switch {
	case resourcepath.Matches("authors/{author}/books/{book}", p):
		return ParseAuthorBookResourcePath(p)
	case resourcepath.Matches("books/{book}", p):
		return ParseBookResourcePath(p)
	}
	return nil, fmt.Errorf("failed to match pattern")
}

type MultipatternResourcePath struct {
	path *resourcepath.ResourcePath
}

func ParseMultipatternResourcePath(p string) (*MultipatternResourcePath, error) {
	var path *resourcepath.ResourcePath
	var err error
	switch {
	case resourcepath.Matches("authors/{author}/books/{book}", p):
		path, err = resourcepath.ParseString("authors/{author}/books/{book}", p)
	case resourcepath.Matches("books/{book}", p):
		path, err = resourcepath.ParseString("books/{book}", p)
	}
	if err != nil {
		return nil, err
	}
	return &MultipatternResourcePath{
		path: path,
	}, nil
}

func (p *MultipatternResourcePath) GetAuthor() string {
	return p.path.Get("author")
}

func (p *MultipatternResourcePath) GetBook() string {
	return p.path.Get("book")
}
