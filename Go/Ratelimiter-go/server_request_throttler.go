package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func fn3() {
	fmt.Println("Starting task")
	delay := time.Duration(rand.Float64()*2000) * time.Millisecond
	time.Sleep(delay)
	fmt.Println("Completed task")
}

func limitConcurrency(tasks []func(), concurrency int) {
	var index int
	var mu sync.Mutex

	worker := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			mu.Lock()
			if index >= len(tasks) {
				mu.Unlock()
				return
			}
			current := index
			index++
			mu.Unlock()
			tasks[current]()
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
}

func demoConcurrencyThrottler() {
	tasks := make([]func(), 24)
	for i := range tasks {
		tasks[i] = fn3
	}
	limitConcurrency(tasks, 3)
}
