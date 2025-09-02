package main

import "fmt"

func main() {
	ch := make(chan string, 2)

	ch <- "Hello"
	msg := <-ch

	ch <- "World"
	msg1 := <-ch

	ch <- "!"
	msg2 := <-ch

	println(msg)
	println(msg1)
	println(msg2)

	close(ch)
	for v := range ch {
		println(v)
	}

	v, ok := <-ch

	println(v, ok)

}

func producer(out chan<- int) { // send-only
	for i := 0; i < 3; i++ {
		out <- i
	}
	close(out)
}

func consumer(in <-chan int) { // receive-only
	for v := range in {
		fmt.Println(v)
	}
}
