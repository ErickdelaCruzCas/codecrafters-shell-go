package command

import (
	"context"
	"io"
)

type Result int

const (
	Ok Result = iota
	Exit
	Error
)

type Command interface {
	Name() string
	Execute(ctx context.Context, args []string, io IO) Result
}

type IO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
