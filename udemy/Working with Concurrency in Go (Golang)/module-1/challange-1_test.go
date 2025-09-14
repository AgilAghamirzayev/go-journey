package main

import "testing"

func Test_updateMessage(t *testing.T) {
	wg.Add(1)
	go updateMessage("Hello")
	wg.Wait()
	if msg != "Hello" {
		t.Errorf("Expected msg to be 'Hello', got %s", msg)
	}
}
