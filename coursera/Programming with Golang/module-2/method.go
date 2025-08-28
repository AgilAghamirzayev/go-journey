package main

import "fmt"

type Author struct {
	name     string
	branch   string
	articles string
	salary   int
}

func (a Author) show() {
	fmt.Println("Author name: ", a.name)
	fmt.Println("Author branch: ", a.branch)
	fmt.Println("Author salary: ", a.salary)
}

func main() {
	a := Author{"Ali", "Computer Science", "100", 100000}
	a.show()

	add := func(a, b int) int {
		return a + b
	}

	result := add(2, 3)
	fmt.Println("Result:", result)
}
