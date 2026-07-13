package commands

import "os"

type ExitCommand struct{}

func (c *ExitCommand) Execute() {
	os.Exit(0)
}
