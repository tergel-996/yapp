package terminal

import (
	"fmt"
	"os"
)

type Kitty struct{}

func (k *Kitty) Name() string { return "kitty" }

func (k *Kitty) Detect() bool {
	_, err := os.Stat("/Applications/kitty.app")
	return err == nil
}

func (k *Kitty) Binary() string {
	return "/usr/bin/open"
}

func (k *Kitty) BuildArgs(cfg LaunchConfig) []string {
	// On macOS, kitty must be launched via `open -na kitty.app --args ...`
	args := []string{"-na", "kitty.app", "--args"}
	args = append(args, "--title", cfg.Title)
	if cfg.FontSize > 0 {
		args = append(args, "-o", fmt.Sprintf("font_size=%d", cfg.FontSize))
	}
	if cfg.NoDecorations {
		args = append(args, "-o", "hide_window_decorations=yes")
	}
	args = append(args, cfg.Command)
	args = append(args, cfg.Args...)
	return args
}
