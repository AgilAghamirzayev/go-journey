package main

import (
	"fmt"
	"time"
)

func main() {
	in := make(chan int)
	out := make(chan int)

	select {
	case v := <-in:
		fmt.Println("got", v)
	case out <- 99:
		fmt.Println("sent 99")
	case <-time.After(200 * time.Millisecond):
		fmt.Println("timeout")
	default:
		fmt.Println("nothing ready (non-blocking)")
	}

}
