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

func (c PwdCommand) Execute(ctx context.Context, args []string) Result {
	currentWorkingDirectory, _ := os.Getwd()
	fmt.Println(currentWorkingDirectory)
	return Ok
}
