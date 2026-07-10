package main

import "os"

const (
	OVERRIDE = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	APPEND   = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

type RedirectStdout struct {
	executor Executor
	filePath string
	isAppend bool
}

func (c *RedirectStdout) Execute() {
	file, err := os.OpenFile(c.filePath, getFlagForRedirect(c.isAppend), 0644)

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

type RedirectStderr struct {
	executor Executor
	filePath string
	isAppend bool
}

func (c *RedirectStderr) Execute() {
	file, err := os.OpenFile(c.filePath, getFlagForRedirect(c.isAppend), 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	temp := os.Stderr
	defer func() {
		os.Stderr = temp
	}()

	os.Stderr = file
	c.executor.Execute()
}

type NoRedirect struct {
	executor Executor
}

func (c *NoRedirect) Execute() {
	c.executor.Execute()
}

func getFlagForRedirect(isAppend bool) int {
	if isAppend {
		return APPEND
	}
	return OVERRIDE
}
