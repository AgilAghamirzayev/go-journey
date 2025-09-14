package main

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, num := range nums {
		if num%2 == 0 {
			println(num, " is even")
		} else {
			println(num, " is odd")
		}
	}

}
