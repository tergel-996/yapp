package terminal

import (
	"fmt"
	"os"
)

type Ghostty struct{}

func (g *Ghostty) Name() string { return "ghostty" }

func (g *Ghostty) DisplayName() string { return "Ghostty" }

func (g *Ghostty) Detect() bool {
	_, err := os.Stat("/Applications/Ghostty.app")
	return err == nil
}

func (g *Ghostty) Binary() string {
	return "/usr/bin/open"
}

func (g *Ghostty) BuildArgs(cfg LaunchConfig) []string {
	// On macOS, Ghostty must be launched via `open -na Ghostty.app --args ...`
	// The `ghostty` binary in PATH is the CLI helper and cannot open windows.
	//
	// Ghostty joins every argv element after `-e` into a single shell
	// string and runs it via `bash -c "exec -l <string>"`, so we must pass
	// exactly one word -- the script path -- to avoid any bash tokenisation
	// breaking our shell expression. The script's shebang handles sh
	// dispatch inside the window.
	args := []string{"-na", "Ghostty.app", "--args"}
	args = append(args, fmt.Sprintf("--title=%s", cfg.Title))
	if cfg.FontSize > 0 {
		args = append(args, fmt.Sprintf("--font-size=%d", cfg.FontSize))
	}
	if cfg.NoDecorations {
		args = append(args, "--window-decoration=false")
	} else {
		args = append(args, "--window-decoration=true")
	}
	args = append(args, "--quit-after-last-window-closed=true")
	args = append(args, "-e", cfg.ScriptPath)
	return args
}
