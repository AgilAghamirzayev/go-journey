package main

import (
	"fmt"
	"os/exec"
)

func main() {

	// Command to be executed
	cmd := exec.Command("ls", "-l") // Example command: list files in the current directory

	// Execute the command and check for errors
	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	// Print the output of the command
	fmt.Println("Command output:")
	fmt.Println(string(output))
}
