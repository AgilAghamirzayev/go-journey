package main

import (
	"log"
	"sync"
)

func printSomething(word string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(word)
}

func main() {

	var wg sync.WaitGroup

	words := []string{
		"apple",
		"banana",
		"cherry",
		"date",
		"elderberry",
		"fig",
		"grape",
		"honeydew",
		"jackfruit",
		"kiwi",
		"lemon",
	}

	wg.Add(len(words))

	for _, word := range words {
		go printSomething(word, &wg)
	}

	wg.Wait()

	wg.Add(1)
	printSomething("Finished", &wg)

}
