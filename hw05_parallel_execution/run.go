package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// максимум 0 ошибок
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	var errCount int
	taskCh := make(chan Task, len(tasks))

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case t := <-taskCh:
					if t == nil {
						return
					}
					if errCount <= m {
						if t() != nil {
							mu.Lock()
							errCount++
							mu.Unlock()
						}
					} else {
						return
					}
				default:
					return
				}
			}
		}()
	}

	for i := range tasks {
		taskCh <- tasks[i]
	}

	close(taskCh)

	wg.Wait()
	if errCount > m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
