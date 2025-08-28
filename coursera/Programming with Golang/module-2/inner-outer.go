package main

func outer() func() int {
	a := 10
	return func() int {
		a++
		return a
	}
}

func main() {
	f := outer()
	println(f())
	println(f())
}
