package command

import (
	"context"
	"fmt"
	"strings"
)

type EchoCommand struct{}

func (c EchoCommand) Name() string {
	return "echo"
}

func (c EchoCommand) Execute(ctx context.Context, args []string) Result {
	fmt.Println(strings.Join(args, " "))
	return Ok
}
