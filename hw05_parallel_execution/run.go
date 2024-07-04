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
	var wg sync.WaitGroup
	tc := make(chan Task)
	var errC atomic.Uint64
	breakWithError := false
	flag := m <= 0
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(&wg, tc, &errC)
	}

	if flag {
		for _, t := range tasks {
			tc <- t
		}
		close(tc)
	} else {
		for _, t := range tasks {
			if errC.Load() < uint64(m) {
				tc <- t
			} else {
				close(tc)
				breakWithError = true
				break
			}
		}
		if !breakWithError {
			close(tc)
		}
	}
	wg.Wait()
	if breakWithError {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(wg *sync.WaitGroup, tc <-chan Task, errC *atomic.Uint64) {
	defer wg.Done()
	for t := range tc {
		err := t()
		if err != nil {
			errC.Add(1)
		}
	}
}
