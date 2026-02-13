package iter

import (
	"errors"
	"sync"
)

const (
	MinChunkSize uint = 256
	MaxChunkSize uint = 4096
)

type ParIter struct {
	errors []error

	wg    *sync.WaitGroup
	errCh chan error
}

func NewParIter() *ParIter {
	return &ParIter{
		errors: []error{},
		wg:     &sync.WaitGroup{},
		errCh:  make(chan error),
	}
}

func (p *ParIter) Wait() error {
	p.wg.Wait()

	close(p.errCh)
	for err := range p.errCh {
		p.errors = append(p.errors, err)
	}

	return errors.Join(p.errors...)
}

type ParIterOptions struct {
	MaxChunkSize uint
	MinChunkSize uint
}

func DefaultParIterOptions() *ParIterOptions {
	return &ParIterOptions{
		MaxChunkSize: MaxChunkSize,
		MinChunkSize: MinChunkSize,
	}
}

func WithMinChunkSize(size uint) ParIterOptions {
	return ParIterOptions{
		MinChunkSize: size,
	}
}

func WithMaxChunkSize(size uint) ParIterOptions {
	return ParIterOptions{
		MaxChunkSize: size,
	}
}

func BuildParIterOptions(opts []ParIterOptions) ParIterOptions {
	o := DefaultParIterOptions()

	for _, opt := range opts {
		o.MaxChunkSize = opt.MaxChunkSize
		o.MinChunkSize = opt.MinChunkSize
	}

	return *o
}
