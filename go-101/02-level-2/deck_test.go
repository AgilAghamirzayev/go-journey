package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d) != 52 {
		t.Errorf("Expected deck length of 52, but got %d", len(d))
	}

	if d[0] != "Ace of Spades" {
		t.Errorf("Expected first card to be 'Ace of Spades', but got '%s'", d[0])
	}

	if d[len(d)-1] != "King of Clubs" {
		t.Errorf("Expected last card to be 'King of Clubs', but got '%s'", d[len(d)-1])
	}
}

func TestDeal(t *testing.T) {
	d := newDeck()
	handSize := 5
	hand, remainingDeck := deal(d, handSize)

	if len(hand) != handSize {
		t.Errorf("Expected hand size of %d, but got %d", handSize, len(hand))
	}

	if len(remainingDeck) != len(d)-handSize {
		t.Errorf("Expected remaining deck size of %d, but got %d", len(d)-handSize, len(remainingDeck))
	}
}

func TestToString(t *testing.T) {
	d := newDeck()
	expectedString := strings.Join(d, ",")
	deckString := d.toString()

	if deckString != expectedString {
		t.Errorf("Expected %s, but got %s", expectedString, deckString)
	}
}

func TestSaveToFileAndNewDeckFromFile(t *testing.T) {
	filename := "_decktesting"
	defer os.Remove(filename)

	deck := newDeck()
	err := deck.saveToFile(filename)
	if err != nil {
		t.Errorf("Error saving deck to file: %v", err)
	}

	loadedDeck := newDeckFromFile(filename)

	if len(loadedDeck) != 52 {
		t.Errorf("Expected deck length of 52, but got %d", len(loadedDeck))
	}

	if !reflect.DeepEqual(deck, loadedDeck) {
		t.Errorf("Expected loaded deck to be equal to saved deck")
	}
}

func TestShuffle(t *testing.T) {
	d := newDeck()
	d.shuffle()

	if len(d) != 52 {
		t.Errorf("Expected 52 cards in deck, but got %d", len(d))
	}

	unshuffledDeck := newDeck()
	if reflect.DeepEqual(d, unshuffledDeck) {
		t.Errorf("Expected shuffled deck to be different from the original deck, but they are the same")
	}
}
