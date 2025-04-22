package pagination

import (
	"testing"

	booksv1 "github.com/blaberg/aep-go/proto/gen/example/books/v1"
	"gotest.tools/v3/assert"
)

func TestParseOffsetPageToken(t *testing.T) {
	t.Parallel()
	t.Run("valid checksums", func(t *testing.T) {
		t.Parallel()
		p := NewPaginator()
		request1 := &booksv1.ListBooksRequest{
			Parent:      "shelves/1",
			MaxPageSize: 10,
		}
		pageToken1, err := p.ParsePageToken(request1)
		assert.NilError(t, err)
		request2 := &booksv1.ListBooksRequest{
			Parent:      "shelves/1",
			MaxPageSize: 20,
			PageToken:   pageToken1.Next(true, request1.GetMaxPageSize()).String(),
		}
		pageToken2, err := p.ParsePageToken(request2)
		assert.NilError(t, err)
		assert.Equal(t, int64(10), pageToken2.Offset)
		request3 := &booksv1.ListBooksRequest{
			Parent:      "shelves/1",
			MaxPageSize: 30,
			PageToken:   pageToken2.Next(true, request2.GetMaxPageSize()).String(),
		}
		pageToken3, err := p.ParsePageToken(request3)
		assert.NilError(t, err)
		assert.Equal(t, int64(30), pageToken3.Offset)
	})

	t.Run("invalid format", func(t *testing.T) {
		t.Parallel()
		p := NewPaginator()
		request := &booksv1.ListBooksRequest{
			Parent:      "shelves/1",
			MaxPageSize: 10,
			PageToken:   "invalid",
		}
		pageToken1, err := p.ParsePageToken(request)
		assert.ErrorContains(t, err, "decode")
		assert.Equal(t, PageToken{}, pageToken1)
	})

	t.Run("invalid checksum", func(t *testing.T) {
		t.Parallel()
		p := NewPaginator()
		pt := &PageToken{
			Offset:   100,
			Checksum: 1234, // invalid
		}
		request := &booksv1.ListBooksRequest{
			Parent:      "shelves/1",
			MaxPageSize: 10,
			PageToken:   pt.String(),
		}
		pageToken1, err := p.ParsePageToken(request)
		assert.ErrorContains(t, err, "checksum")
		assert.Equal(t, PageToken{}, pageToken1)
	})
}
