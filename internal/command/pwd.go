package command

import (
	"context"
	"fmt"
	"os"
)

type PwdCommand struct{}

func (c PwdCommand) Name() string {
	return "pwd"
}

func (c PwdCommand) Execute(ctx context.Context, args []string, io IO) Result {
	currentWorkingDirectory, _ := os.Getwd()
	fmt.Fprintln(io.Stdout, currentWorkingDirectory)
	return Ok
}
