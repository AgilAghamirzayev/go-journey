package main

func main() {
	x := 10

	if x > 6 {
		x := 11
		println(x)
	}

	println(x)
	println(Exported)
	println(unexported)
}
