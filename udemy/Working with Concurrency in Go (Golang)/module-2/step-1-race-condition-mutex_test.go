package main

import (
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	msg = "Hello World"

	wg.Add(2)
	go updateMessage("Good bye", &sync.Mutex{})
	go updateMessage("Good bye again", &sync.Mutex{})
	wg.Wait()

	if msg != "Good bye again" && msg != "Good bye" {
		t.Error("Message should not have been updated")
	}
}
