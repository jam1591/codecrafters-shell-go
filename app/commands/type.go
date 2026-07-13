package commands

import (
	"fmt"
	"os/exec"
	"slices"
)

type TypeCommand struct {
	BuiltInCommands []string
	Command         string
}

func (t *TypeCommand) Execute() {
	if slices.Contains(t.BuiltInCommands, t.Command) {
		fmt.Printf("%s is a shell builtin\n", t.Command)
		return
	}

	path, err := exec.LookPath(t.Command)
	if err != nil {
		fmt.Printf("%s: not found\n", t.Command)
		return
	}

	fmt.Printf("%s is %s\n", t.Command, path)
}
