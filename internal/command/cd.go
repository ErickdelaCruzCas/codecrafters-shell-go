package command

import (
	"context"
	"fmt"
)

type CdCommand struct {
	changeDir func(string) error
}

func NewCdCommand(changeDir func(string) error) CdCommand {
	return CdCommand{
		changeDir: changeDir,
	}
}

func (c CdCommand) Name() string {
	return "cd"
}

func (c CdCommand) Execute(ctx context.Context, args []string) Result {

	if len(args) == 0 {
		// de momento no hacemos nada
		return Ok
	}

	path := args[0]

	if err := c.changeDir(path); err != nil {
		fmt.Println(err)
		return Error
	}

	return Ok

}
