//go:build unix

package commands

import (
	"os"
	"syscall"
)

func redirectStdout(path string, isAppend bool) {
	redirectFD(syscall.Stdout, path, isAppend, true)
}

func redirectStderr(path string, isAppend bool) {
	redirectFD(syscall.Stderr, path, isAppend, false)
}

func redirectFD(fd int, path string, isAppend bool, isStdout bool) {
	file, err := os.OpenFile(path, getFlagForRedirect(isAppend), 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	saved, err := syscall.Dup(fd)
	if err != nil {
		panic(err)
	}
	defer func() {
		syscall.Dup2(saved, fd)
		syscall.Close(saved)
		if isStdout {
			os.Stdout = os.NewFile(uintptr(saved), "/dev/stdout")
		} else {
			os.Stderr = os.NewFile(uintptr(saved), "/dev/stderr")
		}
	}()

	if err := syscall.Dup2(int(file.Fd()), fd); err != nil {
		panic(err)
	}

	if isStdout {
		os.Stdout = file
	} else {
		os.Stderr = file
	}
}
