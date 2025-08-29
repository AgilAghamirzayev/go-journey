package main

import "fmt"

func describeType(v interface{}) {

	switch x := v.(type) {

	case int:
		fmt.Printf("It's an integer: %d\n", x)

	case string:
		fmt.Printf("It's a string: %s\n", x)

	default:
		fmt.Printf("It's something else: %v\n", x)

	}

}

func main() {

	describeType(42)
	describeType("Hello, World")
	describeType(3.14)

}
