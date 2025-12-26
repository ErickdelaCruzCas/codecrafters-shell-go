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

type Redirect struct {
	Stdout string
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

		name, args, redir, err := parseRedirect(tokens)

		if err != nil {
			fmt.Println(err)
			continue
		}

		var (
			oldStdout *os.File
			outFile   *os.File
		)

		if redir.Stdout != "" {
			outFile, err = os.Create(redir.Stdout)
			if err != nil {
				fmt.Println(err)
				continue
			}
			oldStdout = os.Stdout
			os.Stdout = outFile
		}

		cmd, ok := s.commands[name]
		if !ok {

			path, ok := s.IsExecutable(name)
			if ok {
				// Aqui se ejecuta el external
				s.executeExternal(ctx, path, args, name)
				// restaurar stdout
				if outFile != nil {
					os.Stdout = oldStdout
					outFile.Close()
				}
				continue
			}

			fmt.Println(name + ": command not found")
			continue
		}

		// aquí el built in
		result := cmd.Execute(ctx, args)

		if outFile != nil {
			os.Stdout = oldStdout
			outFile.Close()
		}

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

func (*Shell) executeExternal(ctx context.Context, path string, args []string, name string) {
	externalCmd := exec.CommandContext(ctx, path, args...)
	externalCmd.Args[0] = name
	externalCmd.Stdin = os.Stdin
	externalCmd.Stdout = os.Stdout
	externalCmd.Stderr = os.Stderr
	if err := externalCmd.Run(); err != nil {
		// fmt.Println(err)
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

func parseRedirect(tokens []string) (cmd string, args []string, redir Redirect, err error) {
	if len(tokens) == 0 {
		return "", nil, redir, nil
	}

	cmd = tokens[0]

	for i := 1; i < len(tokens); i++ {
		if tokens[i] == ">" || tokens[i] == "1>" {
			if i+1 >= len(tokens) {
				return "", nil, redir, fmt.Errorf("syntax error near >")
			}
			redir.Stdout = tokens[i+1]
			return cmd, args, redir, nil
		}
		args = append(args, tokens[i])
	}

	return cmd, args, redir, nil
}
