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

	var errCounter int64
	var idx int64 = -1
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if atomic.LoadInt64(&errCounter) >= int64(m) {
					return
				}
				i := atomic.AddInt64(&idx, 1)
				if i >= int64(len(tasks)) {
					return
				}
				err := tasks[i]()
				if err != nil {
					atomic.AddInt64(&errCounter, 1)
				}
			}
		}()
	}

	wg.Wait()

	if atomic.LoadInt64(&errCounter) >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
