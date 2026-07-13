package commands

import "os"

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
		redirectStdout(c.Metadata.FilePathStdout, c.Metadata.IsAppend)
	}

	if c.Metadata.FilePathStderr != "" {
		redirectStderr(c.Metadata.FilePathStderr, c.Metadata.IsAppend)
	}

	c.Inner.Execute()
}

func getFlagForRedirect(isAppend bool) int {
	if isAppend {
		return APPEND
	}
	return OVERRIDE
}
