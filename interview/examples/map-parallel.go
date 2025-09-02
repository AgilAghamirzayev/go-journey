package main

import (
	"fmt"
	"time"
)

func main() {
	m := make(map[string]int)

	for i := 1; i <= 10; i++ {
		go func(id int) {
			m[fmt.Sprintf("goroutine-%d", id)] = id
		}(i)
	}

	for s := range m {
		go func(s string) {
			fmt.Println(s)
		}(s)
	}

	time.Sleep(time.Second * 2)

	fmt.Println(m)
}
