package pool

import (
	"sync"
)

type Job struct {
	task func()

	jobsWg *sync.WaitGroup
	errCh  chan error
}

func NewJob(task func(), jobsWg *sync.WaitGroup, errCh chan error) Job {
	return Job{
		task:   task,
		jobsWg: jobsWg,
		errCh:  errCh,
	}
}

func SpawnJob(start, end, maxChunkSize, minChunkSize int, jobsWg *sync.WaitGroup, errCh chan error, cb func(start, end int)) {
	length := end - start
	if length > maxChunkSize && length/2 >= minChunkSize {
		mid := start + length/2

		// TODO: Maybe we can improve performance by calling
		// them inside goroutines, but we must ensure that
		// sync.WaitGroup is incremented safely before
		// parIter.Wait() is called.
		SpawnJob(start, mid, maxChunkSize, minChunkSize, jobsWg, errCh, cb)
		SpawnJob(mid, end, maxChunkSize, minChunkSize, jobsWg, errCh, cb)
		return
	}

	jobsWg.Add(1)

	job := NewJob(func() {
		cb(start, end)
	}, jobsWg, errCh)

	SpentaPool().SendJob(job)
}
