package main

import (
	"fmt"
	"time"
)

func riskyWorker(errCh chan error) {
	defer func() {
		if r := recover(); r != nil {
			errCh <- fmt.Errorf("worker panicked: %v", r)
		}
	}()

	// Panikaya səbəb olacaq
	panic("database connection lost!")
}

func main() {
	errCh := make(chan error)

	go riskyWorker(errCh)

	select {
	case err := <-errCh:
		fmt.Println("Error caught:", err)
	case <-time.After(2 * time.Second):
		fmt.Println("No error, worker finished fine")
	}
}
