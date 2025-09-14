package main

import "fmt"

func main() {

	cards := []string{getCard(1), getCard(2), getCard(3)}
	cards = append(cards, getCard(4))

	for i, card := range cards {
		println("Card", i+1, "is", card)
	}

}

func getCard(a int) string {
	return fmt.Sprintf("Card %d", a)
}
