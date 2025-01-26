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
	// количество ошибок не может быть отрицательным и нулем
	// работает и без этой проверки, но зачем вообще стартовать воркеры, если m <= 0
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)
	var errCount atomic.Int64

	defer func() {
		close(tasksCh)
		wg.Wait()
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go startWorker(&wg, tasksCh, &errCount)
	}

	for _, task := range tasks {
		if errCount.Load() >= int64(m) {
			return ErrErrorsLimitExceeded
		}
		tasksCh <- task
	}

	return nil
}

func startWorker(wg *sync.WaitGroup, tasksCh <-chan Task, errCount *atomic.Int64) {
	defer wg.Done()

	for t := range tasksCh {
		err := t()
		if err != nil {
			errCount.Add(1)
		}
	}
}
