package commands

import (
	"fmt"
	"os"
	"strings"
)

type ChangeDirectoryCommand struct {
	Path string
}

func (c *ChangeDirectoryCommand) Execute() {
	c.Path = strings.Replace(c.Path, "~", os.Getenv("HOME"), 1)

	err := os.Chdir(c.Path)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", c.Path)
		return
	}
}
