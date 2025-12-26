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
	const (
		normal = iota
		singleQuote
		doubleQuote
		escape
	)

	state := normal
	prevState := normal

	var tokens []string
	var token string

	for _, ch := range line {
		switch state {

		case normal:
			switch ch {
			case '\\':
				prevState = normal
				state = escape
			case ' ':
				if token != "" {
					tokens = append(tokens, token)
					token = ""
				}
			case '\'':
				state = singleQuote
			case '"':
				state = doubleQuote
			default:
				token += string(ch)
			}

		case singleQuote:
			if ch == '\'' {
				state = normal
			} else {
				token += string(ch)
			}

		case doubleQuote:
			switch ch {
			case '\\':
				prevState = doubleQuote
				state = escape
			case '"':
				state = normal
			default:
				token += string(ch)
			}

		case escape:
			if prevState == doubleQuote {
				// Dentro de comillas dobles:
				// solo ciertos escapes eliminan el backslash
				switch ch {
				case '"', '\\', ' ':
					token += string(ch)
				default:
					token += "\\" + string(ch)
				}
			} else {
				// Fuera de comillas: el backslash siempre desaparece
				token += string(ch)
			}
			state = prevState
		}
	}

	if token != "" {
		tokens = append(tokens, token)
	}

	return tokens
}
