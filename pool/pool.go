package pool

import (
	"runtime"
	"sync"
)

var (
	_pool *Pool
	_once sync.Once
)

type Pool struct {
	jobs chan Job

	workers int
}

func SpentaPool() *Pool {
	_once.Do(func() {
		_pool = &Pool{
			jobs:    make(chan Job, 100),
			workers: runtime.NumCPU(),
		}

		for range _pool.workers {
			go _pool.worker()
		}
	})

	return _pool
}

func (p *Pool) SendJobs(jobs []Job) {
	for _, job := range jobs {
		p.jobs <- job
	}
}

func (p *Pool) worker() {
	for job := range p.jobs {
		job.task()
	}
}
