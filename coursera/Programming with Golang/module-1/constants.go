package main

import "fmt"

func main() {
	const a = 10

	const (
		Sunday    = 0
		Monday    = 1
		Tuesday   = 2
		Wednesday = 3
		Thursday  = 4
		Friday    = 5
		Saturday  = 6
	)

	const (
		_ = iota
		January
		February
		March
		April
		May
		June
		July
		August
	)

	fmt.Print(a)
	fmt.Print(Sunday)
	fmt.Print(January)

}
