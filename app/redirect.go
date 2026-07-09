package main

import "os"

type RedirectStdout struct {
	executor Executor
	filePath string
}

func (c *RedirectStdout) Execute() {
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

type RedirectStderr struct {
	executor Executor
	filePath string
}

func (c *RedirectStderr) Execute() {
	file, err := os.Create(c.filePath)
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
