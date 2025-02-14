package pagination

type Paginator struct {
	maxPageSize int
	signingKey  uint32
}

type Option func(*Paginator)

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

func WithCustomMaxPageSize(maxPageSize int) Option {
	return func(p *Paginator) {
		p.maxPageSize = maxPageSize
	}
}

func WithCustomSigningKey(key uint32) Option {
	return func(p *Paginator) {
		p.signingKey = key
	}
}
