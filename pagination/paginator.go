package pagination

type Paginator struct {
	maxPageSize int
	checksum    uint32
}

type Option func(*Paginator)

func NewPaginator(options ...Option) *Paginator {
	p := &Paginator{
		maxPageSize: 100,
		checksum:    0x9acb0442,
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

func WithCustomChecksum(checksum uint32) Option {
	return func(p *Paginator) {
		p.checksum = checksum
	}
}
