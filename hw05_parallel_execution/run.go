package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	doneCh := make(chan error, len(tasks))
	defer close(doneCh)

	taskCh := make(chan Task, len(tasks))
	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)

	var isError bool
	mu := sync.Mutex{}

	wg := sync.WaitGroup{}

	defer wg.Wait()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				mu.Lock()
				if isError {
					mu.Unlock()
					break
				}
				mu.Unlock()

				doneCh <- task()
			}
		}()
	}

	var totalErrors int
	var err error
	for i := 0; i < len(tasks); i++ {
		err = <-doneCh
		if err != nil {
			totalErrors++
		}

		if totalErrors == m {
			mu.Lock()
			isError = true
			mu.Unlock()

			return ErrErrorsLimitExceeded
		}
	}

	return nil
}
