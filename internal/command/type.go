package command

import (
	"context"
	"fmt"
)

type TypeCommand struct {
	isBuiltin func(string) bool
}

func NewTypeCommand(isBuiltin func(string) bool) TypeCommand {
	return TypeCommand{
		isBuiltin: isBuiltin,
	}
}

func (c TypeCommand) Name() string {
	return "type"
}

func (c TypeCommand) Execute(ctx context.Context, args []string) Result {
	if len(args) == 0 {
		return Ok
	}

	name := args[0]

	if c.isBuiltin(name) {
		fmt.Printf("%s is a shell builtin\n", name)
		return Ok
	}

	fmt.Printf("%s: not found\n", name)
	return Ok
}
