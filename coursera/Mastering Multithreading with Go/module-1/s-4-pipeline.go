package main

import (
	"fmt"
)

// Stage 1: Generate numbers and send them to a channel
func generator(nums ...int) <-chan int {

	out := make(chan int)

	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()

	return out
}

// Stage 2: Double the input numbers and send to the next stage
func doubler(in <-chan int) <-chan int {

	out := make(chan int)

	go func() {
		defer close(out)

		for num := range in {
			out <- num * 2
		}
	}()

	return out
}

// Stage 3: Print the output numbers
func printer(in <-chan int) {

	for num := range in {
		fmt.Println(num)
	}
}

func main() {

	// Create a pipeline: generator -> doubler -> printer
	input := []int{1, 2, 3, 4, 5}

	generated := generator(input...)
	doubled := doubler(generated)
	printer(doubled)
}
