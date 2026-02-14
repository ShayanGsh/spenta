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
	jobsWg *sync.WaitGroup
	errors []error

	errCh          chan error
	jobsDoneCh     chan struct{}
	postJobsDoneCh chan struct{}

	doneOnce sync.Once
}

func NewParIter() *ParIter {
	p := &ParIter{
		errors:         []error{},
		jobsWg:         &sync.WaitGroup{},
		jobsDoneCh:     make(chan struct{}),
		postJobsDoneCh: make(chan struct{}),
		errCh:          make(chan error),
		doneOnce:       sync.Once{},
	}

	go func() {
		for err := range p.errCh {
			p.errors = append(p.errors, err)
		}
	}()

	go func() {
		p.jobsWg.Wait()
		close(p.jobsDoneCh)
	}()

	return p
}

func (p *ParIter) Wait() error {
	p.doneOnce.Do(func() {
		<-p.jobsDoneCh
		<-p.postJobsDoneCh
		close(p.errCh)
	})

	return errors.Join(p.errors...)
}

func (p *ParIter) postJobsDone() {
	close(p.postJobsDoneCh)
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
		if opt.MaxChunkSize > 0 {
			o.MaxChunkSize = opt.MaxChunkSize
		}
		if opt.MinChunkSize > 0 {
			o.MinChunkSize = opt.MinChunkSize
		}
	}

	return *o
}
