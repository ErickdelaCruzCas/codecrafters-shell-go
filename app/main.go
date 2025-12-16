package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
// var _ = fmt.Print

func main() {

	fmt.Print("$ ")

	command, err := bufio.NewReader(os.Stdin).ReadString('\n')

	comToExecute := command[:len(command)-1]

	// Wait for user input
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
		os.Exit(1)
	}

	if comToExecute == "exit" {
		os.Exit(0)
	} else {
		fmt.Println(command[:len(command)-1] + ": command not found")
		main()
	}

}
