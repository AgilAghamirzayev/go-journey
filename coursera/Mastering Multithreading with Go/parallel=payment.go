package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type PaymentResult struct {
	Service string
	Status  string
	Err     error
}

// mock payment service (simulate network call)
func makePayment(ctx context.Context, service string, wg *sync.WaitGroup, ch chan<- PaymentResult) {
	defer wg.Done()

	// hər servis fərqli cavab vaxtı alsın
	delay := time.Duration(rand.Intn(4)+1) * time.Second
	fmt.Printf("[%s] Payment in progress with %d s delay\n", service, delay)

	select {
	case <-time.After(delay):
		// servis success qaytarır
		ch <- PaymentResult{
			Service: service,
			Status:  "SUCCESS",
			Err:     nil,
		}
	case <-ctx.Done():
		// timeout və ya cancel olundu
		ch <- PaymentResult{
			Service: service,
			Status:  "TIMEOUT",
			Err:     ctx.Err(),
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	services := []string{"Azerisiq", "Azerqaz", "Azersu"}
	ch := make(chan PaymentResult, len(services))
	var wg sync.WaitGroup

	for _, s := range services {
		wg.Add(1)
		go makePayment(ctx, s, &wg, ch)
	}

	// nəticələri toplamaq
	go func() {
		wg.Wait()
		close(ch)
	}()

	// statusları save etmək (simulyasiya olaraq sadəcə ekrana yazırıq)
	for result := range ch {
		if result.Err != nil {
			fmt.Printf("[%s] ERROR: %v\n", result.Service, result.Err)
		} else {
			fmt.Printf("[%s] Payment %s\n", result.Service, result.Status)
		}
	}
}
