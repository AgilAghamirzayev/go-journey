package main

import (
	"sync"
)

var count int

func main() {
	//raceCondition()
	fixedRaceCondition()
}

func fixedRaceCondition() {
	var wg2 sync.WaitGroup
	var mx sync.Mutex

	wg2.Add(20)

	for i := 0; i < 20; i++ {
		go func() {
			defer wg2.Done()
			for i := 0; i < 1000; i++ {
				mx.Lock()
				count += 1
				mx.Unlock()
			}
		}()
	}

	wg2.Wait()

	println(count)
}

func raceCondition() {
	var wg2 sync.WaitGroup
	wg2.Add(20)

	for i := 0; i < 20; i++ {
		go func() {
			defer wg2.Done()
			for i := 0; i < 1000; i++ {
				count += 1
			}
		}()
	}

	wg2.Wait()

	println(count)
}
