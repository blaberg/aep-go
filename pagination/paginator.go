package pagination

// Paginator provides functionality for handling pagination in API requests.
// It maintains configuration for maximum page size and a signing key for token validation.
type Paginator struct {
	maxPageSize int
	signingKey  uint32
}

// Option is a function that configures a Paginator.
type Option func(*Paginator)

// NewPaginator creates a new Paginator with the given options.
// By default, it sets maxPageSize to 100 and uses a predefined signing key.
func NewPaginator(options ...Option) *Paginator {
	p := &Paginator{
		maxPageSize: 100,
		signingKey:  0xefbfde39,
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}

// WithCustomMaxPageSize returns an Option that sets a custom maximum page size for the Paginator.
func WithCustomMaxPageSize(maxPageSize int) Option {
	return func(p *Paginator) {
		p.maxPageSize = maxPageSize
	}
}

// WithCustomSigningKey returns an Option that sets a custom signing key for the Paginator.
// This key is used in the checksum calculation for page tokens.
func WithCustomSigningKey(key uint32) Option {
	return func(p *Paginator) {
		p.signingKey = key
	}
}
