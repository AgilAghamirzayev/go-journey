package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Main function STARTED")

	go func() {
		fmt.Println("Anonymous goroutine STARTED")
		time.Sleep(2 * time.Second)
		fmt.Println("Anonymous goroutine FINISHED")
	}()

	fmt.Println("Main function CONTINUE")
	time.Sleep(3 * time.Second)
	fmt.Println("Main function FINISHED")
}
