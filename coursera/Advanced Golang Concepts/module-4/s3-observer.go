package main

import "fmt"

type Observer interface {
	Update(message string)
}

type EmailSubscriber struct {
	ID string
}

func (e *EmailSubscriber) Update(message string) {
	fmt.Printf("EmailSubscriber %s received update: %s\n", e.ID, message)
}

type Subject interface {
	Register(o Observer)
	Deregister(o Observer)
	NotifyAll(message string)
}

type NewsPublisher struct {
	observers []Observer
}

func (n *NewsPublisher) Register(o Observer) {
	n.observers = append(n.observers, o)
}

func (n *NewsPublisher) Deregister(o Observer) {
	for i, observer := range n.observers {
		if observer == o {
			n.observers = append(n.observers[:i], n.observers[i+1:]...)
			break
		}
	}
}

func (n *NewsPublisher) NotifyAll(message string) {
	for _, observer := range n.observers {
		observer.Update(message)
	}
}

func main() {
	// Create publisher
	newsPublisher := &NewsPublisher{}

	// Create subscribers
	s1 := &EmailSubscriber{ID: "User1"}
	s2 := &EmailSubscriber{ID: "User2"}

	// Register subscribers
	newsPublisher.Register(s1)
	newsPublisher.Register(s2)

	// Notify all
	newsPublisher.NotifyAll("Breaking News: Go design patterns explained!")

	// Deregister one subscriber
	newsPublisher.Deregister(s1)

	// Notify remaining
	newsPublisher.NotifyAll("Update: Observer pattern implemented successfully.")
}
