package resourceid

import (
	"encoding/base32"

	"github.com/google/uuid"
)

var encoder = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)

// New generates a new resource ID.
func New() string {
	id := uuid.New()
	return encoder.EncodeToString(id[:])
}
