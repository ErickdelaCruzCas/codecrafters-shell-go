package shell

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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

func (s *Shell) IsExecutable(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}

func (s *Shell) ChangeDir(path string) error {
	return os.Chdir(path)
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
		tokens := tokenizer(line)

		name := tokens[0]
		args := tokens[1:]

		cmd, ok := s.commands[name]
		if !ok {
			path, ok := s.IsExecutable(name)
			if ok {
				externalCmd := exec.CommandContext(ctx, path, args...)
				externalCmd.Args[0] = name
				externalCmd.Stdin = os.Stdin
				externalCmd.Stdout = os.Stdout
				externalCmd.Stderr = os.Stderr
				if err := externalCmd.Run(); err != nil {
					fmt.Println(err)
				}
				continue
			}

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

func tokenizer(line string) []string {
	tokens := []string{}
	var token string
	singleQuotes := false
	for _, ch := range line {

		// 1️⃣ Comilla: abre o cierra modo quoted
		if ch == '\'' {
			singleQuotes = !singleQuotes
			continue
		}

		// 2️⃣ Espacio fuera de comillas: cierra token
		if ch == ' ' && !singleQuotes {
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
			continue
		}

		// 3️⃣ Cualquier otro carácter
		token += string(ch)

	}

	// 4️⃣ Añadir último token si existe
	if token != "" {
		tokens = append(tokens, token)
	}

	return tokens
}
