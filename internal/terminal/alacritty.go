package terminal

import (
	"fmt"
	"os"
)

type Alacritty struct{}

func (a *Alacritty) Name() string { return "alacritty" }

func (a *Alacritty) DisplayName() string { return "Alacritty" }

func (a *Alacritty) Detect() bool {
	_, err := os.Stat("/Applications/Alacritty.app")
	return err == nil
}

func (a *Alacritty) Binary() string {
	return "/usr/bin/open"
}

func (a *Alacritty) BuildArgs(cfg LaunchConfig) []string {
	// On macOS, Alacritty must be launched via `open -na Alacritty.app --args ...`
	args := []string{"-na", "Alacritty.app", "--args"}
	args = append(args, "--title", cfg.Title)
	if cfg.FontSize > 0 {
		args = append(args, "-o", fmt.Sprintf("font.size=%d", cfg.FontSize))
	}
	if cfg.NoDecorations {
		args = append(args, "-o", "window.decorations=None")
	}
	args = append(args, "-e", cfg.ScriptPath)
	return args
}
