package editor

import (
	"os"

	"golang.org/x/term"
)

type LineEditor struct {
	buffer      []rune
	candidates  []string
	executables []string
}

func New(candidates []string) *LineEditor {
	return &LineEditor{
		buffer:     make([]rune, 0),
		candidates: candidates,
	}
}

func (e *LineEditor) ReadLine() (string, error) {
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	e.buffer = e.buffer[:0]

	for {
		var b [1]byte
		if _, err := os.Stdin.Read(b[:]); err != nil {
			return "", err
		}

		switch b[0] {

		case '\n', '\r':
			os.Stdout.Write([]byte("\r\n"))
			return string(e.buffer), nil

		case '\t':
			e.autocomplete()

		case 127: // backspace
			if len(e.buffer) > 0 {
				e.buffer = e.buffer[:len(e.buffer)-1]
				os.Stdout.Write([]byte("\b \b"))
			}

		default:
			e.buffer = append(e.buffer, rune(b[0]))
			os.Stdout.Write(b[:])
		}
	}
}

func (e *LineEditor) autocomplete() {
	prefix := string(e.buffer)

	matches := make([]string, 0)
	for _, b := range e.candidates {
		if len(b) >= len(prefix) && b[:len(prefix)] == prefix {
			matches = append(matches, b)
		}
	}

	if len(matches) == 1 {
		// completar
		rest := matches[0][len(prefix):]
		for _, r := range rest {
			e.buffer = append(e.buffer, r)
			os.Stdout.Write([]byte(string(r)))
		}
		e.buffer = append(e.buffer, ' ')
		os.Stdout.Write([]byte(" "))
	} else if len(matches) == 0 {
		os.Stdout.Write([]byte{0x07}) // bell
	}
}
