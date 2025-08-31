package main

import (
	"fmt"
	"reflect"
)

type Person2 struct {
	Name   string
	Age    int
	Gender string
}

func main() {
	p := Person2{
		Name:   "John",
		Age:    25,
		Gender: "Male",
	}

	fmt.Println(p)

	structType := reflect.TypeOf(p)
	fmt.Println("Type of p:", structType)

	valueOfPerson := reflect.ValueOf(p)
	nameField := valueOfPerson.FieldByName("Name")

	fmt.Println("Name field:", nameField)
	fmt.Println("Name field:", nameField.Interface())

	//ageField := valueOfPerson.FieldByName("Age")
	//ageField.Set(35)
	//fmt.Println("Age field:", ageField.Interface())

}
