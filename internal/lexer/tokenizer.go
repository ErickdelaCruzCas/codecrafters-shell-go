package lexer

/* =========================
       TOKENIZER
========================= */

func Tokenizer(line string) []string {
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
				switch ch {
				case '"', '\\', ' ':
					token += string(ch)
				default:
					token += "\\" + string(ch)
				}
			} else {
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
