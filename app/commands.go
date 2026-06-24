package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Executor interface {
	Execute()
}

type ExitCommand struct {
}

func (c *ExitCommand) Execute() {
	os.Exit(0)
}

type EchoCommand struct {
	message string
}

func (e *EchoCommand) Execute() {
	fmt.Printf("%s\n", e.message)
}

type ExternalCommand struct {
	command string
	argv    []string
}

func (c *ExternalCommand) Execute() {
	cmd := exec.Command(c.command, c.argv...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("%s: command not found\n", c.command)
		return
	}
}

type PrintCurrentDirectoryCommand struct {
}

func (c *PrintCurrentDirectoryCommand) Execute() {
	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Error getting current directory")
		return
	}

	fmt.Println(cwd)
}

type TypeCommand struct {
	builtInCommands []string
	command         string
}

func (t *TypeCommand) Execute() {
	if slices.Contains(t.builtInCommands, t.command) {
		fmt.Printf("%s is a shell builtin\n", t.command)
	} else {
		path, err := exec.LookPath(t.command)

		if err != nil {
			fmt.Printf("%s: not found\n", t.command)
			return
		}

		fmt.Printf("%s is %s\n", t.command, path)
	}
}

type ChangeDirectoryCommand struct {
	path string
}

func (c *ChangeDirectoryCommand) Execute() {
	c.path = strings.Replace(c.path, "~", os.Getenv("HOME"), 1)

	err := os.Chdir(c.path)

	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", c.path)
		return
	}
}
