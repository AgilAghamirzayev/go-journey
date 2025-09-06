package main

import (
	"fmt"
	"sync"
)

// Transaction represents a financial transaction.
type Transaction struct {
	ID          string  // Unique identifier for the transaction
	Description string  // Description of the transaction
	Amount      float64 // Transaction amount
}

// Ledger represents a ledger for tracking transactions and maintaining the current balance.
type Ledger struct {
	mutex        sync.Mutex
	transactions []Transaction
	balance      float64
}

// AddTransaction adds a new transaction to the ledger.

func (l *Ledger) AddTransaction(transaction Transaction) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.transactions = append(l.transactions, transaction)
	l.balance += transaction.Amount
}

// GetBalance returns the current balance in the ledger.

func (l *Ledger) GetBalance() float64 {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.balance
}

// PrintLedger prints the transactions and current balance in the ledger.

func (l *Ledger) PrintLedger() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	fmt.Println("Ledger Transactions:")

	for _, transaction := range l.transactions {
		fmt.Printf("ID: %s, Description: %s, Amount: %.2f\n", transaction.ID, transaction.Description, transaction.Amount)
	}

	fmt.Printf("Current Balance: %.2f\n", l.balance)
}

func main() {

	// Example usage of the Ledger structure
	myLedger := Ledger{}

	// Add transactions to the ledger
	myLedger.AddTransaction(Transaction{"123456", "Salary", 5000.0})
	myLedger.AddTransaction(Transaction{"789012", "Rent Payment", -1200.0})
	myLedger.AddTransaction(Transaction{"345678", "Utilities", -200.0})

	// Print the ledger
	myLedger.PrintLedger()

	// Get the current balance
	fmt.Printf("Final Balance: %.2f\n", myLedger.GetBalance())
}
