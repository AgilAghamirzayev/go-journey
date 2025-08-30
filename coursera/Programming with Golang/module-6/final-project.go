package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Task struct {
	Description string
	Completed   bool
}

var taskList []Task

const filename = "tasks.json"

// ------------------ File Operations ------------------
func saveToFile(taskList []Task) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(taskList); err != nil {
		fmt.Println("Error encoding tasks:", err)
	}
	fmt.Println("Tasks saved successfully ✅")
}

func loadFromFile() []Task {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("No saved tasks found.")
		return []Task{}
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		fmt.Println("Error decoding tasks:", err)
		return []Task{}
	}
	fmt.Println("Tasks loaded successfully ✅")
	return tasks
}

// ------------------ Core Features ------------------
func addTask() {
	desc := getUserInput("Enter task description: ")
	if strings.TrimSpace(desc) == "" {
		fmt.Println("⚠️ Task description cannot be empty!")
		return
	}
	taskList = append(taskList, Task{Description: desc, Completed: false})
	fmt.Println("Task added successfully ✅")
}

func displayTasks() {
	if len(taskList) == 0 {
		fmt.Println("No tasks available.")
		return
	}
	fmt.Println("\n--- To-Do List ---")
	for i, task := range taskList {
		status := "❌"
		if task.Completed {
			status = "✅"
		}
		fmt.Printf("%d. %20s [%s]\n", i+1, task.Description, status)
	}
}

func markTaskCompleted() {
	index := getIndexFromInput("Enter task number to mark as completed: ")
	if index >= 0 && index < len(taskList) {
		taskList[index].Completed = true
		fmt.Println("Task marked as completed ✅")
	} else {
		fmt.Println("⚠️ Invalid task number!")
	}
}

func removeTask() {
	index := getIndexFromInput("Enter task number to remove: ")
	if index >= 0 && index < len(taskList) {
		taskList = append(taskList[:index], taskList[index+1:]...)
		fmt.Println("Task removed successfully ✅")
	} else {
		fmt.Println("⚠️ Invalid task number!")
	}
}

// ------------------ Helpers ------------------
func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getIndexFromInput(prompt string) int {
	var index int
	input := getUserInput(prompt)
	_, err := fmt.Sscan(input, &index)
	if err != nil || index <= 0 {
		return -1
	}
	return index - 1 // zero-based index
}

func displayMenu() {
	fmt.Println("\n--- To-Do List Menu ---")
	fmt.Println("1. Add Task")
	fmt.Println("2. Display Tasks")
	fmt.Println("3. Mark Task as Completed")
	fmt.Println("4. Remove Task")
	fmt.Println("5. Save Tasks to File")
	fmt.Println("6. Load Tasks from File")
	fmt.Println("7. Exit")
}

// ------------------ Main ------------------

func main() {
	taskList = loadFromFile() // auto-load tasks at start

	for {
		displayMenu()
		choice := getUserInput("Enter your choice: ")

		switch choice {
		case "1":
			addTask()
		case "2":
			displayTasks()
		case "3":
			markTaskCompleted()
		case "4":
			removeTask()
		case "5":
			saveToFile(taskList)
		case "6":
			taskList = loadFromFile()
		case "7":
			fmt.Println("Exiting To-Do List manager. Goodbye 👋")
			return
		default:
			fmt.Println("⚠️ Invalid choice, please try again.")
		}
	}
}
