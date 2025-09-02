package main

import "fmt"

type Task struct {
	ID int
	// Add other fields as needed
}

func workerr(id int, tasks <-chan Task, results chan<- string) {
	for task := range tasks {
		// Process the task (e.g., simulate work)
		results <- fmt.Sprintf("Worker id: %d Processed task ID: %d", id, task.ID)
	}
}

func main() {
	const numWorkers = 4
	tasks := make(chan Task, 10) // Channel for tasks
	results := make(chan string) // Channel for results

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go workerr(i, tasks, results)
	}

	// Add tasks to the channel
	for i := 1; i <= 10; i++ {
		tasks <- Task{ID: i}
	}
	close(tasks) // Close the tasks channel when done

	// Collect results
	for i := 0; i < 10; i++ {
		fmt.Println(<-results)
	}
}
