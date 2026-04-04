package terminal

import (
	"fmt"
	"os"
	"strings"
)

type ITerm struct{}

func (i *ITerm) Name() string { return "iterm" }

func (i *ITerm) Detect() bool {
	_, err := os.Stat("/Applications/iTerm.app")
	return err == nil
}

func (i *ITerm) Binary() string {
	return "/usr/bin/osascript"
}

func (i *ITerm) BuildArgs(cfg LaunchConfig) []string {
	cmd := cfg.Command
	if len(cfg.Args) > 0 {
		cmd += " " + strings.Join(cfg.Args, " ")
	}

	script := fmt.Sprintf(`
tell application "iTerm"
	activate
	set newWindow to (create window with default profile command "%s")
	tell current session of newWindow
		set name to "%s"
	end tell
end tell`, cmd, cfg.Title)

	return []string{"-e", script}
}
