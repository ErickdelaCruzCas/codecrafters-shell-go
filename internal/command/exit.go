package command

import (
	"context"
	"os"
)

type ExitCommand struct{}

func (C ExitCommand) Name() string {
	return "exit"
}

func (c ExitCommand) Execute(ctx context.Context, args []string) Result {
	os.Exit(0)
	return Exit
}
