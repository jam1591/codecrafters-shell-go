package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

var _ = fmt.Print

const (
	ECHO_COMMAND_NAME                    = "echo"
	EXIT_COMMAND_NAME                    = "exit"
	TYPE_COMMAND_NAME                    = "type"
	PRINT_CURRENT_DIRECTORY_COMMAND_NAME = "pwd"
	CHANGE_DIRECTORY_COMMAND_NAME        = "cd"
)

const BUILT_INS = "echo, exit, type, pwd, cd"

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")

		// command, _ := reader.ReadString('\n')
		// var command []byte

		var command []byte
	completion:
		for {
			b, err := reader.ReadByte()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
				os.Exit(1)
			}

			builtInsArray := strings.Split(BUILT_INS, ", ")

			switch b {
			case 127, 8:
				fmt.Print("\b \b")
				continue
			case 13:
				fmt.Println()
				break completion
			case 9:
				partialCmd := strings.TrimSpace(string(command))
				var matches = make([]string, 0)
				for _, cmd := range builtInsArray {
					if strings.HasPrefix(cmd, partialCmd) {
						matches = append(matches, cmd)
					}

					if len(matches) == 1 {
						command = []byte(matches[0])
						command = append(command, ' ')
						fmt.Printf("\r$ %s", string(command))
					} else if len(matches) > 1 {
						fmt.Printf("\r\n%v\r\n$ %s", matches, string(command))
					}
				}

				continue
			}

			command = append(command, b)
			fmt.Printf("%c", b)
		}

		commandFactory := &CommandFactory{parser: &Parser{}}
		executor := commandFactory.NewCommand(strings.TrimSpace(string(command)))
		executor.Execute()
	}
}
