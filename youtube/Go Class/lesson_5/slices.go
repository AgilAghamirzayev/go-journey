package main

import "log"

func main() {

	var a []int
	log.Println(a)

	b := make([]int, 0)
	log.Println(b)

	x := []int{1, 2, 3, 4}
	log.Println(x)

	y := x
	log.Println(y)

	x[0] = 99
	log.Println(x)
	log.Println(y)

	y[1] = 11

	log.Println(x)
	log.Println(y)

}
