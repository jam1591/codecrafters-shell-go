//go:build windows

package commands

import "os"

func redirectStdout(path string, isAppend bool) {
	redirectStream(&os.Stdout, path, isAppend)
}

func redirectStderr(path string, isAppend bool) {
	redirectStream(&os.Stderr, path, isAppend)
}

func redirectStream(stream **os.File, path string, isAppend bool) {
	file, err := os.OpenFile(path, getFlagForRedirect(isAppend), 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	original := *stream
	defer func() { *stream = original }()
	*stream = file
}
