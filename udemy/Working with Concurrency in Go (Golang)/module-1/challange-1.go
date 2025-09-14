package main

import "sync"

var msg string
var wg sync.WaitGroup

func updateMessage(newMsg string) {
	defer wg.Done()
	msg = newMsg
}

func printMessage() {
	println(msg)
}

func main() {

	msg = "Hello World"

	wg.Add(1)
	go updateMessage("Hello Universe")
	wg.Wait()
	printMessage()

	wg.Add(1)
	go updateMessage("Hello Galaxy")
	wg.Wait()
	printMessage()

	wg.Add(1)
	go updateMessage("Hello World")
	wg.Wait()
	printMessage()

}
