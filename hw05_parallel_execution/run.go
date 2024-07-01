package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	errs := make(chan struct{}, len(tasks))
	tc := make(chan Task)
	errCount := 0
	breakWithError := false

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&wg, errs, tc)
	}

	for _, t := range tasks {
		select {
		case <-errs:
			errCount++
			if m > 0 && errCount >= m {
				close(tc)
				breakWithError = true
				break
			} else {
				tc <- t
			}
		default:
			tc <- t
		}
		if breakWithError {
			break
		}
	}
	if !breakWithError {
		close(tc)
	}
	wg.Wait()
	if breakWithError {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(wg *sync.WaitGroup, errs chan<- struct{}, tc <-chan Task) {
	defer wg.Done()
	for t := range tc {
		err := t()
		if err != nil {
			errs <- struct{}{}
		}
	}
}
