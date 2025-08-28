package main

import "fmt"

func main() {
	// Empty slice
	var s []int
	fmt.Println("empty:", s)

	// Slice with values
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("nums:", nums)

	// Slice with make (len=3, cap=5)
	s2 := make([]int, 3, 5)
	fmt.Println("make:", s2, "len:", len(s2), "cap:", cap(s2))

	s2 = append(s2, 6, 7, 8, 9, 10)
	fmt.Println("append:", s2)

	nums1 := []int{10, 20, 30}
	fmt.Println(nums1[0]) // 10
	fmt.Println(nums1[2]) // 30

}
