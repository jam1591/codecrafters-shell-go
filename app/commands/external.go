package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type ExternalCommand struct {
	Command string
	Argv    []string
}

func (c *ExternalCommand) Execute() {
	cmd := exec.Command(c.Command, c.Argv...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return
	}

	fmt.Fprintf(os.Stderr, "%s: command not found\n", c.Command)
}
