package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		n = 1
	}

	taskCh := make(chan Task, n)
	var errCounter int64
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if atomic.LoadInt64(&errCounter) >= int64(m) {
					continue
				}
				err := task()
				if err != nil {
					atomic.AddInt64(&errCounter, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt64(&errCounter) >= int64(m) {
			break
		}
		taskCh <- task
	}
	close(taskCh)
	wg.Wait()

	if atomic.LoadInt64(&errCounter) >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
