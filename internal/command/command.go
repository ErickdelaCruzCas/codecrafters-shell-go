package command

import "context"

type Result int

const (
	Ok Result = iota
	Exit
	Error
)

type Command interface {
	Name() string
	Execute(ctx context.Context, args []string) Result
}
