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

type Completer struct {
	readline.AutoCompleter
	term *readline.Terminal
}

func (b *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	matches, length := b.AutoCompleter.Do(line, pos)

	if len(matches) == 0 {
		fmt.Fprint(os.Stderr, "\a")
	}

	return matches, length
}

func main() {
	prefixCompleter := readline.NewPrefixCompleter(
		readline.PcItem(ECHO_COMMAND_NAME),
		readline.PcItem(EXIT_COMMAND_NAME),
	)

	completer := &Completer{AutoCompleter: prefixCompleter}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: completer,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing readline:", err)
		return
	}
	defer rl.Close()

	completer.term = rl.Terminal

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
