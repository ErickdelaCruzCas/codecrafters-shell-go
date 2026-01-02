package command

import "context"

type ExitCommand struct{}

func (C ExitCommand) Name() string {
	return "exit"
}

func (c ExitCommand) Execute(ctx context.Context, args []string, io IO) Result {
	return Exit
}
