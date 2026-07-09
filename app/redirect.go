package main

import "os"

type Redirect struct {
	executor Executor
	filePath string
}

func (c *Redirect) Execute() {
	file, err := os.Create(c.filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	temp := os.Stdout
	defer func() {
		os.Stdout = temp
	}()

	os.Stdout = file
	c.executor.Execute()
}

type NoRedirect struct {
	executor Executor
}

func (c *NoRedirect) Execute() {
	c.executor.Execute()
}
