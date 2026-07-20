package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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

type State struct {
	tabTime              time.Time
	lastWasAmbiguousBell bool
}

type Completer struct {
	completer readline.AutoCompleter
	state     State
}

func (b *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	matches, length := b.completer.Do(line, pos)

	if len(matches) == 0 {
		// no matches at all — just beep, reset flag
		fmt.Fprint(os.Stderr, "\a")
		b.state.lastWasAmbiguousBell = false
		return matches, length
	}

	if len(matches) == 1 {
		// unambiguous — complete silently, reset flag
		b.state.lastWasAmbiguousBell = false
		return matches, length
	}

	// len(matches) > 1 from here on
	if !b.state.lastWasAmbiguousBell {
		// first TAB: just beep
		fmt.Fprint(os.Stderr, "\a")
		b.state.lastWasAmbiguousBell = true
		return matches, length
	}

	// second TAB: list them
	sorted := make([]string, len(matches))
	for i, m := range matches {
		sorted[i] = string(m)
	}
	sort.Strings(sorted)

	fmt.Println()
	fmt.Println(strings.Join(sorted, "  "))

	b.state.lastWasAmbiguousBell = false

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
		AutoComplete: &Completer{completer: readline.NewPrefixCompleter(completers...)},
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
