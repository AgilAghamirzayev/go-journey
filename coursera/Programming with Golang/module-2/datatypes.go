package main

import "fmt"

func main() {

	mySlice := []int{1, 2, 3}
	fmt.Println(mySlice)

	mySlice = append(mySlice, 4)
	fmt.Println(mySlice)

}
