package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)
	var errCount atomic.Uint64

	defer func() {
		close(tasksCh)
		wg.Wait()
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go startWorker(&wg, tasksCh, &errCount)
	}

	for _, task := range tasks {
		tasksCh <- task

		if errCount.Load() >= uint64(m) {
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func startWorker(wg *sync.WaitGroup, tasksCh <-chan Task, errCount *atomic.Uint64) {
	defer wg.Done()

	for t := range tasksCh {
		err := t()
		if err != nil {
			errCount.Add(1)
		}
	}
}
