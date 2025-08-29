package main

import "fmt"

type MyError struct {
	Code    int
	Message string
}

func (e MyError) Error() string {
	return e.Message
}

func main() {
	err := MyError{Code: 404, Message: "Not found"}

	fmt.Println(err.Error())
}
