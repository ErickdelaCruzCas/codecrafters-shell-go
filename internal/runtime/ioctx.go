package runtime

import (
	"io"
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

/* =========================
     IO CONTEXT
========================= */

type IOContext struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	outFile *os.File
	errFile *os.File
}

func NewIOContext() *IOContext {
	return &IOContext{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (io *IOContext) Apply(redir parser.Redirect) error {
	// stdout
	if redir.Stdout != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if redir.StdoutAppend {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		f, err := os.OpenFile(redir.Stdout, flags, 0644)
		if err != nil {
			return err
		}
		io.outFile = f
		io.Stdout = f
	}

	// stderr
	if redir.Stderr != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if redir.StderrAppend {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		f, err := os.OpenFile(redir.Stderr, flags, 0644)
		if err != nil {
			io.Close()
			return err
		}
		io.errFile = f
		io.Stderr = f
	}

	return nil
}

func (io *IOContext) Close() {
	if io.errFile != nil {
		io.errFile.Close()
		io.errFile = nil
	}
	if io.outFile != nil {
		io.outFile.Close()
		io.outFile = nil
	}
}
