package parser

import "fmt"

type Redirect struct {
	Stdout       string
	StdoutAppend bool
	Stderr       string
	StderrAppend bool
}

/* =========================
     REDIRECT PARSER
========================= */

func ParseRedirect(tokens []string) (cmd string, args []string, redir Redirect, err error) {
	if len(tokens) == 0 {
		return "", nil, redir, nil
	}

	cmd = tokens[0]

	for i := 1; i < len(tokens); i++ {
		switch tokens[i] {

		case ">", "1>":
			if i+1 >= len(tokens) {
				return "", nil, redir, fmt.Errorf("syntax error near >")
			}
			redir.Stdout = tokens[i+1]
			redir.StdoutAppend = false
			i++

		case ">>", "1>>":
			if i+1 >= len(tokens) {
				return "", nil, redir, fmt.Errorf("syntax error near >>")
			}
			redir.Stdout = tokens[i+1]
			redir.StdoutAppend = true
			i++

		case "2>":
			if i+1 >= len(tokens) {
				return "", nil, redir, fmt.Errorf("syntax error near 2>")
			}
			redir.Stderr = tokens[i+1]
			redir.StderrAppend = false
			i++

		case "2>>":
			if i+1 >= len(tokens) {
				return "", nil, redir, fmt.Errorf("syntax error near 2>>")
			}
			redir.Stderr = tokens[i+1]
			redir.StderrAppend = true
			i++

		default:
			args = append(args, tokens[i])
		}
	}

	return cmd, args, redir, nil
}
