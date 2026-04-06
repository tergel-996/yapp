package main

import (
	"os"
	"path/filepath"

	"github.com/tergel/yapp/internal/cli"
)

func main() {
	// When invoked as the .app bundle launcher (Contents/MacOS/Yapp),
	// automatically run the launch subcommand.
	if filepath.Base(os.Args[0]) == "Yapp" {
		os.Args = append([]string{os.Args[0], "launch"}, os.Args[1:]...)
	}
	cli.Execute()
}
