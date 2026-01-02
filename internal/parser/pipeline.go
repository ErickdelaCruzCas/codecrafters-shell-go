package parser

import "fmt"

type CommandLine struct {
	Name  string
	Args  []string
	Redir Redirect
}

func ParsePipeline(tokens []string) ([]CommandLine, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	segments := make([][]string, 0, 1)
	current := make([]string, 0)

	for _, tok := range tokens {
		if tok == "|" {
			if len(current) == 0 {
				return nil, fmt.Errorf("syntax error near |")
			}
			segments = append(segments, current)
			current = make([]string, 0)
			continue
		}
		current = append(current, tok)
	}

	if len(current) == 0 {
		return nil, fmt.Errorf("syntax error near |")
	}
	segments = append(segments, current)

	commands := make([]CommandLine, 0, len(segments))
	for _, segment := range segments {
		name, args, redir, err := ParseRedirect(segment)
		if err != nil {
			return nil, err
		}
		if name == "" {
			return nil, fmt.Errorf("syntax error near |")
		}
		commands = append(commands, CommandLine{
			Name:  name,
			Args:  args,
			Redir: redir,
		})
	}

	return commands, nil
}
