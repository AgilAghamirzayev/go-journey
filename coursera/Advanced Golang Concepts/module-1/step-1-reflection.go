package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name   string
	Age    int
	Gender string
}

func main() {
	p := Person{
		Name:   "John",
		Age:    25,
		Gender: "Male",
	}

	structType := reflect.TypeOf(p)
	fmt.Println("Type of p:", structType)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fmt.Println("Field", i, ":", field.Name, field.Type)
	}
}
