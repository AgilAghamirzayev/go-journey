package main

import (
	"fmt"
)

// CustomMutex is a basic mutex implementation
type CustomMutex struct {
	ch chan struct{} // Using a channel as a signaling mechanism

}

// Lock acquires the lock
func (m *CustomMutex) Lock() {
	// Create a channel with buffer size 1 to act as a lock

	m.ch = make(chan struct{}, 1)
	m.ch <- struct{}{} // Send a signal to acquire the lock

}

// Unlock releases the lock

func (m *CustomMutex) Unlock() {

	<-m.ch // Release the lock by receiving the signal

	close(m.ch) // Close the channel to release resources

}

func main() {

	var mu CustomMutex

	mu.Lock()

	defer mu.Unlock()

	// Critical section

	fmt.Println("Inside critical section")

}
