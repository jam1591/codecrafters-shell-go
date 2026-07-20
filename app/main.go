package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

type State struct {
	isLastBellAmbiguous bool
}

type Completer struct {
	completer readline.AutoCompleter
	state     State
}

func (b *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	matches, length := b.completer.Do(line, pos)
	fmt.Fprintf(os.Stderr, "\n[DEBUG] Do called, matches=%d, isLastBellAmbiguous=%v\n", len(matches), b.state.isLastBellAmbiguous)

	if len(matches) == 0 {
		fmt.Fprint(os.Stderr, "\a")
		b.state.isLastBellAmbiguous = false
		return matches, length
	}

	if len(matches) == 1 {
		b.state.isLastBellAmbiguous = false
		return matches, length
	}

	if !b.state.isLastBellAmbiguous {
		fmt.Fprint(os.Stderr, "\a")
		b.state.isLastBellAmbiguous = true
		fmt.Fprintln(os.Stderr, "[DEBUG] first tab, beeped")
		return matches, length
	}

	sorted := make([]string, len(matches))
	for i, m := range matches {
		sorted[i] = string(m)
	}
	sort.Strings(sorted)
	fmt.Fprintf(os.Stderr, "[DEBUG] second tab, raw matches=%q sorted=%q\n", matches, sorted)

	fmt.Println()
	fmt.Println(strings.Join(sorted, "  "))

	b.state.isLastBellAmbiguous = false

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
