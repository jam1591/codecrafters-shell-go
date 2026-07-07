package main

import (
	"strings"
)

const (
	echoCommandName                  = "echo"
	exitCommandName                  = "exit"
	typeCommandName                  = "type"
	printCurrentDirectoryCommandName = "pwd"
	changeDirectoryCommandName       = "cd"
)

const builtIns = "echo, exit, type, pwd, cd"

type CommandFactory struct {
}

func (f *CommandFactory) NewCommand(cmd string) Executor {
	args := parseTokens(cmd)

	switch args[0] {
	case echoCommandName:
		return &EchoCommand{
			message: strings.Join(args[1:], " "),
		}
	case exitCommandName:
		return &ExitCommand{}
	case typeCommandName:
		return &TypeCommand{
			builtInCommands: strings.Split(builtIns, ", "),
			command:         args[1],
		}
	case printCurrentDirectoryCommandName:
		return &PrintCurrentDirectoryCommand{}
	case changeDirectoryCommandName:
		return &ChangeDirectoryCommand{
			path: args[1],
		}
	default:
		return &ExternalCommand{
			command: args[0],
			argv:    args[1:],
		}
	}
}

func parseTokens(rawCmd string) []string {
	var tokens []string
	var currToken strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	escaping := false

	for _, ch := range rawCmd {
		if inSingleQuote {
			if ch == '\'' {
				inSingleQuote = false
				continue
			}
			currToken.WriteByte(byte(ch))
			continue
		}

		if escaping {
			currToken.WriteByte(byte(ch))
			escaping = false
			continue
		}

		switch ch {
		case '\\':
			escaping = true
		case '"':
			inDoubleQuote = !inDoubleQuote
		case '\'':
			inSingleQuote = true
		case ' ':
			if inDoubleQuote {
				currToken.WriteByte(byte(ch))
				continue
			}
			if currToken.Len() > 0 {
				tokens = append(tokens, currToken.String())
				currToken.Reset()
			}
		default:
			currToken.WriteByte(byte(ch))
		}
	}

	if currToken.Len() > 0 {
		tokens = append(tokens, currToken.String())
	}

	return tokens
}
