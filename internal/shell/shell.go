package shell

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
)

type Shell struct {
	commands map[string]command.Command
}

func New(commands map[string]command.Command) *Shell {
	return &Shell{
		commands: commands,
	}
}

func (s *Shell) IsBuiltin(name string) bool {
	_, ok := s.commands[name]
	return ok
}

func (s *Shell) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("$ ")

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		parts := strings.Split(line, " ")

		name := parts[0]
		args := parts[1:]

		cmd, ok := s.commands[name]
		if !ok {
			fmt.Println(name + ": command not found")
			continue
		}

		result := cmd.Execute(ctx, args)

		switch result {
		case command.Ok:
			// continuar
		case command.Exit:
			return
		case command.Error:
			fmt.Println("command error")
		}
	}
}
