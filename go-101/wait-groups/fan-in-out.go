package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type result struct {
	name string
	val  string
	err  error
}

func main() {
	services := map[string]func(ctx context.Context) (string, error){
		"Azerisiq": func(ctx context.Context) (string, error) {
			time.Sleep(1200 * time.Millisecond)
			return "paid: azerisiq", nil
		},
		"Azerqaz": func(ctx context.Context) (string, error) {
			time.Sleep(2 * time.Second)
			return "paid: azerqaz", nil
		},
		"Azersu": func(ctx context.Context) (string, error) {
			time.Sleep(6 * time.Second)
			return "paid: azersu", nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(services))

	results := make(chan result, len(services))

	for name, fn := range services {
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				results <- result{name: name, err: ctx.Err()}
			default:
				v, err := fn(ctx)
				results <- result{name: name, val: v, err: err}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var ok, failed []string
	for r := range results {
		if r.err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", r.name, r.err))
			continue
		}
		ok = append(ok, fmt.Sprintf("%s → %s", r.name, r.val))
	}

	fmt.Println("Success:", ok)
	fmt.Println("Failed :", failed)
}
