package terminal

import (
	"fmt"
)

type AppleTerminal struct{}

func (a *AppleTerminal) Name() string { return "terminal" }

func (a *AppleTerminal) DisplayName() string { return "Terminal" }

func (a *AppleTerminal) Detect() bool {
	// Terminal.app is always available on macOS
	return true
}

func (a *AppleTerminal) Binary() string {
	return "/usr/bin/osascript"
}

func (a *AppleTerminal) BuildArgs(cfg LaunchConfig) []string {
	// `do script "X"` runs X as a shell expression in a new Terminal tab.
	// Since cfg.ScriptPath is a path to an executable shell script, passing
	// it as the "command" makes the new tab's login shell run the script.
	scriptPath := escapeAppleScriptString(cfg.ScriptPath)
	title := escapeAppleScriptString(cfg.Title)

	script := fmt.Sprintf(`
tell application "Terminal"
	activate
	do script "%s"
	set custom title of front window to "%s"
end tell`, scriptPath, title)

	return []string{"-e", script}
}
