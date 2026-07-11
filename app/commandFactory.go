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
	parser *RedirectParser
}

func NewCommandFactory() *CommandFactory {
	return &CommandFactory{
		parser: &RedirectParser{},
	}
}

func (f *CommandFactory) NewCommand(cmd string) Executor {
	args := parseTokens(cmd)
	if len(args) == 0 {
		return &NoRedirect{executor: nil}
	}

	// Parse redirects
	info := f.parser.Parse(args)
	commandArgs := args[:info.CommandEndIndex]

	// Build executor
	executor := f.buildExecutor(commandArgs)

	// Wrap with redirect
	opts := []RedirectOption{}
	if info.StdoutPath != "" {
		opts = append(opts, WithStdout(info.StdoutPath, info.IsAppend))
	}
	if info.StderrPath != "" {
		opts = append(opts, WithStderr(info.StderrPath, info.IsAppend))
	}

	if len(opts) == 0 {
		return &NoRedirect{executor: executor}
	}

	return Wrap(executor, opts...)
}

func (f *CommandFactory) buildExecutor(args []string) Executor {
	if len(args) == 0 {
		return &NoRedirect{executor: nil}
	}

	switch args[0] {
	case echoCommandName:
		return &EchoCommand{
			message: strings.Join(args[1:], " "),
		}
	case exitCommandName:
		return &ExitCommand{}
	case typeCommandName:
		if len(args) > 1 {
			return &TypeCommand{
				builtInCommands: strings.Split(builtIns, ", "),
				command:         args[1],
			}
		}
		return &TypeCommand{
			builtInCommands: strings.Split(builtIns, ", "),
			command:         "",
		}
	case printCurrentDirectoryCommandName:
		return &PrintCurrentDirectoryCommand{}
	case changeDirectoryCommandName:
		if len(args) > 1 {
			return &ChangeDirectoryCommand{path: args[1]}
		}
		return &ChangeDirectoryCommand{path: ""}
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
