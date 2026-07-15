package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

const (
	ECHO_COMMAND_NAME                    = "echo"
	EXIT_COMMAND_NAME                    = "exit"
	TYPE_COMMAND_NAME                    = "type"
	PRINT_CURRENT_DIRECTORY_COMMAND_NAME = "pwd"
	CHANGE_DIRECTORY_COMMAND_NAME        = "cd"
)

const BUILT_INS = "echo, exit, type, pwd, cd"

func main() {
	completer := readline.NewPrefixCompleter(
		readline.PcItem(ECHO_COMMAND_NAME),
		readline.PcItem(EXIT_COMMAND_NAME),
	)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing readline:", err)
		return
	}
	defer rl.Close()

	commandFactory := &CommandFactory{parser: &Parser{}}

	for {
		command, err := rl.Readline()

		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("Error reading command:", err)
			return
		}

		command = strings.TrimSpace(command)
		executor := commandFactory.NewCommand(command)
		executor.Execute()
	}
}
