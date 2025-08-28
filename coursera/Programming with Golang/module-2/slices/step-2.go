package main

import "fmt"

func main() {
	// Slice with values
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("nums:", nums)

	fmt.Println(nums[1:4]) // [2 3 4]
	fmt.Println(nums[:3])  // [1 2 3]
	fmt.Println(nums[2:])  // [3 4 5]

}
