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
	inSingleQuote := false
	inDoubleQuote := false

	for _, ch := range line {

		// 1️⃣ Single quote: toggle single quote mode (unless inside double quotes)
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}

		// 2️⃣ Double quote: toggle double quote mode (unless inside single quotes)
		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		// 3️⃣ Space outside quotes: ends current token
		if ch == ' ' && !inSingleQuote && !inDoubleQuote {
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
			continue
		}

		// 4️⃣ Any other character (including quotes inside the opposite quote type)
		token += string(ch)

	}

	// 5️⃣ Add last token if it exists
	if token != "" {
		tokens = append(tokens, token)
	}

	return tokens
}
