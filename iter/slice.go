package iter

import (
	"sync"

	"github.com/rouzbehsbz/spenta/pool"
)

func NewSliceParIter[V any](slice *[]V, cb func(start, end int), opts ...ParIterOptions) *ParIter {
	options := BuildParIterOptions(opts)
	length := len(*slice)

	parIter := NewParIter()

	pool.SpawnJob(0, length, int(options.MaxChunkSize), int(options.MinChunkSize), parIter.jobsWg, parIter.errCh, func(start, end int) {
		cb(start, end)
	})

	return parIter
}

func SliceParForEach[V any](slice *[]V, cb func(i int, v V), opts ...ParIterOptions) *ParIter {
	p := NewSliceParIter[V](slice, func(start, end int) {
		for i := start; i < end; i++ {
			cb(i, (*slice)[i])
		}
	}, opts...)

	p.postJobsDone()

	return p
}

func SliceParMap[V any](slice *[]V, cb func(i int, v V) V, opts ...ParIterOptions) *ParIter {
	p := NewSliceParIter[V](slice, func(start, end int) {
		for i := start; i < end; i++ {
			(*slice)[i] = cb(i, (*slice)[i])
		}
	}, opts...)

	p.postJobsDone()

	return p
}

func SliceParFilter[V any](slice *[]V, cb func(i int, v V) bool, opts ...ParIterOptions) *ParIter {
	merge := []V{}
	mu := &sync.Mutex{}

	p := NewSliceParIter(slice, func(start, end int) {
		local := make([]V, 0, end-start)

		copy(local, (*slice)[start:end])

		for i := start; i < end; i++ {
			if cb(i, (*slice)[i]) {
				local = append(local, (*slice)[i])
			}
		}

		mu.Lock()
		merge = append(merge, local...)
		mu.Unlock()
	}, opts...)

	go func() {
		<-p.jobsDoneCh
		*slice = merge
		p.postJobsDone()
	}()

	return p
}
