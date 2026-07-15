package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
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
	fd := int(os.Stdin.Fd())
	reader := bufio.NewReader(os.Stdin)
	builtInsArray := strings.Split(BUILT_INS, ", ")

	for {
		fmt.Print("$ ")

		oldState, err := term.MakeRaw(fd)
		if err != nil {
			panic(err)
		}

		var command []byte
	completion:
		for {
			b, err := reader.ReadByte()
			if err != nil {
				term.Restore(fd, oldState)
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
				os.Exit(1)
			}

			switch b {
			case 13: // Enter
				fmt.Print("\r\n")
				break completion

			case 127, 8: // Backspace
				if len(command) > 0 {
					command = command[:len(command)-1]
					fmt.Print("\b \b")
				}

			case 9: // Tab
				partialCmd := strings.TrimSpace(string(command))
				var matches []string
				for _, cmd := range builtInsArray {
					if strings.HasPrefix(cmd, partialCmd) {
						matches = append(matches, cmd)
					}
				}
				if len(matches) == 1 {
					command = []byte(matches[0] + " ")
					// \r  -> return to column 0
					// \033[K -> erase to end of line (clears any stale chars from a longer previous draw)
					fmt.Print("\r\033[K$ " + string(command))
				} else if len(matches) > 1 {
					fmt.Print("\r\n" + strings.Join(matches, "  ") + "\r\n$ " + string(command))
				}

			default:
				command = append(command, b)
				fmt.Printf("%c", b)
			}
		}

		term.Restore(fd, oldState)

		commandFactory := &CommandFactory{parser: &Parser{}}
		executor := commandFactory.NewCommand(strings.TrimSpace(string(command)))
		executor.Execute()
	}
}
