// package main

// import (
// 	"github.com/codecrafters-io/shell-starter-go/internal/command"
// 	"github.com/codecrafters-io/shell-starter-go/internal/shell"
// )

// func main() {
// 	commands := map[string]command.Command{
// 		"exit": command.ExitCommand{},
// 		"echo": command.EchoCommand{},
// 		"pwd":  command.PwdCommand{},
// 	}

// 	sh := shell.New(commands)

// 	commands["type"] = command.NewTypeCommand(sh.IsBuiltin, sh.IsExecutable)
// 	commands["cd"] = command.NewCdCommand(sh.ChangeDir)

// 	sh.Run()
// }