package main

import "os"

const (
	OVERRIDE = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	APPEND   = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

type RedirectMetadata struct {
	filePathStdout  string
	filePathStderr  string
	isAppend        bool
	commandEndIndex int
}

type RedirectParser struct {
	metadata RedirectMetadata
}

func (p *RedirectParser) ParseInfo(args []string) RedirectMetadata {
	for i, arg := range args {
		switch arg {
		case ">>", "1>>":
			p.metadata.filePathStdout = args[i+1]
			p.metadata.commandEndIndex = i
			p.metadata.isAppend = true
		case "2>>":
			p.metadata.filePathStderr = args[i+1]
			p.metadata.commandEndIndex = i
			p.metadata.isAppend = true
		case ">", "1>":
			p.metadata.filePathStdout = args[i+1]
			p.metadata.commandEndIndex = i
			p.metadata.isAppend = false
		case "2>":
			p.metadata.filePathStderr = args[i+1]
			p.metadata.commandEndIndex = i
			p.metadata.isAppend = false
		default:
			continue
		}
	}
	return p.metadata
}

type RedirectCommand struct {
	metadata RedirectMetadata
	inner    Executor
}

func (c *RedirectCommand) Execute() {
	if c.metadata.filePathStdout != "" {
		file, err := os.OpenFile(c.metadata.filePathStdout, getFlagForRedirect(c.metadata.isAppend), 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		temp := os.Stdout
		defer func() { os.Stdout = temp }()
		os.Stdout = file
	}

	if c.metadata.filePathStderr != "" {
		file, err := os.OpenFile(c.metadata.filePathStderr, getFlagForRedirect(c.metadata.isAppend), 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		temp := os.Stderr
		defer func() { os.Stderr = temp }()
		os.Stderr = file
	}

	c.inner.Execute()
}

type RedirectStdout struct {
	executor Executor
	filePath string
	isAppend bool
}

func getFlagForRedirect(isAppend bool) int {
	if isAppend {
		return APPEND
	}
	return OVERRIDE
}
