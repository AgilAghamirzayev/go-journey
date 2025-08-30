package main

import "fmt"

func main() {

	ch := make(chan int)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
			}
		}()

		panic("Something went wrong!")

		ch <- 42

	}()

	select {
	case value := <-ch:
		fmt.Println("Received from channel:", value)
	}
	fmt.Println("Program continues after recovery")
}
