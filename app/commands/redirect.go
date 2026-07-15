package commands

import (
	"fmt"
	"os"
)

const (
	OVERRIDE = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	APPEND   = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

type RedirectMetadata struct {
	FilePathStdout  string
	FilePathStderr  string
	IsAppend        bool
	CommandEndIndex int
}

type RedirectCommand struct {
	Metadata RedirectMetadata
	Inner    Executor
}

func (c *RedirectCommand) Execute() {
	if c.Metadata.FilePathStdout != "" {
		file, err := os.OpenFile(c.Metadata.FilePathStdout, getFlagForRedirect(c.Metadata.IsAppend), 0643)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\r\n", c.Metadata.FilePathStdout, err)
			return
		}
		defer file.Close()
		tempOut := os.Stdout
		defer func() { os.Stdout = tempOut }()
		os.Stdout = file
	}
	if c.Metadata.FilePathStderr != "" {
		file, err := os.OpenFile(c.Metadata.FilePathStderr, getFlagForRedirect(c.Metadata.IsAppend), 0643)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\r\n", c.Metadata.FilePathStderr, err)
			return
		}
		defer file.Close()
		tempErr := os.Stderr
		defer func() { os.Stderr = tempErr }()
		os.Stderr = file
	}

	c.Inner.Execute()
}

func getFlagForRedirect(isAppend bool) int {
	if isAppend {
		return APPEND
	}
	return OVERRIDE
}
