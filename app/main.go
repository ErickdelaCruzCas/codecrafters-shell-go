package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/history"
	"github.com/codecrafters-io/shell-starter-go/internal/shell"
)

func main() {
	historyStore := history.New()
	historyFile := os.Getenv("HISTFILE")
	if historyFile != "" {
		if err := historyStore.LoadFrom(historyFile); err != nil && !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	commands := map[string]command.Command{
		"exit":    command.ExitCommand{},
		"echo":    command.EchoCommand{},
		"pwd":     command.PwdCommand{},
		"history": command.NewHistoryCommand(historyStore),
	}

	sh := shell.New(commands, historyStore)

	commands["type"] = command.NewTypeCommand(sh.IsBuiltin, sh.IsExecutable)
	commands["cd"] = command.NewCdCommand(sh.ChangeDir)

	sh.Run()

	if historyFile != "" {
		if err := historyStore.WriteTo(historyFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
