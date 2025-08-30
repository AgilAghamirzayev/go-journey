package main

import (
	"fmt"
	"time"
)

func task(resultCh chan string, errCh chan error) {
	defer func() {
		if r := recover(); r != nil {
			errCh <- fmt.Errorf("panic: %v", r)
		}
	}()

	// Təsadüfi panic
	panic("something went wrong!")
	// normalda: resultCh <- "Task completed"
}

func main() {
	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go task(resultCh, errCh)

	select {
	case res := <-resultCh:
		fmt.Println("Got result:", res)
	case err := <-errCh:
		fmt.Println("Recovered from panic:", err)
	case <-time.After(2 * time.Second):
		fmt.Println("Timeout, no response")
	}
}
