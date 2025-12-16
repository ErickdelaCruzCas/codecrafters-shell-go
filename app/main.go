package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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
			fmt.Println(strings.Join(splitCommand[1:], " "))
		default:
			fmt.Println(keyword + ": command not found")
		}
	}
}
