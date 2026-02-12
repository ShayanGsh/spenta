package iter

import (
	"sync"

	"github.com/rouzbehsbz/spenta/pool"
)

func NewSliceParIter[T comparable](slice *[]T, cb func(i int)) *ParIter {
	len, chunkSize, chunkCount := SliceChunk(slice, pool.SpentaPool().Workers)

	wg := &sync.WaitGroup{}
	wg.Add(chunkCount)

	jobs := pool.NewSliceJobs(len, chunkCount, chunkSize, wg, func(i int) {
		cb(i)
	})

	go pool.SpentaPool().SendJobs(jobs)

	return &ParIter{
		wg: wg,
	}
}

func SliceParForEach[T comparable](slice *[]T, cb func(e T)) *ParIter {
	return NewSliceParIter[T](slice, func(i int) {
		cb((*slice)[i])
	})
}

func SliceParMap[T comparable](slice *[]T, cb func(e T) T) *ParIter {
	return NewSliceParIter[T](slice, func(i int) {
		(*slice)[i] = cb((*slice)[i])
	})
}

func SliceParFilter[T comparable](slice *[]T, cb func(e T) bool) *ParIter {
	return NewSliceParIter[T](slice, func(i int) {
		s := *slice
		j := 0
		for _, v := range s {
			if cb(v) {
				s[j] = v
				j++
			}
		}
		*slice = s[:j]
	})
}
