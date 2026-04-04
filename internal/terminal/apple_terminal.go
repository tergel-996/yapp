package terminal

import (
	"fmt"
	"strings"
)

type AppleTerminal struct{}

func (a *AppleTerminal) Name() string { return "terminal" }

func (a *AppleTerminal) Detect() bool {
	// Terminal.app is always available on macOS
	return true
}

func (a *AppleTerminal) Binary() string {
	return "/usr/bin/osascript"
}

func (a *AppleTerminal) BuildArgs(cfg LaunchConfig) []string {
	cmd := cfg.Command
	if len(cfg.Args) > 0 {
		cmd += " " + strings.Join(cfg.Args, " ")
	}

	script := fmt.Sprintf(`
tell application "Terminal"
	activate
	do script "%s"
	set custom title of front window to "%s"
end tell`, cmd, cfg.Title)

	return []string{"-e", script}
}
