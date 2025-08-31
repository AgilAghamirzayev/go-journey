package main

import (
	"fmt"
	"time"
)

type CustomError struct {
	Code      int
	Message   string
	Timestamp time.Time
	RequestID string
}

func (e CustomError) Error() string {
	return fmt.Sprintf("Error: %d: %s  (RequestID: %s, Timestamp: %s)", e.Code, e.Message, e.RequestID, e.Timestamp.String())
}

func main() {
	e := CustomError{Code: 404, Message: "Not found", Timestamp: time.Now(), RequestID: "1234567890"}
	fmt.Println(e.Error())
}
