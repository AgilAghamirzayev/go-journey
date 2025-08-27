package main

import (
	"fmt"
	"os"
)

func main() {
	demonstrateDefer()
	demonstratePanicAndRecover()
}

func demonstrateDefer() {
	fmt.Println("Defer statement example:")

	file, err := os.Open("example.txt")

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	fmt.Println("File opened successfully")
}

func demonstratePanicAndRecover() {
	fmt.Println("Panic and recover statement example:")

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic:", err)
		}
	}()

	fmt.Println("Before panic")

	panic("This is a panic")

	fmt.Println("After panic")
}
