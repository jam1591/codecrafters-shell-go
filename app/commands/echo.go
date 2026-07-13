package commands

import "fmt"

type EchoCommand struct {
	Message string
}

func (e *EchoCommand) Execute() {
	fmt.Printf("%s\n", e.Message)
}
