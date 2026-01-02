package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/editor"
	"github.com/codecrafters-io/shell-starter-go/internal/history"
	"github.com/codecrafters-io/shell-starter-go/internal/lexer"
	"github.com/codecrafters-io/shell-starter-go/internal/parser"
	shellruntime "github.com/codecrafters-io/shell-starter-go/internal/runtime"
)

type Shell struct {
	commands map[string]command.Command
	history  *history.Store
}

type runner struct {
	start func() error
	wait  func()
}

type pipeSetup struct {
	ioCtx      *shellruntime.IOContext
	pipeWriter *io.PipeWriter
	closeStdin io.Closer
	closePipe  bool
}

func New(commands map[string]command.Command, historyStore *history.Store) *Shell {
	return &Shell{
		commands: commands,
		history:  historyStore,
	}
}

/* =========================
         RUN
========================= */

func (s *Shell) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	editor := editor.New(s.builtinNames(), s.executablesInPath())

	for {
		fmt.Print("$ ")
		os.Stdout.Sync()

		if s.history != nil {
			editor.SetHistory(s.history.List())
		}

		line, err := editor.ReadLine()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if s.history != nil && strings.TrimSpace(line) != "" {
			s.history.Add(line)
		}

		tokens := lexer.Tokenize(line)
		commands, err := parser.ParsePipeline(tokens)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(commands) == 0 {
			continue
		}
		shouldExit := s.executePipeline(ctx, commands)
		if shouldExit {
			return
		}
	}
}

/* =========================
      EXECUTION
========================= */

func (s *Shell) executePipeline(ctx context.Context, pipeline []parser.CommandLine) bool {
	var exitRequested int32

	runners := make([]runner, 0, len(pipeline))

	var prevReader io.Reader = os.Stdin

	for i, cmdLine := range pipeline {
		var pipeReader *io.PipeReader
		var pipeWriter *io.PipeWriter
		if i < len(pipeline)-1 {
			pipeReader, pipeWriter = io.Pipe()
		}

		setup, err := s.preparePipelineIO(prevReader, pipeWriter, cmdLine.Redir)
		if err != nil {
			fmt.Println(err)
			return false
		}

		if builtin, ok := s.commands[cmdLine.Name]; ok {
			runners = append(runners, s.newBuiltinRunner(ctx, builtin, cmdLine.Args, setup, &exitRequested))
		} else if path, ok := s.IsExecutable(cmdLine.Name); ok {
			runners = append(runners, s.newExternalRunner(ctx, path, cmdLine.Args, cmdLine.Name, setup))
		} else {
			fmt.Println(cmdLine.Name + ": command not found")
			if setup.pipeWriter != nil {
				setup.pipeWriter.Close()
			}
			setup.ioCtx.Close()
			return false
		}

		if pipeReader != nil {
			prevReader = pipeReader
		}
	}

	for i, r := range runners {
		if err := r.start(); err != nil {
			fmt.Println(err)
			for j := 0; j < i; j++ {
				runners[j].wait()
			}
			return false
		}
	}

	for _, r := range runners {
		r.wait()
	}

	return atomic.LoadInt32(&exitRequested) == 1
}

func (s *Shell) preparePipelineIO(prevReader io.Reader, pipeWriter *io.PipeWriter, redir parser.Redirect) (pipeSetup, error) {
	ioCtx := shellruntime.NewIOContext()
	ioCtx.Stdin = prevReader
	if pipeWriter != nil {
		ioCtx.Stdout = pipeWriter
	}
	if err := ioCtx.Apply(redir); err != nil {
		if pipeWriter != nil {
			pipeWriter.Close()
		}
		ioCtx.Close()
		return pipeSetup{}, err
	}

	if pipeWriter != nil && ioCtx.Stdout != pipeWriter {
		pipeWriter.Close()
		pipeWriter = nil
	}

	var closeStdin io.Closer
	if r, ok := prevReader.(*io.PipeReader); ok {
		closeStdin = r
	}

	closePipe := pipeWriter != nil && ioCtx.Stdout == pipeWriter

	return pipeSetup{
		ioCtx:      ioCtx,
		pipeWriter: pipeWriter,
		closeStdin: closeStdin,
		closePipe:  closePipe,
	}, nil
}

func (s *Shell) closePipelineIO(setup pipeSetup) {
	if setup.closePipe {
		setup.pipeWriter.Close()
	}
	if setup.closeStdin != nil {
		setup.closeStdin.Close()
	}
	setup.ioCtx.Close()
}

func (s *Shell) newBuiltinRunner(
	ctx context.Context,
	builtin command.Command,
	args []string,
	setup pipeSetup,
	exitFlag *int32,
) runner {
	done := make(chan command.Result, 1)

	return runner{
		start: func() error {
			go func() {
				result := builtin.Execute(ctx, args, command.IO{
					Stdin:  setup.ioCtx.Stdin,
					Stdout: setup.ioCtx.Stdout,
					Stderr: setup.ioCtx.Stderr,
				})
				if result == command.Exit {
					atomic.StoreInt32(exitFlag, 1)
				}
				s.closePipelineIO(setup)
				done <- result
			}()
			return nil
		},
		wait: func() {
			<-done
		},
	}
}

func (s *Shell) newExternalRunner(
	ctx context.Context,
	path string,
	args []string,
	name string,
	setup pipeSetup,
) runner {
	externalCmd := exec.CommandContext(ctx, path, args...)
	externalCmd.Args[0] = name
	externalCmd.Stdin = setup.ioCtx.Stdin
	externalCmd.Stdout = setup.ioCtx.Stdout
	externalCmd.Stderr = setup.ioCtx.Stderr

	return runner{
		start: func() error {
			return externalCmd.Start()
		},
		wait: func() {
			_ = externalCmd.Wait()
			s.closePipelineIO(setup)
		},
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

func (s *Shell) builtinNames() []string {
	names := make([]string, 0, len(s.commands))
	for name := range s.commands {
		names = append(names, name)
	}
	return names
}

func (s *Shell) executablesInPath() []string {
	seen := make(map[string]struct{})
	result := []string{}

	pathEnv := os.Getenv("PATH")
	dirs := filepath.SplitList(pathEnv)

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			// fullPath := filepath.Join(dir, name)

			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Unix: ejecutable si tiene bit x
			if info.Mode()&0111 != 0 {
				if _, ok := seen[name]; !ok {
					seen[name] = struct{}{}
					result = append(result, name)
				}
			}
		}
	}

	return result
}
