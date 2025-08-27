package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID  int
	URL string
}

type Result struct {
	ID  int
	Len int
	Err error
}

func worker(jobs <-chan Job, results chan<- Result) {
	for j := range jobs {
		// pretend: fetch URL (omitted) → len
		time.Sleep(20 * time.Millisecond)
		results <- Result{ID: j.ID, Len: len(j.URL)}
	}
}

func main() {
	const workers = 5
	jobs := make(chan Job, 100)
	results := make(chan Result, 100)

	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 1; w <= workers; w++ {
		go func(id int) { defer wg.Done(); worker(jobs, results) }(w)
	}

	// enqueue jobs
	for i, u := range []string{"https://a", "https://bb", "https://ccc"} {
		jobs <- Job{ID: i + 1, URL: u}
	}
	close(jobs)

	// close results when workers finish
	go func() { wg.Wait(); close(results) }()

	for r := range results {
		fmt.Printf("job %d len=%d err=%v\n", r.ID, r.Len, r.Err)
	}
}

func fanIn[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(chs))
	for _, ch := range chs {
		go func(c <-chan T) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
