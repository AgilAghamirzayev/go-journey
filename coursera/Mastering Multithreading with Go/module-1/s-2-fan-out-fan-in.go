package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// worker function
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Simulate work
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		result := job * 2
		fmt.Printf("Worker %d processed job %d -> %d\n", id, job, result)
		results <- result
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	jobs := make(chan int, 10)
	results := make(chan int, 10)

	var wg sync.WaitGroup

	// FAN-OUT: start 3 worker goroutines
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// send 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs) // no more jobs

	// Wait for all workers to finish (in background)
	go func() {
		wg.Wait()
		close(results)
	}()

	// FAN-IN: collect results into one stream
	for res := range results {
		fmt.Println("Result:", res)
	}
}
