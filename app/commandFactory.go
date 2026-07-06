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
	var curr int
	var tokens []string
	var currToken string

	for curr < len(rawCmd) {
		switch rawCmd[curr] {
		case '\\':
			curr += 1
			currToken = currToken + rawCmd[curr:curr+1]
			curr += 1
		case '"':
			curr += 1
			temp := curr
			for ; curr < len(rawCmd) && rawCmd[curr] != '"'; curr += 1 {
				if rawCmd[curr] == '\\' {
					curr += 1
					currToken = currToken + rawCmd[curr:curr+1]
					curr += 1
				}
			}
			currToken += rawCmd[temp:curr]
			curr += 1
		case '\'':
			curr += 1
			temp := curr
			for ; curr < len(rawCmd) && rawCmd[curr] != '\''; curr += 1 {
			}
			currToken += rawCmd[temp:curr]
			curr += 1
		case ' ':
			if currToken != "" {
				tokens = append(tokens, currToken)
			}
			currToken = ""
			curr += 1
		default:
			currToken += string(rawCmd[curr])
			curr += 1
		}
	}

	if currToken != "" {
		tokens = append(tokens, currToken)
	}

	return tokens
}
