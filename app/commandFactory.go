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

	redirectIndex := len(args)
	for i, arg := range args {
		if arg == ">" || arg == "1>" {
			redirectIndex = i
			break
		}
	}

	var executor Executor
	switch args[0] {
	case echoCommandName:
		executor = &EchoCommand{
			message: strings.Join(args[1:redirectIndex], " "),
		}
	case exitCommandName:
		executor = &ExitCommand{}
	case typeCommandName:
		executor = &TypeCommand{
			builtInCommands: strings.Split(builtIns, ", "),
			command:         args[1],
		}
	case printCurrentDirectoryCommandName:
		executor = &PrintCurrentDirectoryCommand{}
	case changeDirectoryCommandName:
		executor = &ChangeDirectoryCommand{
			path: args[1],
		}
	default:
		executor = &ExternalCommand{
			command: args[0],
			argv:    args[1:redirectIndex],
		}
	}

	if redirectIndex == len(args) {
		return &NoRedirect{
			executor: executor,
		}
	}

	return &Redirect{
		executor: executor,
		filePath: args[redirectIndex+1],
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

		if ch == '\\' {
			escaping = true
			continue
		}

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = true
			continue
		}

		if ch == ' ' && !inDoubleQuote {
			if currToken.Len() > 0 {
				tokens = append(tokens, currToken.String())
				currToken.Reset()
			}
			continue
		}

		currToken.WriteByte(byte(ch))
	}

	if currToken.Len() > 0 {
		tokens = append(tokens, currToken.String())
	}

	return tokens
}
