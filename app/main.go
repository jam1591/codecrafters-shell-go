package main

import (
	"fmt"
	"os"
	"path/filepath"
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
}

func (b *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	matches, length := b.AutoCompleter.Do(line, pos)

	if len(matches) == 0 {
		fmt.Fprint(os.Stderr, "\a")
	}

	return matches, length
}

func main() {
	var completers []readline.PrefixCompleterInterface
	completers = append(completers, readline.PcItem(ECHO_COMMAND_NAME))
	completers = append(completers, readline.PcItem(EXIT_COMMAND_NAME))

	for _, path := range filepath.SplitList(os.Getenv("PATH")) {
		files, _ := os.ReadDir(path)
		for _, f := range files {
			info, _ := f.Info()
			if !info.IsDir() && info.Mode().Perm()&0111 != 0 {
				completers = append(completers, readline.PcItem(info.Name()))
			}
		}
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: &Completer{readline.NewPrefixCompleter(completers...)},
	})

	if err != nil {
		panic(err)
	}
	defer rl.Close()

	commandFactory := &CommandFactory{parser: &Parser{}}

	for {
		command, err := rl.Readline()

		if err != nil {
			break
		}

		executor := commandFactory.NewCommand(strings.TrimSpace(command))
		executor.Execute()
	}
}
