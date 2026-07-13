package main

import (
	"shell-starter-go/app/commands"
	"strings"
)

type ParsedCommand struct {
	Command  string
	Args     []string
	Redirect commands.RedirectMetadata
}

type Parser struct{}

func (p *Parser) Parse(rawCmd string) ParsedCommand {
	tokens := p.tokenize(rawCmd)
	redirect := p.parseRedirect(tokens)

	return ParsedCommand{
		Command:  tokens[0],
		Args:     tokens[1:redirect.CommandEndIndex],
		Redirect: redirect,
	}
}

func (p *Parser) tokenize(rawCmd string) []string {
	var tokens []string
	var current strings.Builder

	inSingleQuote := false
	inDoubleQuote := false
	escaping := false

	for _, ch := range rawCmd {
		switch {
		case inSingleQuote:
			if ch == '\'' {
				inSingleQuote = false
				continue
			}
			current.WriteRune(ch)
		case escaping:
			current.WriteRune(ch)
			escaping = false
		case ch == '\\':
			escaping = true
		case ch == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
		case ch == '\'' && !inDoubleQuote:
			inSingleQuote = true
		case ch == ' ' && !inDoubleQuote:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (p *Parser) parseRedirect(args []string) commands.RedirectMetadata {
	metadata := commands.RedirectMetadata{CommandEndIndex: len(args)}

	for i, arg := range args {
		switch arg {
		case ">>", "1>>":
			metadata.FilePathStdout = args[i+1]
			metadata.CommandEndIndex = i
			metadata.IsAppend = true
			return metadata
		case "2>>":
			metadata.FilePathStderr = args[i+1]
			metadata.CommandEndIndex = i
			metadata.IsAppend = true
			return metadata
		case ">", "1>":
			metadata.FilePathStdout = args[i+1]
			metadata.CommandEndIndex = i
			metadata.IsAppend = false
			return metadata
		case "2>":
			metadata.FilePathStderr = args[i+1]
			metadata.CommandEndIndex = i
			metadata.IsAppend = false
			return metadata
		}
	}

	return metadata
}
