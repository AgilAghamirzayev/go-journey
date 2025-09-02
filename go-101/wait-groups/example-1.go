package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func worker(id int) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go worker(i)
	}

	wg.Done()
	wg.Wait()
}
