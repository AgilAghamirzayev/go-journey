package main

import "fmt"

func main() {
	println(destroy())
}

func init() {
	fmt.Println("This is executed before main() function")
}

func destroy() string {
	return "This is executed after main() function"
}
