package command

import (
	"context"
	"fmt"

	"github.com/codecrafters-io/shell-starter-go/internal/history"
)

type HistoryCommand struct {
	store *history.Store
}

func NewHistoryCommand(store *history.Store) HistoryCommand {
	return HistoryCommand{
		store: store,
	}
}

func (h HistoryCommand) Name() string {
	return "history"
}

func (h HistoryCommand) Execute(ctx context.Context, args []string, io IO) Result {
	entries := h.store.List()
	for i, entry := range entries {
		fmt.Fprintf(io.Stdout, "%5d  %s\n", i+1, entry)
	}

	return Ok
}
