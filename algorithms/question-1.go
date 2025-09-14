package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan *int, 4)
	array := []int{1, 2, 3, 4}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for _, value := range array {
			v := value
			ch <- &v
		}
		defer close(ch)
	}()

	go func() {
		defer wg.Done()
		for value := range ch {
			fmt.Println(*value) // what will be printed here?
		}
	}()

	wg.Wait()
}
