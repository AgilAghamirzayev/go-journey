package main

import (
	"log"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(newMsg string, m *sync.Mutex) {
	defer wg.Done()

	m.Lock()
	msg = newMsg
	m.Unlock()
}

func main() {
	msg = "Hello World"

	var mutex sync.Mutex

	wg.Add(2)
	go updateMessage("Hello Universe", &mutex)
	go updateMessage("Hello Galaxy", &mutex)
	wg.Wait()

	log.Println(msg)
}
