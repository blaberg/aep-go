package resourceid

import (
	"encoding/base32"

	"github.com/google/uuid"
)

var encoder = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz0123456789").WithPadding(base32.NoPadding)

func New() string {
	id := uuid.New()
	return encoder.EncodeToString(id[:])
}
