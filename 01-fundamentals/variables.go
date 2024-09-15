package main

import "fmt"

func main() {
	a := 12
	println(a)

	b, c, d, _, f := 0, 1, 2, 4, "happiness"
	println(b, c, d, f)

	var g int
	println(g)

	var isCold bool
	println(isCold)

	adams := 43
	println("Adams is", adams, "years old.")
	fmt.Printf("Adams is %d years old.\n", adams)

	fmt.Printf("43 as binary: %b \n", adams)
	fmt.Printf("44 as hexadecimal: %x \n", adams)

	x := 54
	z := 1152167890277787419
	y := 55.11

	println()
	fmt.Printf("x: %d, y: %.2f\n", x, y)
	fmt.Printf("x: %T, y: %T\n", x, y)
	fmt.Printf("z: %d, type: %T\n", z, z)

	var name string = "Response from server"
	fmt.Printf("name: %s\n", name)

}
