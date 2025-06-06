package pagination

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

// PageToken represents a pagination token with an offset and a checksum used for validating pagination sequences.
type PageToken struct {
	Offset   int64
	Checksum uint32
}

// Next generates the next PageToken by incrementing the offset using the given page size if more pages are available.
// If hasMore is false, it returns nil indicating there are no more pages.
func (p *PageToken) Next(hasMore bool, pageSize int32) *PageToken {
	if !hasMore {
		return nil
	}
	return &PageToken{
		Checksum: p.Checksum,
		Offset:   p.Offset + int64(pageSize),
	}
}

// String returns the Base64 URL-encoded string representation of the PageToken's offset and checksum fields.
func (p *PageToken) String() string {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint64(buf[:8], uint64(p.Offset))
	binary.BigEndian.PutUint32(buf[8:12], p.Checksum)
	return base64.URLEncoding.EncodeToString(buf)
}

// ParsePageToken decodes and validates a pagination token from the given request, returning the corresponding PageToken.
func (p *Paginator) ParsePageToken(request Request) (PageToken, error) {
	checksum, err := calculateRequestChecksum(request)
	if err != nil {
		return PageToken{}, err
	}
	checksum ^= p.signingKey
	if request.GetPageToken() == "" {
		return PageToken{
			Offset:   0,
			Checksum: checksum,
		}, nil
	}
	var pageToken PageToken
	if err := decodePageToken(request.GetPageToken(), &pageToken); err != nil {
		return PageToken{}, fmt.Errorf("failed to decode pagetokenpageToken")
	}

	if pageToken.Checksum != checksum {
		return PageToken{}, fmt.Errorf("missmatch checksum")
	}
	return pageToken, nil
}

func decodePageToken(s string, p *PageToken) error {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}
	p.Offset = int64(binary.BigEndian.Uint64(data[0:8]))
	p.Checksum = binary.BigEndian.Uint32(data[8:12])
	return nil
}
