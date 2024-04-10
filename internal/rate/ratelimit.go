package rate

import (
	"fmt"
	"sync"
)

type JobFunc[T any] func(int) T

type RateLimiter[T any] struct {
	jobFunction  JobFunc[T]
	dataChan     chan T
	queueChan    chan int
	limit        int
	numberOfJobs int
	waitGroup    sync.WaitGroup
}

func New[T any](limit, numberOfJobs int, jobFn JobFunc[T]) *RateLimiter[T] {
	return &RateLimiter[T]{
		jobFunction:  jobFn,
		limit:        limit,
		numberOfJobs: numberOfJobs,
		queueChan:    make(chan int, limit),
		dataChan:     make(chan T, numberOfJobs),
		waitGroup:    sync.WaitGroup{},
	}
}

func (r *RateLimiter[T]) Spawn() *RateLimiter[T] {
	for i := 0; i < r.limit; i++ {
		r.waitGroup.Add(1)
		go func() {
			defer r.waitGroup.Done()
			for j := range r.queueChan {
				r.dataChan <- r.jobFunction(j)
			}
		}()
	}

	return r
}

func (r *RateLimiter[T]) Run() chan T {
	fmt.Println("feeding...")
	for i := 0; i < r.numberOfJobs; i++ {
		r.queueChan <- i
	}
	close(r.queueChan)
	r.waitGroup.Wait()
	close(r.dataChan)
	return r.dataChan
}
