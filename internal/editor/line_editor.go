package editor

import (
	"os"
	"strings"

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
	buf := string(e.buffer)

	// 1. separar head y token activo
	lastSpace := strings.LastIndex(buf, " ")

	var head, token string
	if lastSpace == -1 {
		head = ""
		token = buf
	} else {
		head = buf[:lastSpace+1] // incluye el espacio
		token = buf[lastSpace+1:]
	}

	// 2. buscar matches sobre el token
	matches := make([]string, 0)
	for _, c := range e.candidates {
		if strings.HasPrefix(c, token) {
			matches = append(matches, c)
		}
	}

	// 3. sin matches → bell
	if len(matches) == 0 {
		os.Stdout.Write([]byte{0x07})
		return
	}

	// 4. único match → reemplazar token + espacio
	if len(matches) == 1 {
		match := matches[0]

		// borrar token actual del buffer
		suffix := match[len(token):]

		// escribir match completo
		for _, r := range suffix {
			e.buffer = append(e.buffer, r)
			os.Stdout.Write([]byte(string(r)))
		}

		// añadir espacio final
		e.buffer = append(e.buffer, ' ')
		os.Stdout.Write([]byte(" "))
		return
	}

	// 5. múltiples matches → por ahora, no hacemos nada
	// (más adelante se listan)
}
