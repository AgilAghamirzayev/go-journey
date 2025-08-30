package main

import (
	"fmt"
	"time"
)

func safeWorker(id int, errCh chan error) {
	defer func() {
		if r := recover(); r != nil {
			errCh <- fmt.Errorf("worker %d panicked: %v", id, r)
		}
	}()

	if id == 2 {
		panic("critical bug in worker 2")
	}

	time.Sleep(1 * time.Second)
	fmt.Printf("Worker %d finished\n", id)
	errCh <- nil
}

func main() {
	errCh := make(chan error, 3)

	for i := 1; i <= 3; i++ {
		go safeWorker(i, errCh)
	}

	for i := 0; i < 3; i++ {
		select {
		case err := <-errCh:
			if err != nil {
				fmt.Println("Handled:", err)
			} else {
				fmt.Println("Worker completed successfully")
			}
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout waiting for worker")
		}
	}
}
