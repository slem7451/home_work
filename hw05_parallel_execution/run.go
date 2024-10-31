package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	taskCh := make(chan Task, len(tasks))
	errorCount := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < n; i++ {
		wg.Add(1)
		go run(taskCh, &errorCount, m, &wg, &mu)
	}

	for _, t := range tasks {
		mu.Lock()
		if errorCount >= m {
			mu.Unlock()
			close(taskCh)
			return ErrErrorsLimitExceeded
		}
		mu.Unlock()
		taskCh <- t
	}

	close(taskCh)
	wg.Wait()

	if errorCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func run(taskCh <-chan Task, errorCount *int, m int, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for {
		mu.Lock()
		if *errorCount >= m {
			mu.Unlock()
			return
		}
		mu.Unlock()

		task, ok := <-taskCh
		if !ok {
			return
		}
		err := task()
		if err != nil {
			mu.Lock()
			*errorCount++
			mu.Unlock()
		}
	}
}
