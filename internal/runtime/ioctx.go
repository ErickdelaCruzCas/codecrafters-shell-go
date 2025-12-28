package runtime

import (
	"os"

	"github.com/codecrafters-io/shell-starter-go/internal/parser"
)

/* =========================
     IO CONTEXT
========================= */

type IOContext struct {
	oldStdout *os.File
	oldStderr *os.File
	outFile   *os.File
	errFile   *os.File
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
		io.oldStdout = os.Stdout
		io.outFile = f
		os.Stdout = f
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
			io.Restore()
			return err
		}
		io.oldStderr = os.Stderr
		io.errFile = f
		os.Stderr = f
	}

	return nil
}

func (io *IOContext) Restore() {
	if io.errFile != nil {
		os.Stderr = io.oldStderr
		io.errFile.Close()
		io.errFile = nil
	}
	if io.outFile != nil {
		os.Stdout = io.oldStdout
		io.outFile.Close()
		io.outFile = nil
	}
}
