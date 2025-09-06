package main

import (
	"fmt"
	"sync"
)

// Barrier represents a reusable synchronization barrier
type Barrier struct {
	count int
	wg    sync.WaitGroup
	mu    sync.Mutex
}

// NewBarrier creates a new instance of Barrier

func NewBarrier(n int) *Barrier {
	return &Barrier{
		count: n,
	}
}

// Wait waits for all goroutines to reach the barrier before proceeding
func (b *Barrier) Wait() {

	b.mu.Lock()
	b.count--

	if b.count == 0 {

		// If count becomes 0, reset counter and signal all waiting goroutines
		b.count = 3 // Set it according to the number of goroutines to synchronize
		b.mu.Unlock()
		b.wg.Done() // Signal completion to all waiting goroutines
		return
	}

	b.mu.Unlock()
	b.wg.Add(1) // Add one more goroutine to wait
	b.wg.Wait() // Wait until signaled by the last goroutine
}

func main() {

	barrier := NewBarrier(3) // Create a barrier for 3 goroutines

	for i := 0; i < 3; i++ {
		go func(id int) {
			for j := 0; j < 3; j++ {
				fmt.Printf("Goroutine %d - Task %d\n", id, j)
				barrier.Wait() // Synchronize at the barrier
			}
		}(i)
	}

	// Wait for all goroutines to complete
	barrier.Wait()
	fmt.Println("All goroutines finished")
}
