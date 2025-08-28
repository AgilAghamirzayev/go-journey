package main

import "fmt"

type Task struct {
	Title string
	Done  bool
}

var tasks []Task

func addTask(title string) {
	tasks = append(tasks, Task{title, false})
}

func removeTask(index int) {
	tasks = append(tasks[:index], tasks[index+1:]...)
}

func displayTasks() {
	fmt.Println("Your Tasks:")
	for i, task := range tasks {
		fmt.Println(i, task.Title, task.Done)
	}
}

func markTaskDone(index int) {
	tasks[index].Done = true
}

func main() {
	addTask("Learn Golang")
	addTask("Build a project")
	addTask("Deploy to production")
	displayTasks()

	markTaskDone(0)
	displayTasks()

	removeTask(0)
	displayTasks()
}

/*
1. **Define a Task Structure:**
2. **Create a Slice to Hold Tasks:**
3. **Add Tasks:**
4. **Remove Tasks:**
5. **Display Tasks:**
6. **Example Usage:**
   Here’s how you might use these functions in your main program.

   ```go
   func main() {
       addTask("Learn Golang")
       addTask("Build a project")
       displayTasks()

       removeTask(0) // Remove the first task
       displayTasks()
   }
   ```

*/
