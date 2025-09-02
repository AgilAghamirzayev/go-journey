package main

import "fmt"

func main() {
	array := []int{1, 2, 3, 4}
	newArray := make([]*int, 4)

	for i, value := range array {
		newArray[i] = &value
	}

	for _, value := range newArray {
		fmt.Println(*value)
	}
}
