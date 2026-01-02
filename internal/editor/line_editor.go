package editor

import (
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
)

type LineEditor struct {
	buffer      []rune
	builtins    []string
	executables []string

	lastWasTab bool
	history    []string
	histIndex  int
	savedInput []rune
}

func New(candidates []string, excutables []string) *LineEditor {
	return &LineEditor{
		buffer:      make([]rune, 0),
		builtins:    candidates,
		executables: excutables,
		histIndex:   -1,
	}
}

func (e *LineEditor) SetHistory(entries []string) {
	e.history = entries
	if e.histIndex >= len(entries) {
		e.histIndex = -1
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
	e.histIndex = -1
	e.savedInput = nil

	for {
		var b [1]byte
		if _, err := os.Stdin.Read(b[:]); err != nil {
			return "", err
		}

		switch b[0] {
		case 27: // ESC
			seq, err := e.readEscapeSeq()
			if err != nil {
				return "", err
			}
			if len(seq) == 2 && seq[0] == '[' {
				switch seq[1] {
				case 'A':
					e.historyUp()
				case 'B':
					e.historyDown()
				}
			}

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

func (e *LineEditor) readEscapeSeq() ([]byte, error) {
	buf := make([]byte, 2)
	for i := 0; i < 2; i++ {
		var b [1]byte
		if _, err := os.Stdin.Read(b[:]); err != nil {
			return nil, err
		}
		buf[i] = b[0]
	}
	return buf, nil
}

func (e *LineEditor) historyUp() {
	if len(e.history) == 0 {
		return
	}
	if e.histIndex == -1 {
		e.savedInput = append([]rune(nil), e.buffer...)
		e.histIndex = len(e.history) - 1
	} else if e.histIndex > 0 {
		e.histIndex--
	}
	e.buffer = []rune(e.history[e.histIndex])
	e.redraw()
}

func (e *LineEditor) historyDown() {
	if e.histIndex == -1 {
		return
	}
	if e.histIndex < len(e.history)-1 {
		e.histIndex++
		e.buffer = []rune(e.history[e.histIndex])
		e.redraw()
		return
	}
	e.histIndex = -1
	e.buffer = append([]rune(nil), e.savedInput...)
	e.redraw()
}

func (e *LineEditor) redraw() {
	os.Stdout.Write([]byte("\r\033[K"))
	os.Stdout.Write([]byte("$ "))
	os.Stdout.Write([]byte(string(e.buffer)))
	e.lastWasTab = false
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

	sort.Strings(matches)

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
