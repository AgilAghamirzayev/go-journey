package main

import (
	"fmt"
	"time"
)

type Subscriber struct {
	ID   string
	Chan chan string
}

type Publisher struct {
	subscribers []Subscriber
}

func (p *Publisher) Register(sub Subscriber) {
	p.subscribers = append(p.subscribers, sub)
}

func (p *Publisher) Deregister(sub Subscriber) {
	for i, s := range p.subscribers {
		if s.ID == sub.ID {
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			close(s.Chan) // close channel for cleanup
			break
		}
	}
}

func (p *Publisher) NotifyAll(message string) {
	for _, sub := range p.subscribers {
		go func(s Subscriber) {
			s.Chan <- message
		}(sub)
	}
}

func main() {
	publisher := &Publisher{}

	// Create subscribers
	sub1 := Subscriber{ID: "User1", Chan: make(chan string)}
	sub2 := Subscriber{ID: "User2", Chan: make(chan string)}

	// Register subscribers
	publisher.Register(sub1)
	publisher.Register(sub2)

	// Listen concurrently
	go func() {
		for msg := range sub1.Chan {
			fmt.Printf("Subscriber %s received: %s\n", sub1.ID, msg)
		}
	}()
	go func() {
		for msg := range sub2.Chan {
			fmt.Printf("Subscriber %s received: %s\n", sub2.ID, msg)
		}
	}()

	// Publish updates
	publisher.NotifyAll("Breaking News: Observer Pattern with channels!")
	time.Sleep(1 * time.Second)

	// Deregister one
	publisher.Deregister(sub1)

	// Publish again
	publisher.NotifyAll("Update: User1 unsubscribed, only User2 gets this.")
	time.Sleep(1 * time.Second)
}
