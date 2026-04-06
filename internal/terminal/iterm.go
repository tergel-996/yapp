package terminal

import (
	"fmt"
	"os"
)

type ITerm struct{}

func (i *ITerm) Name() string { return "iterm" }

func (i *ITerm) DisplayName() string { return "iTerm" }

func (i *ITerm) Detect() bool {
	_, err := os.Stat("/Applications/iTerm.app")
	return err == nil
}

func (i *ITerm) Binary() string {
	return "/usr/bin/osascript"
}

func (i *ITerm) BuildArgs(cfg LaunchConfig) []string {
	// iTerm2's `create window with default profile command "..."` tokenises
	// its `command` argument via execvp-like splitting, which would mangle
	// a path containing spaces. `write text` pipes text into the session's
	// shell verbatim, so the new shell in the window executes the script
	// path like any other command.
	scriptPath := escapeAppleScriptString(cfg.ScriptPath)
	title := escapeAppleScriptString(cfg.Title)

	script := fmt.Sprintf(`
tell application "iTerm"
	activate
	create window with default profile
	tell current session of current window
		set name to "%s"
		write text "%s"
	end tell
end tell`, title, scriptPath)

	return []string{"-e", script}
}
