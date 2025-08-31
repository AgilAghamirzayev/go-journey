package main

import (
	"fmt"
	"sync"
)

// Singleton struct (could hold DB connection, config, etc.)
type Singleton struct {
	Value string
}

var instance *Singleton
var once sync.Once

// GetInstance ensures only one instance is created
func GetInstance() *Singleton {
	once.Do(func() {
		fmt.Println("Creating new Singleton instance...")
		instance = &Singleton{Value: "I am the only instance"}
	})
	return instance
}

func main() {
	// Simulate multiple calls
	s1 := GetInstance()
	s2 := GetInstance()
	s3 := Singleton{Value: "I am the second instance"}

	fmt.Println("s1 value:", s1.Value)
	fmt.Println("s2 value:", s2.Value)

	// Verify both are same instance
	if s1 == s2 {
		fmt.Println("Both variables point to the same Singleton instance ✅")
	}

	fmt.Println(s3)
}
