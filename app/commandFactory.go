package main

import (
	"shell-starter-go/app/commands"
	"strings"
)

type CommandFactory struct {
	parser *Parser
}

func (f *CommandFactory) NewCommand(cmd string) commands.Executor {
	parsed := f.parser.Parse(cmd)

	var executor commands.Executor
	switch parsed.Command {
	case ECHO_COMMAND_NAME:
		executor = &commands.EchoCommand{
			Message: strings.Join(parsed.Args, " "),
		}
	case EXIT_COMMAND_NAME:
		executor = &commands.ExitCommand{}
	case TYPE_COMMAND_NAME:
		executor = &commands.TypeCommand{
			BuiltInCommands: strings.Split(BUILT_INS, ", "),
			Command:         parsed.Args[0],
		}
	case PRINT_CURRENT_DIRECTORY_COMMAND_NAME:
		executor = &commands.PrintCurrentDirectoryCommand{}
	case CHANGE_DIRECTORY_COMMAND_NAME:
		executor = &commands.ChangeDirectoryCommand{
			Path: parsed.Args[0],
		}
	default:
		executor = &commands.ExternalCommand{
			Command: parsed.Command,
			Argv:    parsed.Args,
		}
	}

	return &commands.RedirectCommand{
		Metadata: parsed.Redirect,
		Inner:    executor,
	}
}
