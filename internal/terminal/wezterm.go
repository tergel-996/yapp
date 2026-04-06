package terminal

import (
	"fmt"
	"os"
	"os/exec"
)

type WezTerm struct{}

func (w *WezTerm) Name() string { return "wezterm" }

func (w *WezTerm) DisplayName() string { return "WezTerm" }

func (w *WezTerm) Detect() bool {
	if _, err := exec.LookPath("wezterm"); err == nil {
		return true
	}
	_, err := os.Stat("/Applications/WezTerm.app")
	return err == nil
}

func (w *WezTerm) Binary() string {
	path, err := exec.LookPath("wezterm")
	if err != nil {
		return "wezterm"
	}
	return path
}

func (w *WezTerm) BuildArgs(cfg LaunchConfig) []string {
	// WezTerm's --config flag is a global option that must come before the
	// subcommand (`wezterm --config K=V start -- ...`). Each --config takes
	// one "name=value" argument where the value is parsed as Lua.
	//
	// Note: WezTerm does not expose a CLI flag for the window title on
	// macOS. Titles are set via terminal escape sequences from the running
	// program, so cfg.Title is intentionally not honored here. `--class` is
	// still passed to set the OS-level window class.
	args := []string{}
	if cfg.FontSize > 0 {
		args = append(args, "--config", fmt.Sprintf("font_size=%d", cfg.FontSize))
	}
	if cfg.NoDecorations {
		// The Lua value must be a quoted string; pass the double-quotes
		// through literally since exec.Command does not invoke a shell.
		args = append(args, "--config", `window_decorations="NONE"`)
	}
	args = append(args,
		"start",
		"--class", cfg.Title,
		"--",
		cfg.ScriptPath,
	)
	return args
}
