package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/editor"
	"github.com/codecrafters-io/shell-starter-go/internal/lexer"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
	shellruntime "github.com/codecrafters-io/shell-starter-go/internal/runtime"
)

type Shell struct {
	commands map[string]command.Command
}

func New(commands map[string]command.Command) *Shell {
	return &Shell{commands: commands}
}

/* =========================
         RUN
========================= */

func (s *Shell) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	editor := editor.New(s.builtinNames())

	for {
		fmt.Print("$ ")
		os.Stdout.Sync()

		line, err := editor.ReadLine()
		if err != nil {
			fmt.Println(err)
			continue
		}

		tokens := lexer.Tokenize(line)
		name, args, redir, err := parser.ParseRedirect(tokens)
		if err != nil {
			fmt.Println(err)
			continue
		}
		io := &shellruntime.IOContext{}
		if err := io.Apply(redir); err != nil {
			fmt.Println(err)
			continue
		}
		shouldExit := s.execute(ctx, name, args)
		io.Restore()
		if shouldExit {
			return
		}
	}
}

/* =========================
      EXECUTION
========================= */

func (s *Shell) execute(ctx context.Context, name string, args []string) bool {
	if cmd, ok := s.commands[name]; ok {
		result := cmd.Execute(ctx, args)
		return result == command.Exit
	}

	if path, ok := s.IsExecutable(name); ok {
		s.executeExternal(ctx, path, args, name)
		return false
	}

	fmt.Println(name + ": command not found")
	return false
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

func (*Shell) executeExternal(ctx context.Context, path string, args []string, name string) {
	externalCmd := exec.CommandContext(ctx, path, args...)
	externalCmd.Args[0] = name
	externalCmd.Stdin = os.Stdin
	externalCmd.Stdout = os.Stdout
	externalCmd.Stderr = os.Stderr
	_ = externalCmd.Run()
}

func (s *Shell) builtinNames() []string {
	names := make([]string, 0, len(s.commands))
	for name := range s.commands {
		names = append(names, name)
	}
	return names
}
