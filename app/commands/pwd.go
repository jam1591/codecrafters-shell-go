package commands

import (
	"fmt"
	"os"
)

type PrintCurrentDirectoryCommand struct{}

func (c *PrintCurrentDirectoryCommand) Execute() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
		return
	}

	fmt.Println(cwd)
}
