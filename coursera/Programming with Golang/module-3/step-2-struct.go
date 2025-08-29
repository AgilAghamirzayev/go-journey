package main

import "fmt"

type Contact struct {
	Email string
	Phone string
}

type Person struct {
	FirstName   string
	LastName    string
	ContactInfo Contact
}

func main() {
	alice := Person{
		FirstName: "Alice",
		LastName:  "Smith",
		ContactInfo: Contact{
			Email: "alice@gmail.com",
			Phone: "+254779",
		},
	}

	fmt.Println(alice)
}
