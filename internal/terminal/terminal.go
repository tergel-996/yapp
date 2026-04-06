package terminal

// LaunchConfig describes what to run inside a terminal window.
//
// ScriptPath is the absolute path to an executable /bin/sh script. The
// caller (internal/cli/launch.go) is responsible for writing the script --
// that is where the yazi invocation is wrapped with a trap that removes the
// Yapp marker file on exit, keeping Yapp.app alive in Cmd+Tab for as long
// as yazi is running.
//
// We pass a *path* (not a shell expression) because terminals like Ghostty
// concatenate their -e arguments into a bash -c string, which breaks any
// inner `/bin/sh -c` quoting. A single script path is one argv word under
// every terminal we support, so no tokenisation pitfalls.
type LaunchConfig struct {
	ScriptPath    string
	Title         string
	FontSize      int
	NoDecorations bool
}

type Terminal interface {
	Name() string
	// DisplayName returns the macOS .app display name of the terminal,
	// suitable for use in `tell application "<DisplayName>" to activate`
	// in AppleScript. This is how the Yapp Dock-click handler refocuses
	// the spawned terminal window without needing to track PIDs.
	DisplayName() string
	Detect() bool
	Binary() string
	BuildArgs(cfg LaunchConfig) []string
}
