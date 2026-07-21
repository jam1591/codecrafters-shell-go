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

func (c *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	matches, length := c.completer.Do(line, pos)

	unique := make(map[string][]rune)
	for _, m := range matches {
		unique[string(m)] = m
	}

	matches = matches[:0]
	for _, m := range unique {
		matches = append(matches, m)
	}

	if len(matches) == 0 {
		fmt.Fprint(os.Stderr, "\a")
		c.state.isLastBellAmbiguous = false
		return matches, length
	}

	if len(matches) == 1 {
		c.state.isLastBellAmbiguous = false
		return matches, length
	}

	if !c.state.isLastBellAmbiguous {
		fmt.Fprint(os.Stderr, "\a")
		c.state.isLastBellAmbiguous = true
		return nil, 0
	}

	prefix := string(line)
	full := make([]string, len(matches))
	for i, m := range matches {
		full[i] = prefix + string(m)
	}
	sort.Strings(full)

	result := make([][]string, len(full))
	for i, s := range full {
		result[i] = strings.Split(s, "_")
	}

	fmt.Println()
	fmt.Println(strings.Join(full, "  "))
	c.state.isLastBellAmbiguous = false

	for _, s := range result {
		if strings.Join(s[:len(s)-1], "_") == string(line) {
			return [][]rune{[]rune(strings.Join(s, "_"))}, 1
		}
	}

	fmt.Print("$ " + prefix)
	return nil, 0
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
