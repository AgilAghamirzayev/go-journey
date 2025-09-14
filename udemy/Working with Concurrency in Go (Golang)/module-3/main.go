package module_3

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

var philosophers = []Philosopher{
	{name: "Plato", rightFork: 0, leftFork: 4},
	{name: "Socrates", rightFork: 1, leftFork: 0},
	{name: "Aristotle", rightFork: 2, leftFork: 1},
	{name: "Pascal", rightFork: 3, leftFork: 2},
	{name: "Locke", rightFork: 4, leftFork: 3},
}

func main() {

}
