package main

import "fmt"

func main() {
	fruitCount1 := map[string]int{}
	fruitCount2 := make(map[string]int)

	fmt.Println(fruitCount1)
	fmt.Println(fruitCount2)

	fruitCount1["apple"] = 1
	fruitCount1["banana"] = 2
	fruitCount1["orange"] = 3
	fruitCount1["grape"] = 4
	fruitCount1["mango"] = 5
	fruitCount1["pineapple"] = 6
	fruitCount1["strawberry"] = 7
	fruitCount1["watermelon"] = 8
	fruitCount1["cherry"] = 9
	fruitCount1["lemon"] = 10

	fruitCount2["apple"] = 1
	fruitCount2["banana"] = 2
	fruitCount2["orange"] = 3
	fruitCount2["grape"] = 4
	fruitCount2["mango"] = 5
	fruitCount2["pineapple"] = 6
	fruitCount2["strawberry"] = 7
	fruitCount2["watermelon"] = 8
	fruitCount2["cherry"] = 9
	fruitCount2["lemon"] = 10

	fmt.Println(fruitCount1)
	fmt.Println(fruitCount2)

	delete(fruitCount1, "apple")
	fmt.Println(fruitCount1)

	count, exist := fruitCount1["lemon"]
	fmt.Println(count, exist)
}
