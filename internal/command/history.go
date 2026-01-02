package command

import (
	"context"
	"fmt"
	"strconv"

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
	if len(args) >= 2 {
		switch args[0] {
		case "-r":
			if err := h.store.LoadFrom(args[1]); err != nil {
				fmt.Fprintln(io.Stderr, err)
			}
			return Ok
		case "-w":
			if err := h.store.WriteTo(args[1]); err != nil {
				fmt.Fprintln(io.Stderr, err)
			}
			return Ok
		case "-a":
			if err := h.store.AppendTo(args[1]); err != nil {
				fmt.Fprintln(io.Stderr, err)
			}
			return Ok
		}
	}

	entries := h.store.List()
	start := 0
	if len(args) > 0 {
		limit, err := strconv.Atoi(args[0])
		if err != nil || limit < 0 {
			return Ok
		}
		if limit < len(entries) {
			start = len(entries) - limit
		}
	}

	for i, entry := range entries[start:] {
		fmt.Fprintf(io.Stdout, "%5d  %s\n", start+i+1, entry)
	}

	return Ok
}
