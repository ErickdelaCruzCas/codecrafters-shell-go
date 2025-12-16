package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printPrompt() {
	fmt.Fprint(os.Stdout, "$ ")
}

func main() {
	for {
		printPrompt()
		handleInput()
	}
}

func handleInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("$ ")
		if !scanner.Scan() {
			return
		}

		command := scanner.Text()
		splitCommand := strings.Split(command, " ")

		keyword := splitCommand[0]

		switch keyword {
		case "exit":
			os.Exit(1)
		case "echo":
			fmt.Println(splitCommand[1])
		default:
			// By default some CLI tools return only the keyword
			fmt.Println(keyword + ": command not found")
		}
	}

}
