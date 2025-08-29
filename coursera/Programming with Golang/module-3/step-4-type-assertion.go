package main

import "fmt"

func main() {

	var x interface{}
	x = 42

	if value, ok := x.(int); ok {
		fmt.Println("x is an int:", value)
	} else {
		fmt.Println("x is not an int")
	}

	x = "Hello, World"

	if value, ok := x.(string); ok {
		fmt.Println("x is a string:", value)
	} else {
		fmt.Println("x is not a string")
	}

}
