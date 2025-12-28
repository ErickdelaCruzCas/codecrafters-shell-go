package editor

import (
	"os"
	"strings"

	"golang.org/x/term"
)

type LineEditor struct {
	buffer      []rune
	builtins    []string
	executables []string

	lastWasTab bool
}

func New(candidates []string, excutables []string) *LineEditor {
	return &LineEditor{
		buffer:      make([]rune, 0),
		builtins:    candidates,
		executables: excutables,
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

	var token string
	if lastSpace == -1 {
		token = buf
	} else {
		token = buf[lastSpace+1:]
	}

	// 2. buscar primero en builtins
	matches := make([]string, 0)
	for _, c := range e.builtins {
		if strings.HasPrefix(c, token) {
			matches = append(matches, c)
		}
	}

	// si no hay matches en builtins, buscar en ejecutables
	if len(matches) == 0 {
		for _, c := range e.executables {
			if strings.HasPrefix(c, token) {
				matches = append(matches, c)
			}
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
	if len(matches) > 1 {

		lcp := longestCommonPrefix(matches)

		// ¿el LCP añade algo nuevo?
		if len(lcp) > len(token) {
			suffix := lcp[len(token):]

			for _, r := range suffix {
				e.buffer = append(e.buffer, r)
				os.Stdout.Write([]byte(string(r)))
			}

			e.lastWasTab = false
			return
		}

		if !e.lastWasTab {
			os.Stdout.Write([]byte{0x07})
			e.lastWasTab = true
			return
		}
		// salto de línea limpio (raw mode)
		e.listCandidates(matches)
	}
}

func longestCommonPrefix(candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}

	prefix := candidates[0]

	for _, s := range candidates[1:] {
		for !strings.HasPrefix(s, prefix) {
			if prefix == "" {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}

	return prefix
}

func (e *LineEditor) listCandidates(matches []string) {
	os.Stdout.Write([]byte("\r\n"))

	for _, m := range matches {
		os.Stdout.Write([]byte(m))
		os.Stdout.Write([]byte("  "))
	}

	// nueva línea
	os.Stdout.Write([]byte("\r\n"))

	// redibujar prompt + buffer
	os.Stdout.Write([]byte("$ "))
	os.Stdout.Write([]byte(string(e.buffer)))

	e.lastWasTab = false
}
