package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" sensitive:"true"`
	FullName string `json:"full_name"`
	Age      int    `json:"age" validate:"min=18"`
}

func Describe(i any) {
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)

	// Unwrap pointers
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		fmt.Println("Describe: not a struct")
		return
	}

	fmt.Println("Type:", t.Name())
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := v.Field(i)
		fmt.Printf("- %s (%s) json=%q tags=%v value=%v\n",
			sf.Name, sf.Type, sf.Tag.Get("json"), sf.Tag, fv.Interface())
	}
}

func main() {
	u := User{ID: 1, Email: "a@b.com", FullName: "Agil", Age: 25}
	Describe(u)
}
