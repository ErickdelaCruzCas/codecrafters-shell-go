package command

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = home + path[1:]
	}

	if err := c.changeDir(path); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", path)
		return Ok
	}

	return Ok

}
