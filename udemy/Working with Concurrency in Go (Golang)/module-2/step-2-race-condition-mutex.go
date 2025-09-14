package main

import (
	"fmt"
	"sync"
)

var wg1 sync.WaitGroup

type Income struct {
	Source string
	Amount int
}

func main() {
	var bankBalance int
	var balance sync.Mutex

	fmt.Printf("Initial bank Balance: $%d.00\n ", bankBalance)

	incomes := []Income{
		{Source: "Main job", Amount: 100},
		{Source: "Gifts", Amount: 200},
		{Source: "Part time job", Amount: 300},
		{Source: "Investment", Amount: 400},
	}

	wg1.Add(len(incomes))

	for i, income := range incomes {
		go func(i int, income Income) {
			defer wg1.Done()

			for week := 1; week <= 52; week++ {
				balance.Lock()
				temp := bankBalance
				temp += income.Amount
				bankBalance = temp
				balance.Unlock()
				fmt.Printf("On week %d, you earned $%d.00, from %s\n", week, income.Amount, income.Source)
			}
		}(i, income)
	}

	wg1.Wait()

	fmt.Printf("Final bank Balance: $%d.00\n ", bankBalance)
}
